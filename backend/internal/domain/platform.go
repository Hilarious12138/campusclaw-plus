package domain

import (
	"context"
	"encoding/json"
	"errors"
	"time"
)

var ErrUnauthenticated = errors.New("unauthenticated")
var ErrForbidden = errors.New("forbidden")
var ErrNotFound = errors.New("not found")
var ErrInvalid = errors.New("invalid argument")
var ErrConflict = errors.New("conflict")
var ErrLimited = errors.New("rate limited")

type Identity struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	ClassID  string `json:"class_id"`
}
type Claims struct {
	Subject, ID string
	ExpiresAt   time.Time
}
type Credentials interface {
	Hash(string) (string, error)
	Match(string, string) bool
	Issue(string) (string, error)
	Verify(string) (Claims, error)
}
type LoginUser struct {
	Identity
	Hash string
}
type Audit struct {
	ID         string    `json:"id"`
	ActorID    string    `json:"actor_id,omitempty"`
	ClassID    string    `json:"class_id,omitempty"`
	OccurredAt time.Time `json:"occurred_at"`
	Operation  string    `json:"operation"`
	ObjectType string    `json:"object_type"`
	ObjectID   string    `json:"object_id"`
	Outcome    string    `json:"outcome"`
	HTTPStatus int       `json:"http_status"`
	SkillID    string    `json:"skill_id,omitempty"`
	Usage      any       `json:"usage"`
	DurationMS int64     `json:"duration_ms"`
	RequestID  string    `json:"request_id"`
	TaskID     string    `json:"task_id,omitempty"`
	TurnID     string    `json:"turn_id,omitempty"`
}
type AuditFilter struct {
	ActorID, Operation    string
	From, Until           time.Time
	Status, Limit, Offset int
}
type AuthStore interface {
	LoginUser(context.Context, string) (LoginUser, error)
	Identity(context.Context, Claims) (Identity, error)
	Revoke(context.Context, Claims, Audit) error
	AppendAudit(context.Context, Audit) error
	AuditEvents(context.Context, Identity, AuditFilter) ([]Audit, error)
	LoginAllowed(context.Context, string, time.Time) (bool, error)
	LoginFailure(context.Context, string, time.Time, bool) error
}
type Resource struct {
	ID        string          `json:"id"`
	ClassID   string          `json:"class_id"`
	CreatedBy string          `json:"created_by"`
	Kind      string          `json:"kind"`
	Revision  int             `json:"revision"`
	Body      json.RawMessage `json:"body"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type ResourceStore interface {
	Resources(context.Context, Identity, string) ([]Resource, error)
	Resource(context.Context, Identity, string, string) (Resource, error)
	SaveResource(context.Context, Identity, Resource, Audit) (Resource, error)
	DeleteResource(context.Context, Identity, string, string, Audit) error
}

// MarshalJSON keeps absent identities and skills explicit in seven-column audit records.
func (a Audit) MarshalJSON() ([]byte, error) {
	type plain Audit
	b, e := json.Marshal(plain(a))
	if e != nil {
		return nil, e
	}
	var m map[string]any
	if e = json.Unmarshal(b, &m); e != nil {
		return nil, e
	}
	if a.ActorID == "" {
		m["actor_id"] = nil
	}
	if a.ClassID == "" {
		m["class_id"] = nil
	}
	if a.SkillID == "" {
		m["skill_id"] = nil
	}
	if a.Usage == nil {
		m["usage"] = "unknown"
	}
	return json.Marshal(m)
}
