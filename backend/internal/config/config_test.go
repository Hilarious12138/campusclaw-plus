package config

import (
	"strings"
	"testing"
)

func valid() map[string]string {
	return map[string]string{"JWT_SECRET": strings.Repeat("j", 32), "APP_ORIGIN": "http://localhost:3001", "SEED_TEACHER_A_PASSWORD": strings.Repeat("s", 24), "SEED_STUDENT_A_PASSWORD": strings.Repeat("t", 24), "SEED_TEACHER_B_PASSWORD": strings.Repeat("u", 24), "SEED_STUDENT_B_PASSWORD": strings.Repeat("v", 24), "BROKER_ADDR": "broker:5672", "BROKER_USER": "plus", "BROKER_PASSWORD": strings.Repeat("b", 24), "EMBEDDING_BASE_URL": "http://fixture:8091/v1", "EMBEDDING_API_KEY": strings.Repeat("e", 24), "EMBEDDING_MODEL": "test-embedding", "EMBEDDING_DIM": "8", "QDRANT_URL": "http://qdrant:6333", "CHAT_BASE_URL": "http://fixture:8091/v1", "CHAT_API_KEY": strings.Repeat("c", 24), "CHAT_MODEL": "fixture-only", "MCP_AUTH_TOKEN": strings.Repeat("m", 32), "MCP_ALLOWED_BASE_URLS": "http://mcp-demo:8090", "TASK_TIMEOUT_SECONDS": "120", "TASK_LEASE_SECONDS": "150", "TASK_MAX_ATTEMPTS": "3", "TASK_RETRY_BASE_MS": "5000", "TASK_POLL_MS": "1000", "AI_ENABLED": "true", "EMBEDDING_TIMEOUT_SECONDS": "30", "APP_ENV": "test", "APP_VERSION": "test", "HTTP_ADDR": ":8080", "DB_ADDR": "db:3306", "DB_NAME": "test", "DB_USER": "test", "DB_PASSWORD": strings.Repeat("x", 24), "DB_TIMEOUT_MS": "1500", "SHUTDOWN_TIMEOUT_SECONDS": "10"}
}
func TestNormalConfigS02(t *testing.T) {
	m := valid()
	if _, e := Load(func(k string) string { return m[k] }); e != nil {
		t.Fatal(e)
	}
}
func TestFailureMissingEveryKeyS02(t *testing.T) {
	for k := range valid() {
		t.Run(k, func(t *testing.T) {
			m := valid()
			delete(m, k)
			if _, e := Load(func(k string) string { return m[k] }); e == nil {
				t.Fatal("missing config accepted")
			}
		})
	}
}
func TestBoundaryTimeoutS03(t *testing.T) {
	for _, v := range []struct {
		s  string
		ok bool
	}{{"99", false}, {"100", true}, {"10000", true}, {"10001", false}, {"not-a-number", false}} {
		t.Run(v.s, func(t *testing.T) {
			m := valid()
			m["DB_TIMEOUT_MS"] = v.s
			_, e := Load(func(k string) string { return m[k] })
			if (e == nil) != v.ok {
				t.Fatal("timeout boundary")
			}
		})
	}
}
func TestFailureConfigDoesNotLeak(t *testing.T) {
	m := valid()
	m["DB_PASSWORD"] = "sensitive"
	_, e := Load(func(k string) string { return m[k] })
	if e == nil || strings.Contains(e.Error(), m["DB_PASSWORD"]) {
		t.Fatal("config leak")
	}
}

func TestAIDisabledDoesNotRequireProviderConfiguration(t *testing.T) {
	m := valid()
	m["AI_ENABLED"] = "false"
	for _, key := range []string{"EMBEDDING_BASE_URL", "EMBEDDING_API_KEY", "EMBEDDING_MODEL", "EMBEDDING_DIM", "EMBEDDING_TIMEOUT_SECONDS", "CHAT_BASE_URL", "CHAT_API_KEY", "CHAT_MODEL"} {
		delete(m, key)
	}
	if _, e := Load(func(k string) string { return m[k] }); e != nil {
		t.Fatal(e)
	}
}
