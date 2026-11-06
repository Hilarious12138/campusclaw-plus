package service

import (
	"campusclaw-plus/internal/domain"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"
)

type Auth struct {
	Store          domain.AuthStore
	Credentials    domain.Credentials
	Now            func() time.Time
	FingerprintKey []byte
	DummyHash      string
}

func (a Auth) Login(ctx context.Context, username, password, source string) (string, domain.Identity, error) {
	var empty domain.Identity
	if len(username) > 100 || len(password) > 1024 || username == "" || password == "" {
		return "", empty, domain.ErrUnauthenticated
	}
	username = strings.TrimSpace(username)
	h := hmac.New(sha256.New, a.FingerprintKey)
	h.Write([]byte(strings.ToLower(username) + "\x00" + source))
	key := hex.EncodeToString(h.Sum(nil))
	allowed, e := a.Store.LoginAllowed(ctx, key, a.Now())
	if e != nil {
		return "", empty, e
	}
	if !allowed {
		return "", empty, domain.ErrLimited
	}
	u, e := a.Store.LoginUser(ctx, username)
	if e != nil && !errors.Is(e, domain.ErrNotFound) {
		return "", empty, e
	}
	hash := u.Hash
	if hash == "" {
		hash = a.DummyHash
	}
	match := a.Credentials.Match(hash, password) && u.ID != ""
	if e = a.Store.LoginFailure(ctx, key, a.Now(), match); e != nil {
		return "", empty, e
	}
	if !match {
		return "", empty, domain.ErrUnauthenticated
	}
	value, e := a.Credentials.Issue(u.ID)
	return value, u.Identity, e
}
func (a Auth) Identify(ctx context.Context, value string) (domain.Identity, domain.Claims, error) {
	c, e := a.Credentials.Verify(value)
	if e != nil {
		return domain.Identity{}, c, e
	}
	u, e := a.Store.Identity(ctx, c)
	return u, c, e
}
func (a Auth) Revoke(ctx context.Context, c domain.Claims, event domain.Audit) error {
	return a.Store.Revoke(ctx, c, event)
}
func (a Auth) Audit(ctx context.Context, event domain.Audit) error {
	return a.Store.AppendAudit(ctx, event)
}
func (a Auth) Events(ctx context.Context, u domain.Identity, f domain.AuditFilter) ([]domain.Audit, error) {
	if u.Role != "teacher" {
		return nil, domain.ErrForbidden
	}
	return a.Store.AuditEvents(ctx, u, f)
}
