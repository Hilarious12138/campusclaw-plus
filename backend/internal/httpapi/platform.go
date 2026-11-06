package httpapi

import (
	"campusclaw-plus/internal/domain"
	"campusclaw-plus/internal/service"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Platform struct {
	Auth      service.Auth
	Origin    string
	Resources service.Resources
}

func statusFor(e error) (int, string) {
	switch {
	case e == nil:
		return 200, "OK"
	case errors.Is(e, domain.ErrUnauthenticated):
		return 401, "UNAUTHENTICATED"
	case errors.Is(e, domain.ErrForbidden):
		return 403, "FORBIDDEN"
	case errors.Is(e, domain.ErrNotFound):
		return 404, "NOT_FOUND"
	case errors.Is(e, domain.ErrInvalid):
		return 400, "INVALID_ARGUMENT"
	case errors.Is(e, domain.ErrConflict):
		return 409, "CONFLICT"
	case errors.Is(e, domain.ErrLimited):
		return 429, "RATE_LIMITED"
	default:
		return 503, "DEPENDENCY_UNAVAILABLE"
	}
}
func decode(w http.ResponseWriter, r *http.Request, out any, max int64) error {
	mt, _, e := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if e != nil || mt != "application/json" {
		return domain.ErrInvalid
	}
	r.Body = http.MaxBytesReader(w, r.Body, max)
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if e = d.Decode(out); e != nil {
		return domain.ErrInvalid
	}
	var extra any
	if d.Decode(&extra) != io.EOF {
		return domain.ErrInvalid
	}
	return nil
}
func (p *Platform) serve(w http.ResponseWriter, r *http.Request, id string) bool {
	kind, resourceID := "", ""
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) >= 2 && parts[0] == "api" && (parts[1] == "agents" || parts[1] == "skills" || parts[1] == "prompts") && len(parts) <= 3 {
		kind = parts[1]
		if len(parts) == 3 {
			resourceID = parts[2]
		}
	}
	known := r.URL.Path == "/api/auth/login" || r.URL.Path == "/api/auth/logout" || r.URL.Path == "/api/me" || r.URL.Path == "/api/audit-events"
	known = known || kind != ""
	if !known {
		return false
	}
	start := time.Now()
	var u domain.Identity
	var claims domain.Claims
	var data any
	var e error
	method := "GET"
	if strings.HasPrefix(r.URL.Path, "/api/auth/") {
		method = "POST"
	}
	event := domain.Audit{OccurredAt: start.UTC(), Operation: r.Method + " " + r.URL.Path, ObjectType: "identity", ObjectID: "", RequestID: id, Usage: "unknown"}
	login := r.URL.Path == "/api/auth/login"
	if !login {
		h := r.Header.Values("Authorization")
		if len(h) != 1 {
			e = domain.ErrUnauthenticated
		} else {
			parts := strings.Fields(h[0])
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				e = domain.ErrUnauthenticated
			} else {
				u, claims, e = p.Auth.Identify(r.Context(), parts[1])
			}
		}
	}
	if e == nil && kind == "" && r.Method != method {
		e = domain.ErrInvalid
	}
	if e == nil && r.Method != "GET" && r.Header.Get("Origin") != "" && r.Header.Get("Origin") != p.Origin {
		e = domain.ErrForbidden
	}
	if e == nil && r.URL.Path != "/api/audit-events" && r.URL.RawQuery != "" {
		e = domain.ErrInvalid
	}
	audited := false
	if e == nil && kind != "" {
		event.ActorID = u.ID
		event.ClassID = u.ClassID
		event.ObjectType = kind
		event.ObjectID = resourceID
		event.HTTPStatus = 200
		event.Outcome = "OK"
		switch r.Method {
		case "GET":
			if resourceID == "" {
				data, e = p.Resources.List(r.Context(), u, kind)
			} else {
				data, e = p.Resources.Read(r.Context(), u, kind, resourceID)
			}
		case "POST", "PUT":
			if u.Role != "teacher" {
				e = domain.ErrForbidden
				break
			}
			if (r.Method == "POST" && resourceID != "") || (r.Method == "PUT" && resourceID == "") {
				e = domain.ErrInvalid
				break
			}
			var raw json.RawMessage
			e = decode(w, r, &raw, 65536)
			revision := 0
			if r.Method == "PUT" {
				revision, _ = strconv.Atoi(r.Header.Get("If-Match"))
				if revision < 1 {
					e = domain.ErrInvalid
				}
			}
			if e == nil {
				data, e = p.Resources.Save(r.Context(), u, domain.Resource{ID: resourceID, Kind: kind, Revision: revision, Body: raw}, event)
				audited = e == nil
			}
		case "DELETE":
			if resourceID == "" {
				e = domain.ErrInvalid
			} else {
				e = p.Resources.Delete(r.Context(), u, kind, resourceID, event)
				audited = e == nil
				data = map[string]bool{"deleted": e == nil}
			}
		default:
			e = domain.ErrInvalid
		}
	}
	if e == nil && kind == "" {
		switch r.URL.Path {
		case "/api/auth/login":
			var b struct {
				Username string `json:"username"`
				Password string `json:"password"`
			}
			e = decode(w, r, &b, 8192)
			if e == nil {
				host, _, _ := net.SplitHostPort(r.RemoteAddr)
				var token string
				token, u, e = p.Auth.Login(r.Context(), b.Username, b.Password, host)
				if e == nil {
					data = map[string]any{"access_token": token, "token_type": "Bearer", "expires_in": 1800, "user": u}
				}
			}
		case "/api/me":
			data = u
		case "/api/auth/logout":
			data = map[string]bool{"logged_out": true}
		case "/api/audit-events":
			var f domain.AuditFilter
			f, e = auditFilter(r)
			if e == nil {
				data, e = p.Auth.Events(r.Context(), u, f)
			}
		}
	}
	event.ActorID = u.ID
	event.ClassID = u.ClassID
	event.HTTPStatus, event.Outcome = statusFor(e)
	event.DurationMS = time.Since(start).Milliseconds()
	if r.URL.Path == "/api/auth/logout" && e == nil {
		e = p.Auth.Revoke(r.Context(), claims, event)
	} else if !audited {
		if ae := p.Auth.Audit(r.Context(), event); ae != nil {
			e = ae
		}
	}
	status, code := statusFor(e)
	out := envelope{Code: code, RequestID: id}
	if e == nil {
		out.Data = data
	} else {
		out.Error = &problem{Message: "请求未完成，请检查权限、参数或服务状态。"}
	}
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(out)
	return true
}
func auditFilter(r *http.Request) (domain.AuditFilter, error) {
	f := domain.AuditFilter{Limit: 20}
	v := r.URL.Query()
	for k, x := range v {
		if len(x) != 1 {
			return f, domain.ErrInvalid
		}
		switch k {
		case "actor_id", "operation", "from", "until", "status", "limit", "offset":
		default:
			return f, domain.ErrInvalid
		}
	}
	f.ActorID = v.Get("actor_id")
	f.Operation = v.Get("operation")
	if len(f.ActorID) > 32 || len(f.Operation) > 128 {
		return f, domain.ErrInvalid
	}
	for k, target := range map[string]*int{"limit": &f.Limit, "offset": &f.Offset, "status": &f.Status} {
		if s := v.Get(k); s != "" {
			n, e := strconv.Atoi(s)
			if e != nil {
				return f, domain.ErrInvalid
			}
			*target = n
		}
	}
	if f.Limit < 1 || f.Limit > 100 || f.Offset < 0 || f.Offset > 100000 || f.Status != 0 && (f.Status < 100 || f.Status > 599) {
		return f, domain.ErrInvalid
	}
	for k, target := range map[string]*time.Time{"from": &f.From, "until": &f.Until} {
		if s := v.Get(k); s != "" {
			t, e := time.Parse(time.RFC3339, s)
			if e != nil {
				return f, domain.ErrInvalid
			}
			*target = t
		}
	}
	if !f.From.IsZero() && !f.Until.IsZero() && f.Until.Before(f.From) {
		return f, domain.ErrInvalid
	}
	return f, nil
}
