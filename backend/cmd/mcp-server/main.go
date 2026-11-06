package main

import (
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type rpcRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}
type tool struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Schema      map[string]any `json:"inputSchema"`
}

func tools() []tool {
	schema := map[string]any{"type": "object", "properties": map[string]any{"question": map[string]any{"type": "string", "minLength": 2, "maxLength": 500}}, "required": []string{"question"}, "additionalProperties": false}
	return []tool{{"formula_hint", "给出方程、面积等公式的适用条件及下一步思考提示，不解出最终答案。", schema}, {"unit_convert", "解析明确单位换算，如 1 km 转 m；支持 m/cm/km 和 g/kg。", schema}}
}

var conversion = regexp.MustCompile(`(?i)([-+]?(?:\d+(?:\.\d*)?|\.\d+))\s*(km|cm|kg|m|g)\s*(?:转为?|换算为?|转换为?|到|to|=|→)\s*(km|cm|kg|m|g)\b`)

func runTool(name, question string) (string, bool) {
	if name == "formula_hint" {
		switch {
		case strings.Contains(question, "方程") || strings.Contains(question, "等式"):
			return "方程提示：保持等式两边相等，先确认未知数和常数项；思考两边同时进行什么相同运算能消去常数项。", false
		case strings.Contains(question, "面积"):
			return "面积提示：先识别图形与已知长度；长方形面积=长×宽，三角形面积=底×高÷2。先检查长度单位是否一致。", false
		default:
			return "先列出已知量、所求量与适用公式，检查公式的前提条件，再决定下一步；需要更具体的题目才能选定公式。", false
		}
	}
	if name != "unit_convert" {
		return "未知工具", true
	}
	matches := conversion.FindAllStringSubmatch(question, -1)
	if len(matches) != 1 {
		return "请给出明确换算，例如 1 km 转 m；仅支持长度 m/cm/km 和质量 g/kg。", true
	}
	m := matches[0]
	value, e := strconv.ParseFloat(m[1], 64)
	if e != nil || math.IsInf(value, 0) || math.IsNaN(value) {
		return "数值无效。", true
	}
	from, to := strings.ToLower(m[2]), strings.ToLower(m[3])
	factor := map[string]float64{"m": 1, "cm": .01, "km": 1000, "g": 1, "kg": 1000}
	mass := func(s string) bool { return s == "g" || s == "kg" }
	if mass(from) != mass(to) {
		return "单位维度不同，不能在长度和质量之间换算。", true
	}
	result := value * factor[from] / factor[to]
	if math.IsInf(result, 0) || math.IsNaN(result) {
		return "换算结果超出范围。", true
	}
	return fmt.Sprintf("单位换算：%g %s = %g %s。换算系数为 %g，请检查题目中的量与单位。", value, from, result, to, factor[from]/factor[to]), false
}
func handler(token string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status":"ok"}`)
	})
	mux.HandleFunc("POST /mcp", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-store")
		if subtle.ConstantTimeCompare([]byte(r.Header.Get("Authorization")), []byte("Bearer "+token)) != 1 {
			w.Header().Set("WWW-Authenticate", "Bearer")
			w.WriteHeader(401)
			fmt.Fprint(w, `{"error":"MCP Bearer token required"}`)
			return
		}
		var req rpcRequest
		reply := func(result any, code int, message string) {
			var id any
			if len(req.ID) > 0 {
				_ = json.Unmarshal(req.ID, &id)
			}
			body := map[string]any{"jsonrpc": "2.0", "id": id}
			if code != 0 {
				body["error"] = map[string]any{"code": code, "message": message}
			} else {
				body["result"] = result
			}
			_ = json.NewEncoder(w).Encode(body)
		}
		d := json.NewDecoder(http.MaxBytesReader(w, r.Body, 64<<10))
		var raw json.RawMessage
		if e := d.Decode(&raw); e != nil {
			reply(nil, -32700, "Invalid JSON request")
			return
		}
		var extra any
		if d.Decode(&extra) != io.EOF {
			reply(nil, -32700, "Invalid JSON request")
			return
		}
		if json.Unmarshal(raw, &req) != nil {
			reply(nil, -32600, "Expected one JSON-RPC request object")
			return
		}
		if req.JSONRPC != "2.0" || req.Method == "" {
			reply(nil, -32600, "Expected JSON-RPC 2.0 request")
			return
		}
		if len(req.ID) > 0 {
			var id any
			_ = json.Unmarshal(req.ID, &id)
			switch id.(type) {
			case string, float64:
			default:
				req.ID = nil
				reply(nil, -32600, "Request id must be a string or number")
				return
			}
		}
		if req.Method == "notifications/initialized" && len(req.ID) == 0 {
			w.WriteHeader(202)
			return
		}
		if len(req.ID) == 0 {
			reply(nil, -32600, "Request id required")
			return
		}
		switch req.Method {
		case "initialize":
			var p struct {
				Version string `json:"protocolVersion"`
			}
			if json.Unmarshal(req.Params, &p) != nil || p.Version == "" {
				reply(nil, -32602, "protocolVersion required")
				return
			}
			reply(map[string]any{"protocolVersion": "2025-03-26", "capabilities": map[string]any{"tools": map[string]any{}}, "serverInfo": map[string]string{"name": "campus-tools", "version": "0.6.0"}}, 0, "")
		case "ping":
			reply(map[string]any{}, 0, "")
		case "tools/list":
			reply(map[string]any{"tools": tools()}, 0, "")
		case "tools/call":
			var p struct {
				Name      string `json:"name"`
				Arguments struct {
					Question string `json:"question"`
				} `json:"arguments"`
			}
			params := json.NewDecoder(strings.NewReader(string(req.Params)))
			params.DisallowUnknownFields()
			if params.Decode(&p) != nil || (p.Name != "formula_hint" && p.Name != "unit_convert") || len([]rune(strings.TrimSpace(p.Arguments.Question))) < 2 || len([]rune(p.Arguments.Question)) > 500 {
				reply(nil, -32602, "Unknown tool or invalid question (2-500 characters required)")
				return
			}
			text, bad := runTool(p.Name, p.Arguments.Question)
			reply(map[string]any{"content": []any{map[string]string{"type": "text", "text": text}}, "isError": bad}, 0, "")
		default:
			reply(nil, -32601, "Method not found")
		}
	})
	return mux
}
func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	token := os.Getenv("MCP_AUTH_TOKEN")
	if len(token) < 32 {
		logger.Error("invalid_configuration", "module", "mcp", "event", "failure", "request_id", "system")
		os.Exit(1)
	}
	server := http.Server{Addr: ":8090", Handler: handler(token), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second, WriteTimeout: 10 * time.Second, MaxHeaderBytes: 8192}
	if e := server.ListenAndServe(); e != nil && e != http.ErrServerClosed {
		logger.Error("server_failed", "module", "mcp", "event", "failure", "request_id", "system")
		os.Exit(1)
	}
}
