package httpapi

import (
	"campusclaw-plus/internal/domain"
	"campusclaw-plus/internal/service"
	"context"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type authFake struct {
	events  []domain.Audit
	fail    bool
	revoked bool
}

func (f *authFake) LoginUser(context.Context, string) (domain.LoginUser, error) {
	return domain.LoginUser{}, domain.ErrNotFound
}
func (f *authFake) Identity(_ context.Context, c domain.Claims) (domain.Identity, error) {
	if f.revoked {
		return domain.Identity{}, domain.ErrUnauthenticated
	}
	return domain.Identity{ID: c.Subject, Role: "student", ClassID: strings.Repeat("a", 32)}, nil
}
func (f *authFake) Revoke(_ context.Context, _ domain.Claims, v domain.Audit) error {
	if f.fail {
		return errors.New("audit offline")
	}
	f.revoked = true
	f.events = append(f.events, v)
	return nil
}
func (f *authFake) AppendAudit(_ context.Context, v domain.Audit) error {
	if f.fail {
		return errors.New("audit offline")
	}
	f.events = append(f.events, v)
	return nil
}
func (f *authFake) AuditEvents(context.Context, domain.Identity, domain.AuditFilter) ([]domain.Audit, error) {
	return f.events, nil
}
func (f *authFake) LoginAllowed(context.Context, string, time.Time) (bool, error) { return true, nil }
func (f *authFake) LoginFailure(context.Context, string, time.Time, bool) error   { return nil }

type credsFake struct{}

func (credsFake) Hash(string) (string, error)  { return "hash", nil }
func (credsFake) Match(string, string) bool    { return false }
func (credsFake) Issue(string) (string, error) { return "opaque", nil }
func (credsFake) Verify(s string) (domain.Claims, error) {
	if s != "valid" {
		return domain.Claims{}, domain.ErrUnauthenticated
	}
	return domain.Claims{Subject: strings.Repeat("b", 32)}, nil
}
func TestAuthorizationAndAuditFailures(t *testing.T) {
	for _, c := range []struct {
		path, method, token, origin, body string
		status                            int
	}{
		{"/api/me", "GET", "", "", "", 401},
		{"/api/me?token=valid", "GET", "", "", "", 401},
		{"/api/me?class_id=other", "GET", "valid", "", "", 400},
		{"/api/me", "GET", "valid", "", "", 200},
		{"/api/audit-events", "GET", "valid", "", "", 403},
		{"/api/skills", "POST", "valid", "", "{}", 403},
		{"/api/auth/logout", "POST", "valid", "https://evil.invalid", "{}", 403},
	} {
		f := &authFake{}
		p := Platform{Auth: service.Auth{Store: f, Credentials: credsFake{}}, Origin: "http://localhost"}
		r := httptest.NewRequest(c.method, c.path, strings.NewReader(c.body))
		if c.token != "" {
			r.Header.Set("Authorization", "Bearer "+c.token)
		}
		r.Header.Set("Origin", c.origin)
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		p.serve(w, r, "test")
		if w.Code != c.status || len(f.events) != 1 || f.events[0].HTTPStatus != c.status {
			t.Fatal(c.path, w.Code, w.Body.String(), f.events)
		}
		if c.status >= 400 && strings.Contains(w.Body.String(), `"data"`) {
			t.Fatal("rejection leaked data")
		}
	}
	f := &authFake{fail: true}
	p := Platform{Auth: service.Auth{Store: f, Credentials: credsFake{}}, Origin: "http://localhost"}
	w := httptest.NewRecorder()
	p.serve(w, httptest.NewRequest("GET", "/api/me", nil), "test")
	if w.Code != 503 {
		t.Fatal("audit failure hidden")
	}
}
