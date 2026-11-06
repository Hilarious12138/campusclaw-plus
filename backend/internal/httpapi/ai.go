package httpapi

import (
	"campusclaw-plus/internal/domain"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func (p *Platform) ai(w http.ResponseWriter, r *http.Request, u domain.Identity, event domain.Audit) (any, error, bool, int) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	event.ActorID = u.ID
	event.ClassID = u.ClassID
	event.HTTPStatus = 200
	event.Outcome = "OK"
	if r.URL.Path == "/api/mcp/servers" {
		if r.Method == "GET" {
			v, e := p.MCP.Store.MCPServers(r.Context(), u)
			return v, e, false, 0
		}
		if r.Method == "POST" {
			if u.Role != "teacher" {
				return nil, domain.ErrForbidden, false, 0
			}
			var v struct {
				Name    string `json:"name"`
				BaseURL string `json:"base_url"`
			}
			if decode(w, r, &v, 8192) != nil {
				return nil, domain.ErrInvalid, false, 0
			}
			s, e := p.MCP.Register(r.Context(), u, domain.MCPServer{Name: v.Name, BaseURL: v.BaseURL}, event)
			return s, e, e == nil, 0
		}
	}
	if len(parts) == 5 && parts[1] == "mcp" && parts[2] == "servers" {
		if u.Role != "teacher" {
			return nil, domain.ErrForbidden, false, 0
		}
		if r.Method == "POST" && (parts[4] == "connect" || parts[4] == "disconnect") {
			v, e := p.MCP.Connect(r.Context(), u, parts[3], parts[4] == "connect", event)
			return v, e, e == nil, 0
		}
		if r.Method == "GET" && parts[4] == "tools" {
			v, e := p.MCP.Store.MCPServer(r.Context(), u, parts[3])
			if e != nil {
				return nil, e, false, 0
			}
			if !v.Connected {
				return []domain.MCPTool{}, nil, false, 0
			}
			return v.Tools, nil, false, 0
		}
	}
	if r.URL.Path == "/api/conversations" {
		if r.Method == "GET" {
			v, e := p.Conversations.Store.Conversations(r.Context(), u)
			return v, e, false, 0
		}
		if r.Method == "POST" {
			var b struct {
				AgentID string `json:"agent_id"`
			}
			if decode(w, r, &b, 8192) != nil {
				return nil, domain.ErrInvalid, false, 0
			}
			v, e := p.Conversations.Store.CreateConversation(r.Context(), u, b.AgentID, event)
			return v, e, e == nil, 0
		}
	}
	if len(parts) == 3 && parts[1] == "conversations" && r.Method == "GET" {
		v, turns, e := p.Conversations.Store.Conversation(r.Context(), u, parts[2])
		return map[string]any{"conversation": v, "turns": turns}, e, false, 0
	}
	if r.URL.Path == "/api/skills/import" && r.Method == "POST" {
		if u.Role != "teacher" {
			return nil, domain.ErrForbidden, false, 0
		}
		r.Body = http.MaxBytesReader(w, r.Body, 70000)
		if r.ParseMultipartForm(70000) != nil {
			return nil, domain.ErrInvalid, false, 0
		}
		defer r.MultipartForm.RemoveAll()
		file, h, e := r.FormFile("file")
		if e != nil {
			return nil, domain.ErrInvalid, false, 0
		}
		defer file.Close()
		if !strings.HasSuffix(strings.ToLower(h.Filename), ".md") && !strings.HasSuffix(strings.ToLower(h.Filename), ".txt") {
			return nil, domain.ErrInvalid, false, 0
		}
		raw, e := io.ReadAll(io.LimitReader(file, 65537))
		if e != nil || len(raw) > 65536 {
			return nil, domain.ErrInvalid, false, 0
		}
		skill, e := p.Conversations.ParseSkill(string(raw))
		return map[string]any{"preview": skill, "requires_confirmation": true}, e, false, 0
	}
	return nil, domain.ErrNotFound, false, 0
}
func (p *Platform) stream(w http.ResponseWriter, r *http.Request, u domain.Identity, event domain.Audit, id, conv string) (bool, error) {
	var b struct {
		Content string `json:"content"`
	}
	if decode(w, r, &b, 16384) != nil {
		return false, domain.ErrInvalid
	}
	started := false
	emit := func(kind string, v any) error {
		if !started {
			started = true
			w.Header().Set("Content-Type", "text/event-stream")
			w.Header().Set("X-Accel-Buffering", "no")
			w.WriteHeader(200)
		}
		raw, e := json.Marshal(v)
		if e != nil {
			return e
		}
		if _, e = io.WriteString(w, "event: "+kind+"\ndata: "+string(raw)+"\n\n"); e != nil {
			return e
		}
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		return nil
	}
	event.ActorID = u.ID
	event.ClassID = u.ClassID
	event.HTTPStatus = 200
	event.Outcome = "OK"
	e := p.Conversations.Send(r.Context(), u, conv, b.Content, event, emit)
	if started && e != nil {
		_, code := statusFor(e)
		_ = emit("error", map[string]string{"code": code, "message": "回答未完成，请稍后重试。"})
	}
	return started, e
}
