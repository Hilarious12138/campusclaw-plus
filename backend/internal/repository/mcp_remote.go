package repository

import (
	"bytes"
	"campusclaw-plus/internal/domain"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type MCPRemote struct {
	AllowedURLs []string
	Token       string
}

var toolName = regexp.MustCompile(`^[a-z][a-z0-9_]{0,63}$`)

func (m MCPRemote) Allowed(base string) bool {
	u, e := url.Parse(base)
	if e != nil || u.User != nil || u.RawQuery != "" || u.Fragment != "" || u.Path != "" && u.Path != "/" || (u.Scheme != "http" && u.Scheme != "https") {
		return false
	}
	for _, v := range m.AllowedURLs {
		if strings.TrimRight(v, "/") == strings.TrimRight(base, "/") {
			return true
		}
	}
	return false
}
func (m MCPRemote) rpc(ctx context.Context, base, method string, params, out any) error {
	if !m.Allowed(base) || len(m.Token) < 32 {
		return domain.ErrForbidden
	}
	transport := &http.Transport{Proxy: nil, ResponseHeaderTimeout: 5 * time.Second, DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
		host, port, e := net.SplitHostPort(address)
		if e != nil {
			return nil, e
		}
		ips, e := net.DefaultResolver.LookupIPAddr(ctx, host)
		if e != nil {
			return nil, e
		}
		for _, ip := range ips {
			if host != "mcp-demo" && (ip.IP.IsPrivate() || ip.IP.IsLoopback() || ip.IP.IsLinkLocalUnicast() || ip.IP.IsUnspecified()) {
				continue
			}
			return (&net.Dialer{Timeout: 5 * time.Second}).DialContext(ctx, network, net.JoinHostPort(ip.IP.String(), port))
		}
		return nil, domain.ErrForbidden
	}}
	defer transport.CloseIdleConnections()
	client := http.Client{Transport: transport, Timeout: 5 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	request := map[string]any{"jsonrpc": "2.0", "id": "plus", "method": method, "params": params}
	if method == "notifications/initialized" {
		delete(request, "id")
	}
	raw, _ := json.Marshal(request)
	r, e := http.NewRequestWithContext(ctx, "POST", strings.TrimRight(base, "/")+"/mcp", bytes.NewReader(raw))
	if e != nil {
		return e
	}
	r.Header.Set("Authorization", "Bearer "+m.Token)
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json, text/event-stream")
	r.Header.Set("MCP-Protocol-Version", "2025-03-26")
	res, e := client.Do(r)
	if e != nil {
		return errors.New("MCP unavailable")
	}
	defer res.Body.Close()
	if method == "notifications/initialized" {
		if res.StatusCode == 202 || res.StatusCode == 204 {
			return nil
		}
		return errors.New("MCP initialization failed")
	}
	if res.StatusCode != 200 {
		return errors.New("MCP failed")
	}
	b, e := io.ReadAll(io.LimitReader(res.Body, 65537))
	if e != nil || len(b) > 65536 {
		return errors.New("MCP response too large")
	}
	var v struct {
		JSONRPC string          `json:"jsonrpc"`
		ID      string          `json:"id"`
		Result  json.RawMessage `json:"result"`
		Error   json.RawMessage `json:"error"`
	}
	if json.Unmarshal(b, &v) != nil || v.JSONRPC != "2.0" || v.ID != "plus" || len(v.Error) > 0 {
		return errors.New("MCP invalid response")
	}
	return json.Unmarshal(v.Result, out)
}
func (m MCPRemote) Discover(ctx context.Context, base string) ([]domain.MCPTool, error) {
	var init map[string]any
	if e := m.rpc(ctx, base, "initialize", map[string]any{"protocolVersion": "2025-03-26", "capabilities": map[string]any{}, "clientInfo": map[string]any{"name": "campusclaw-plus", "version": "1"}}, &init); e != nil {
		return nil, e
	}
	if e := m.rpc(ctx, base, "notifications/initialized", map[string]any{}, nil); e != nil {
		return nil, e
	}
	var out struct {
		Tools []domain.MCPTool `json:"tools"`
	}
	if e := m.rpc(ctx, base, "tools/list", map[string]any{}, &out); e != nil {
		return nil, e
	}
	if len(out.Tools) > 32 {
		return nil, domain.ErrInvalid
	}
	seen := map[string]bool{}
	for _, t := range out.Tools {
		if !toolName.MatchString(t.Name) || seen[t.Name] {
			return nil, domain.ErrInvalid
		}
		seen[t.Name] = true
	}
	return out.Tools, nil
}
func (m MCPRemote) Call(ctx context.Context, base, name, question string) (string, error) {
	var out struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
		IsError bool `json:"isError"`
	}
	if e := m.rpc(ctx, base, "tools/call", map[string]any{"name": name, "arguments": map[string]any{"question": question}}, &out); e != nil {
		return "", e
	}
	if out.IsError {
		return "", errors.New("MCP tool error")
	}
	var texts []string
	for _, c := range out.Content {
		if c.Type == "text" {
			texts = append(texts, c.Text)
		}
	}
	v := strings.Join(texts, "\n")
	if v == "" || len(v) > 8192 {
		return "", errors.New("MCP invalid result")
	}
	return v, nil
}
