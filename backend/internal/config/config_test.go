package config

import (
	"strings"
	"testing"
)

func valid() map[string]string {
	return map[string]string{"APP_ENV": "test", "APP_VERSION": "test", "HTTP_ADDR": ":8080", "DB_ADDR": "db:3306", "DB_NAME": "test", "DB_USER": "test", "DB_PASSWORD": strings.Repeat("x", 24), "DB_TIMEOUT_MS": "1500", "SHUTDOWN_TIMEOUT_SECONDS": "10"}
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
