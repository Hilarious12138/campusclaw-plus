package repository

import (
	"campusclaw-plus/internal/domain"
	"context"
	"database/sql"
	"errors"
)

func scanResource(row interface{ Scan(...any) error }) (domain.Resource, error) {
	var v domain.Resource
	e := row.Scan(&v.ID, &v.ClassID, &v.CreatedBy, &v.Kind, &v.Revision, &v.Body, &v.CreatedAt, &v.UpdatedAt)
	if errors.Is(e, sql.ErrNoRows) {
		e = domain.ErrNotFound
	}
	return v, e
}

const resourceFields = "id,class_id,created_by,kind,revision,body,created_at,updated_at"

func (m *MySQL) Resources(ctx context.Context, u domain.Identity, kind string) ([]domain.Resource, error) {
	rows, e := m.db.QueryContext(ctx, "SELECT "+resourceFields+" FROM resources WHERE class_id=? AND kind=? ORDER BY created_at DESC,id DESC LIMIT 100", u.ClassID, kind)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []domain.Resource{}
	for rows.Next() {
		v, e := scanResource(rows)
		if e != nil {
			return nil, e
		}
		out = append(out, v)
	}
	return out, rows.Err()
}
func (m *MySQL) Resource(ctx context.Context, u domain.Identity, kind, id string) (domain.Resource, error) {
	return scanResource(m.db.QueryRowContext(ctx, "SELECT "+resourceFields+" FROM resources WHERE class_id=? AND kind=? AND id=?", u.ClassID, kind, id))
}
func (m *MySQL) SaveResource(ctx context.Context, u domain.Identity, v domain.Resource, event domain.Audit) (domain.Resource, error) {
	if u.Role != "teacher" {
		return v, domain.ErrForbidden
	}
	tx, e := m.db.BeginTx(ctx, nil)
	if e != nil {
		return v, e
	}
	defer tx.Rollback()
	if v.ID == "" {
		v.ID = NewID()
		_, e = tx.ExecContext(ctx, "INSERT INTO resources(id,class_id,created_by,kind,revision,body,created_at,updated_at) VALUES(?,?,?,?,1,?,UTC_TIMESTAMP(6),UTC_TIMESTAMP(6))", v.ID, u.ClassID, u.ID, v.Kind, []byte(v.Body))
	} else {
		var result sql.Result
		result, e = tx.ExecContext(ctx, "UPDATE resources SET body=?,revision=revision+1,updated_at=UTC_TIMESTAMP(6) WHERE id=? AND class_id=? AND kind=? AND revision=?", []byte(v.Body), v.ID, u.ClassID, v.Kind, v.Revision)
		if e == nil {
			n, _ := result.RowsAffected()
			if n == 0 {
				var n int
				if e = tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM resources WHERE id=? AND class_id=? AND kind=?", v.ID, u.ClassID, v.Kind).Scan(&n); e == nil {
					if n == 0 {
						e = domain.ErrNotFound
					} else {
						e = domain.ErrConflict
					}
				}
			}
		}
	}
	if e != nil {
		return v, e
	}
	event.ObjectType = v.Kind
	event.ObjectID = v.ID
	if e = appendAudit(ctx, tx, event); e != nil {
		return v, e
	}
	if e = tx.Commit(); e != nil {
		return v, e
	}
	return m.Resource(ctx, u, v.Kind, v.ID)
}
func (m *MySQL) DeleteResource(ctx context.Context, u domain.Identity, kind, id string, event domain.Audit) error {
	if u.Role != "teacher" {
		return domain.ErrForbidden
	}
	tx, e := m.db.BeginTx(ctx, nil)
	if e != nil {
		return e
	}
	defer tx.Rollback()
	res, e := tx.ExecContext(ctx, "DELETE FROM resources WHERE id=? AND class_id=? AND kind=?", id, u.ClassID, kind)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrNotFound
	}
	event.ObjectType = kind
	event.ObjectID = id
	if e = appendAudit(ctx, tx, event); e != nil {
		return e
	}
	return tx.Commit()
}
