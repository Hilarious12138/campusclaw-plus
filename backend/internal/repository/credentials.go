package repository

import (
	"campusclaw-plus/internal/domain"
	"crypto/rand"
	"encoding/hex"
	"github.com/golang-jwt/jwt/v5"
	"strings"
	"time"
)

type Credentials struct {
	Secret []byte
	Now    func() time.Time
}

func NewID() string {
	var b [16]byte
	if _, e := rand.Read(b[:]); e != nil {
		panic("random source unavailable")
	}
	return hex.EncodeToString(b[:])
}
func (c Credentials) Hash(p string) (string, error) { return passwordHash(p) }
func (c Credentials) Match(h, p string) bool        { return passwordMatches(h, p) }
func (c Credentials) Issue(id string) (string, error) {
	now := c.Now().UTC()
	claims := jwt.RegisteredClaims{Issuer: "campusclaw-plus", Audience: jwt.ClaimStrings{"campusclaw-plus-api"}, Subject: id, ID: NewID(), IssuedAt: jwt.NewNumericDate(now), NotBefore: jwt.NewNumericDate(now), ExpiresAt: jwt.NewNumericDate(now.Add(30 * time.Minute))}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(c.Secret)
}
func (c Credentials) Verify(value string) (domain.Claims, error) {
	var out domain.Claims
	if len(value) > 4096 || strings.TrimSpace(value) != value {
		return out, domain.ErrUnauthenticated
	}
	v := &jwt.RegisteredClaims{}
	t, e := jwt.ParseWithClaims(value, v, func(*jwt.Token) (any, error) { return c.Secret, nil }, jwt.WithValidMethods([]string{"HS256"}), jwt.WithIssuer("campusclaw-plus"), jwt.WithAudience("campusclaw-plus-api"), jwt.WithExpirationRequired(), jwt.WithIssuedAt(), jwt.WithTimeFunc(c.Now))
	if e != nil || t == nil || !t.Valid || len(v.Subject) != 32 || len(v.ID) != 32 || v.IssuedAt == nil || v.NotBefore == nil || v.ExpiresAt == nil || v.ExpiresAt.Sub(v.IssuedAt.Time) > 30*time.Minute || !v.ExpiresAt.After(v.IssuedAt.Time) {
		return out, domain.ErrUnauthenticated
	}
	if _, e = hex.DecodeString(v.Subject); e != nil {
		return out, domain.ErrUnauthenticated
	}
	if _, e = hex.DecodeString(v.ID); e != nil {
		return out, domain.ErrUnauthenticated
	}
	return domain.Claims{Subject: v.Subject, ID: v.ID, ExpiresAt: v.ExpiresAt.Time}, nil
}
