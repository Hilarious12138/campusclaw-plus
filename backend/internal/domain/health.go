package domain

import (
	"context"
	"time"
)

type Dependency interface{ Ping(context.Context) error }
type Health struct {
	Status       string            `json:"status"`
	Version      string            `json:"version"`
	Time         time.Time         `json:"time"`
	Dependencies map[string]string `json:"dependencies"`
}
