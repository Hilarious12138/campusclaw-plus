package httpapi

import (
	"campusclaw-plus/internal/service"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log/slog"
	"net/http"
	"regexp"
	"time"
)

var safeID = regexp.MustCompile(`^[A-Za-z0-9_-]{1,64}$`)

type envelope struct {
	Code      string   `json:"code"`
	Data      any      `json:"data,omitempty"`
	Error     *problem `json:"error,omitempty"`
	RequestID string   `json:"request_id"`
}
type problem struct {
	Message string `json:"message"`
}
type response struct {
	http.ResponseWriter
	status int
}

func (w *response) WriteHeader(n int) { w.status = n; w.ResponseWriter.WriteHeader(n) }
func Handler(s service.Health, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		id := r.Header.Get("X-Request-ID")
		if !safeID.MatchString(id) {
			var b [16]byte
			_, _ = rand.Read(b[:])
			id = hex.EncodeToString(b[:])
		}
		w.Header().Set("X-Request-ID", id)
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		rw := &response{w, 200}
		defer func() {
			logger.Info("request", "module", "httpapi", "event", "request", "request_id", id, "method", r.Method, "path", r.URL.Path, "status", rw.status, "duration_ms", time.Since(start).Milliseconds())
		}()
		reply := func(status int, code, message string, data any) {
			v := envelope{Code: code, Data: data, RequestID: id}
			if message != "" {
				v.Error = &problem{message}
			}
			rw.WriteHeader(status)
			_ = json.NewEncoder(rw).Encode(v)
		}
		if r.URL.Path != "/health" && r.URL.Path != "/live" {
			reply(404, "NOT_FOUND", "资源不存在或未开放，请检查访问路径。", nil)
			return
		}
		if r.Method != "GET" {
			w.Header().Set("Allow", "GET")
			reply(405, "METHOD_NOT_ALLOWED", "该入口仅允许GET，请按接口说明重试。", nil)
			return
		}
		if r.URL.RawQuery != "" {
			reply(400, "INVALID_ARGUMENT", "该入口不接受查询参数，请移除参数重试。", nil)
			return
		}
		if r.URL.Path == "/live" {
			reply(200, "OK", "", map[string]string{"status": "alive"})
			return
		}
		h, ok := s.Check(r.Context())
		if !ok {
			reply(503, "DEPENDENCY_UNAVAILABLE", "数据库暂不可用，请检查服务后重试。", h)
			return
		}
		reply(200, "OK", "", h)
	})
}
