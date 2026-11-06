package repository

import (
	"campusclaw-plus/internal/domain"
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

//go:embed platform-schema.sql
var platformSchema string

func (m *MySQL) Migrate(ctx context.Context, credentials domain.Credentials, seedPassword string) error {
	for _, q := range strings.Split(platformSchema, ";") {
		if strings.TrimSpace(q) != "" {
			if _, e := m.db.ExecContext(ctx, q); e != nil {
				return errors.New("platform migration unavailable")
			}
		}
	}
	var n int
	if e := m.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version=?", "platform-auth-v1").Scan(&n); e != nil {
		return e
	}
	if n > 0 {
		return nil
	}
	hash, e := credentials.Hash(seedPassword)
	if e != nil {
		return e
	}
	tx, e := m.db.BeginTx(ctx, nil)
	if e != nil {
		return e
	}
	defer tx.Rollback()
	for _, c := range []string{"a", "b"} {
		cid := strings.Repeat(c, 32)
		if _, e = tx.ExecContext(ctx, "INSERT INTO classes(id,name) VALUES(?,?) ON DUPLICATE KEY UPDATE id=id", cid, "班级 "+strings.ToUpper(c)); e != nil {
			return e
		}
		for _, role := range []string{"teacher", "student"} {
			if _, e = tx.ExecContext(ctx, "INSERT INTO users(id,username,password_hash,role,class_id,created_at) VALUES(?,?,?,?,?,UTC_TIMESTAMP(6)) ON DUPLICATE KEY UPDATE id=id", NewID(), role+"_"+c, hash, role, cid); e != nil {
				return e
			}
		}
	}
	if _, e = tx.ExecContext(ctx, "INSERT INTO schema_migrations VALUES(?,UTC_TIMESTAMP(6)) ON DUPLICATE KEY UPDATE version=version", "platform-auth-v1"); e != nil {
		return e
	}
	return tx.Commit()
}
func (m *MySQL) LoginUser(ctx context.Context, name string) (domain.LoginUser, error) {
	var u domain.LoginUser
	e := m.db.QueryRowContext(ctx, "SELECT id,username,role,class_id,password_hash FROM users WHERE username=? AND active=true", name).Scan(&u.ID, &u.Username, &u.Role, &u.ClassID, &u.Hash)
	if errors.Is(e, sql.ErrNoRows) {
		e = domain.ErrNotFound
	}
	return u, e
}
func (m *MySQL) Identity(ctx context.Context, c domain.Claims) (domain.Identity, error) {
	var u domain.Identity
	e := m.db.QueryRowContext(ctx, "SELECT id,username,role,class_id FROM users WHERE id=? AND active=true AND NOT EXISTS(SELECT 1 FROM revoked_tokens WHERE jti=? AND expires_at>UTC_TIMESTAMP(6))", c.Subject, c.ID).Scan(&u.ID, &u.Username, &u.Role, &u.ClassID)
	if errors.Is(e, sql.ErrNoRows) {
		e = domain.ErrUnauthenticated
	}
	return u, e
}

type executor interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
}

func appendAudit(ctx context.Context, q executor, v domain.Audit) error {
	if v.ID == "" {
		v.ID = NewID()
	}
	if v.OccurredAt.IsZero() {
		v.OccurredAt = time.Now().UTC()
	}
	if v.Usage == nil {
		v.Usage = "unknown"
	}
	for _, s := range []*string{&v.Operation, &v.ObjectType, &v.ObjectID, &v.Outcome, &v.SkillID, &v.RequestID, &v.TaskID, &v.TurnID} {
		if len(*s) > 256 {
			*s = (*s)[:256]
		}
	}
	raw, e := json.Marshal(v)
	if e != nil || len(raw) > 8192 {
		return errors.New("invalid audit event")
	}
	_, e = q.ExecContext(ctx, "INSERT INTO audit_events(id,class_id,actor_id,occurred_at,operation,http_status,event) VALUES(?,NULLIF(?,''),NULLIF(?,''),?,?,?,?)", v.ID, v.ClassID, v.ActorID, v.OccurredAt.UTC(), v.Operation, v.HTTPStatus, raw)
	return e
}
func (m *MySQL) AppendAudit(ctx context.Context, v domain.Audit) error {
	return appendAudit(ctx, m.db, v)
}
func (m *MySQL) Revoke(ctx context.Context, c domain.Claims, v domain.Audit) error {
	tx, e := m.db.BeginTx(ctx, nil)
	if e != nil {
		return e
	}
	defer tx.Rollback()
	if _, e = tx.ExecContext(ctx, "INSERT INTO revoked_tokens(jti,expires_at) VALUES(?,?) ON DUPLICATE KEY UPDATE jti=jti", c.ID, c.ExpiresAt.UTC()); e != nil {
		return e
	}
	if e = appendAudit(ctx, tx, v); e != nil {
		return e
	}
	return tx.Commit()
}
func (m *MySQL) AuditEvents(ctx context.Context, u domain.Identity, f domain.AuditFilter) ([]domain.Audit, error) {
	if u.Role != "teacher" {
		return nil, domain.ErrForbidden
	}
	q := "SELECT event FROM audit_events WHERE class_id=?"
	args := []any{u.ClassID}
	if f.ActorID != "" {
		q += " AND actor_id=?"
		args = append(args, f.ActorID)
	}
	if f.Operation != "" {
		q += " AND operation=?"
		args = append(args, f.Operation)
	}
	if f.Status != 0 {
		q += " AND http_status=?"
		args = append(args, f.Status)
	}
	if !f.From.IsZero() {
		q += " AND occurred_at>=?"
		args = append(args, f.From.UTC())
	}
	if !f.Until.IsZero() {
		q += " AND occurred_at<=?"
		args = append(args, f.Until.UTC())
	}
	q += " ORDER BY occurred_at DESC,id DESC LIMIT ? OFFSET ?"
	args = append(args, f.Limit, f.Offset)
	rows, e := m.db.QueryContext(ctx, q, args...)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []domain.Audit{}
	for rows.Next() {
		var raw []byte
		var v domain.Audit
		if e = rows.Scan(&raw); e != nil {
			return nil, e
		}
		if e = json.Unmarshal(raw, &v); e != nil {
			return nil, e
		}
		out = append(out, v)
	}
	return out, rows.Err()
}
func (m *MySQL) LoginAllowed(ctx context.Context, key string, now time.Time) (bool, error) {
	var n int
	var start time.Time
	e := m.db.QueryRowContext(ctx, "SELECT failures,window_start FROM login_failures WHERE fingerprint=?", key).Scan(&n, &start)
	if errors.Is(e, sql.ErrNoRows) {
		return true, nil
	}
	return n < 5 || now.Sub(start) >= 15*time.Minute, e
}
func (m *MySQL) LoginFailure(ctx context.Context, key string, now time.Time, success bool) error {
	if success {
		_, e := m.db.ExecContext(ctx, "DELETE FROM login_failures WHERE fingerprint=?", key)
		return e
	}
	_, e := m.db.ExecContext(ctx, "INSERT INTO login_failures VALUES(?,1,?) ON DUPLICATE KEY UPDATE failures=IF(window_start<=?,1,failures+1),window_start=IF(window_start<=?,VALUES(window_start),window_start)", key, now.UTC(), now.Add(-15*time.Minute).UTC(), now.Add(-15*time.Minute).UTC())
	return e
}
