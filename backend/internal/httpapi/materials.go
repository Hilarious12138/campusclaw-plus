package httpapi

import (
	"campusclaw-plus/internal/domain"

	"io"
	"net/http"
	"strconv"
	"strings"
)

func (p *Platform) business(w http.ResponseWriter, r *http.Request, u domain.Identity, event domain.Audit) (any, error, bool, int) {
	path := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	event.ActorID = u.ID
	event.ClassID = u.ClassID
	event.HTTPStatus = 202
	event.Outcome = "ACCEPTED"
	if r.URL.Path == "/api/materials" {
		if r.Method == "GET" {
			f, e := pagination(r)
			if e != nil {
				return nil, e, false, 0
			}
			v, e := p.Materials.List(r.Context(), u, f.Limit, f.Offset)
			return v, e, false, 0
		}
		if r.Method != "POST" {
			return nil, domain.ErrInvalid, false, 0
		}
		if u.Role != "teacher" {
			return nil, domain.ErrForbidden, false, 0
		}
		if r.URL.RawQuery != "" {
			return nil, domain.ErrInvalid, false, 0
		}
		r.Body = http.MaxBytesReader(w, r.Body, (2<<20)+65536)
		if r.ParseMultipartForm((2<<20)+65536) != nil {
			return nil, domain.ErrInvalid, false, 0
		}
		defer r.MultipartForm.RemoveAll()
		if len(r.MultipartForm.Value) != 0 || len(r.MultipartForm.File) != 1 || len(r.MultipartForm.File["file"]) != 1 {
			return nil, domain.ErrInvalid, false, 0
		}
		f, h, e := r.FormFile("file")
		if e != nil {
			return nil, domain.ErrInvalid, false, 0
		}
		defer f.Close()
		raw, e := io.ReadAll(io.LimitReader(f, (2<<20)+1))
		if e != nil {
			return nil, domain.ErrInvalid, false, 0
		}
		if len(r.Header.Values("Idempotency-Key")) != 1 {
			return nil, domain.ErrInvalid, false, 0
		}
		m, t, e := p.Materials.Create(r.Context(), u, h.Filename, string(raw), r.Header.Get("Idempotency-Key"), event)
		if e != nil {
			return nil, e, false, 0
		}
		return map[string]any{"material_id": m.ID, "task_id": t.ID, "status": t.Status, "status_url": "/api/tasks/" + t.ID}, nil, true, 202
	}
	if len(path) == 3 && path[1] == "tasks" && r.Method == "GET" {
		v, e := p.Materials.Task(r.Context(), u, path[2])
		return v, e, false, 0
	}
	if len(path) >= 3 && path[1] == "materials" && r.Method == "GET" {
		v, e := p.Materials.Read(r.Context(), u, path[2])
		if e != nil {
			return nil, e, false, 0
		}
		if len(path) == 4 && path[3] == "file" {
			return map[string]any{"title": v.Title, "content": v.Content}, nil, false, 0
		}
		if len(path) == 3 {
			return v, nil, false, 0
		}
	}
	if r.URL.Path == "/api/retrieval/search" && r.Method == "POST" {
		var b struct {
			Query string `json:"query"`
			Mode  string `json:"mode"`
		}
		if decode(w, r, &b, 8192) != nil {
			return nil, domain.ErrInvalid, false, 0
		}
		v, e := p.Retrieval.Search(r.Context(), u, b.Query, b.Mode, nil)
		return v, e, false, 0
	}
	if len(path) == 4 && path[1] == "retrieval" && path[2] == "chunks" && r.Method == "GET" {
		v, e := p.Retrieval.Store.Chunk(r.Context(), u, path[3])
		return v, e, false, 0
	}
	return nil, domain.ErrNotFound, false, 0
}
func pagination(r *http.Request) (domain.AuditFilter, error) {
	f := domain.AuditFilter{Limit: 20}
	for k, v := range r.URL.Query() {
		if len(v) != 1 {
			return f, domain.ErrInvalid
		}
		n, e := strconv.Atoi(v[0])
		if e != nil {
			return f, domain.ErrInvalid
		}
		switch k {
		case "limit":
			f.Limit = n
		case "offset":
			f.Offset = n
		default:
			return f, domain.ErrInvalid
		}
	}
	if f.Limit < 1 || f.Limit > 100 || f.Offset < 0 || f.Offset > 100000 {
		return f, domain.ErrInvalid
	}
	return f, nil
}
