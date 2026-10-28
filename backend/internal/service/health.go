package service

import (
	"campusclaw-plus/internal/domain"
	"context"
	"time"
)

type Health struct {
	DB      domain.Dependency
	Version string
	Timeout time.Duration
	Now     func() time.Time
}

func (s Health) Check(ctx context.Context) (domain.Health, bool) {
	ctx, cancel := context.WithTimeout(ctx, s.Timeout)
	defer cancel()
	h := domain.Health{Status: "ready", Version: s.Version, Time: s.Now().UTC(), Dependencies: map[string]string{"mysql": "up"}}
	if s.DB.Ping(ctx) != nil {
		h.Status = "not_ready"
		h.Dependencies["mysql"] = "down"
		return h, false
	}
	return h, true
}
