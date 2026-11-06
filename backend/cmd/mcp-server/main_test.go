package main

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestProtocol(t *testing.T) {
	token := strings.Repeat("x", 32)
	h := handler(token)
	for _, tc := range []struct {
		body, auth string
		status     int
		contains   string
	}{
		{`{"jsonrpc":"2.0","id":1,"method":"tools/list"}`, "", 401, "error"},
		{`{"jsonrpc":"2.0","id":1,"method":"tools/list"}`, "Bearer wrong", 401, "error"},
		{`{"jsonrpc":"2.0","id":1,"method":"tools/list"}`, "Bearer " + token, 200, "formula_hint"},
		{`broken`, "Bearer " + token, 200, "-32700"},
		{`[]`, "Bearer " + token, 200, "-32600"},
		{`{"jsonrpc":"2.0","id":true,"method":"tools/list"}`, "Bearer " + token, 200, `"id":null`},
		{`{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"unit_convert","arguments":{"question":"1 km 转 m","extra":1}}}`, "Bearer " + token, 200, "-32602"},
		{`{"jsonrpc":"2.0","id":1,"method":"missing"}`, "Bearer " + token, 200, "-32601"},
		{`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-03-26"}}`, "Bearer " + token, 200, "protocolVersion"},
		{`{"jsonrpc":"2.0","method":"notifications/initialized"}`, "Bearer " + token, 202, ""},
		{`{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"unit_convert","arguments":{"question":"1 km 转 m"}}}`, "Bearer " + token, 200, "1000"},
	} {
		r := httptest.NewRequest("POST", "/mcp", strings.NewReader(tc.body))
		r.Header.Set("Authorization", tc.auth)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Code != tc.status || !strings.Contains(w.Body.String(), tc.contains) {
			t.Fatal(tc.status, w.Code, w.Body.String())
		}
		if w.Code == 200 {
			var result any
			if json.Unmarshal(w.Body.Bytes(), &result) != nil {
				t.Fatal("invalid JSON")
			}
		}
	}
}
func TestConvert(t *testing.T) {
	for _, q := range []string{"1 kg 转 m", "请换算", "1 m 转 kg"} {
		if _, bad := runTool("unit_convert", q); !bad {
			t.Fatal(q)
		}
	}
	if text, bad := runTool("unit_convert", "100 cm to m"); bad || !strings.Contains(text, "1 m") {
		t.Fatal(text)
	}
}
