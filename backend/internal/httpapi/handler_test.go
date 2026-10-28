package httpapi

import (
	"bytes"
	"campusclaw-plus/internal/service"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type probe func(context.Context) error

func (p probe) Ping(c context.Context) error { return p(c) }
func setup(p probe) (*bytes.Buffer, service.Health) {
	b := &bytes.Buffer{}
	return b, service.Health{DB: p, Version: "test-version", Timeout: 100 * time.Millisecond, Now: func() time.Time { return time.Date(2026, 10, 28, 0, 0, 0, 0, time.UTC) }}
}
func TestNormalHealthLogS08S10(t *testing.T) {
	b, s := setup(func(context.Context) error { return nil })
	r := httptest.NewRequest("GET", "/health", nil)
	r.Header.Set("X-Request-ID", "verify-lesson07")
	w := httptest.NewRecorder()
	Handler(s, slog.New(slog.NewJSONHandler(b, nil))).ServeHTTP(w, r)
	var v struct {
		Code, RequestID string
		Data            struct {
			Version      string
			Dependencies map[string]string
		}
	}
	if e := json.Unmarshal(w.Body.Bytes(), &v); e != nil {
		t.Fatal(e)
	}
	var raw map[string]any
	json.Unmarshal(w.Body.Bytes(), &raw)
	if w.Code != 200 || v.Code != "OK" || v.Data.Version != "test-version" || v.Data.Dependencies["mysql"] != "up" || raw["request_id"] != "verify-lesson07" || w.Header().Get("X-Request-ID") != "verify-lesson07" {
		t.Fatal(w.Body.String())
	}
	var log map[string]any
	if json.Unmarshal(b.Bytes(), &log) != nil || log["request_id"] != "verify-lesson07" {
		t.Fatal("log correlation")
	}
	for _, k := range []string{"time", "level", "module", "event", "request_id", "method", "path", "status", "duration_ms"} {
		if _, ok := log[k]; !ok {
			t.Fatal("missing log field", k)
		}
	}
}
func TestFailurePathsS04S05S06S07(t *testing.T) {
	cases := []struct {
		method, path, code string
		status             int
	}{{"GET", "/health?class_id=other", "INVALID_ARGUMENT", 400}, {"GET", "/api/teachers/materials", "NOT_FOUND", 404}, {"GET", "/api/missing", "NOT_FOUND", 404}, {"GET", "/uploads/guessed.txt", "NOT_FOUND", 404}, {"POST", "/health", "METHOD_NOT_ALLOWED", 405}}
	for _, v := range cases {
		t.Run(v.method+v.path, func(t *testing.T) {
			b, s := setup(func(context.Context) error { t.Fatal("rejected request probed database"); return nil })
			w := httptest.NewRecorder()
			Handler(s, slog.New(slog.NewJSONHandler(b, nil))).ServeHTTP(w, httptest.NewRequest(v.method, v.path, nil))
			var out envelope
			json.Unmarshal(w.Body.Bytes(), &out)
			if w.Code != v.status || out.Code != v.code || out.Data != nil {
				t.Fatal(w.Code, w.Body.String())
			}
		})
	}
}
func TestFailureDependencyS09(t *testing.T) {
	b, s := setup(func(c context.Context) error { <-c.Done(); return errors.New("private-driver-detail") })
	handler := Handler(s, slog.New(slog.NewJSONHandler(b, nil)))
	w := httptest.NewRecorder()
	start := time.Now()
	handler.ServeHTTP(w, httptest.NewRequest("GET", "/health", nil))
	if w.Code != 503 || !strings.Contains(w.Body.String(), "DEPENDENCY_UNAVAILABLE") || strings.Contains(w.Body.String(), "private-driver-detail") || time.Since(start) > time.Second {
		t.Fatal(w.Body.String())
	}
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, httptest.NewRequest("GET", "/live", nil))
	if w.Code != 200 {
		t.Fatal("liveness depends on DB")
	}
}
func TestBoundaryRequestIDS11(t *testing.T) {
	for _, id := range []string{"x", strings.Repeat("a", 64), strings.Repeat("a", 65), "bad id"} {
		t.Run(id, func(t *testing.T) {
			b, s := setup(func(context.Context) error { return nil })
			r := httptest.NewRequest("GET", "/live", nil)
			r.Header.Set("X-Request-ID", id)
			w := httptest.NewRecorder()
			Handler(s, slog.New(slog.NewJSONHandler(b, nil))).ServeHTTP(w, r)
			got := w.Header().Get("X-Request-ID")
			if !safeID.MatchString(got) || (safeID.MatchString(id) && got != id) || (!safeID.MatchString(id) && got == id) {
				t.Fatal("unsafe request ID")
			}
		})
	}
}
