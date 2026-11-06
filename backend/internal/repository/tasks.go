package repository

import (
	"campusclaw-plus/internal/domain"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"strconv"
)

func materialHash(title, body string) string {
	h := sha256.Sum256([]byte(title + "\x00" + body))
	return hex.EncodeToString(h[:])
}
func stableID(value string) string {
	h := sha256.Sum256([]byte(value))
	return hex.EncodeToString(h[:16])
}
func (m *MySQL) CreateMaterial(ctx context.Context, u domain.Identity, title, body, key string, event domain.Audit) (domain.Material, domain.Task, error) {
	var mat domain.Material
	var task domain.Task
	if u.Role != "teacher" {
		return mat, task, domain.ErrForbidden
	}
	hash := materialHash(title, body)
	mid := stableID(u.ClassID + hash + "index-v1")
	tid := stableID(mid + "index-v1")
	tx, e := m.db.BeginTx(ctx, nil)
	if e != nil {
		return mat, task, e
	}
	defer tx.Rollback()
	// Serializes per-class quota and idempotency checks, not external work.
	var cid string
	if e = tx.QueryRowContext(ctx, "SELECT id FROM classes WHERE id=? FOR UPDATE", u.ClassID).Scan(&cid); e != nil {
		return mat, task, e
	}
	var oldHash, oldID string
	e = tx.QueryRowContext(ctx, "SELECT input_hash,task_id FROM task_requests WHERE class_id=? AND actor_id=? AND idempotency_key=?", u.ClassID, u.ID, key).Scan(&oldHash, &oldID)
	if e == nil {
		if oldHash != hash {
			return mat, task, domain.ErrConflict
		}
		tid = oldID
	} else if !errors.Is(e, sql.ErrNoRows) {
		return mat, task, e
	} else {
		var count, exists int
		if e = tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM tasks WHERE id=? AND class_id=?", tid, u.ClassID).Scan(&exists); e != nil {
			return mat, task, e
		}
		if exists == 0 {
			if e = tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM tasks WHERE class_id=? AND status IN ('pending','running')", u.ClassID).Scan(&count); e != nil {
				return mat, task, e
			}
			if count >= 100 {
				return mat, task, domain.ErrLimited
			}
			if _, e = tx.ExecContext(ctx, "INSERT INTO materials VALUES(?,?,?,?,?,?,?,'',UTC_TIMESTAMP(6),UTC_TIMESTAMP(6))", mid, u.ClassID, u.ID, title, body, hash, tid); e != nil {
				return mat, task, e
			}
			if _, e = tx.ExecContext(ctx, "INSERT INTO tasks(id,class_id,created_by,material_id,input_hash,status,stage,next_run_at,created_at,updated_at) VALUES(?,?,?,?,?,'pending','queued',UTC_TIMESTAMP(6),UTC_TIMESTAMP(6),UTC_TIMESTAMP(6))", tid, u.ClassID, u.ID, mid, hash); e != nil {
				return mat, task, e
			}
			if e = enqueue(ctx, tx, tid, "work", 0, 0); e != nil {
				return mat, task, e
			}
		}
		if _, e = tx.ExecContext(ctx, "INSERT INTO task_requests VALUES(?,?,?,?,?)", u.ClassID, u.ID, key, hash, tid); e != nil {
			return mat, task, e
		}
	}
	event.ObjectType = "material"
	event.ObjectID = mid
	event.TaskID = tid
	if e = appendAudit(ctx, tx, event); e != nil {
		return mat, task, e
	}
	if e = tx.Commit(); e != nil {
		return mat, task, e
	}
	task, e = m.Task(ctx, u, tid)
	if e != nil {
		return mat, task, e
	}
	mat, e = m.Material(ctx, u, task.MaterialID)
	return mat, task, e
}
func enqueue(ctx context.Context, q executor, tid, kind string, generation, delay int) error {
	_, e := q.ExecContext(ctx, "INSERT INTO task_outbox(id,task_id,kind,generation,available_at) VALUES(?,?,?,?,TIMESTAMPADD(MICROSECOND,?,UTC_TIMESTAMP(6))) ON DUPLICATE KEY UPDATE id=id", NewID(), tid, kind, generation, delay*1000)
	return e
}

const matFields = "m.id,m.class_id,m.created_by,m.title,m.content,m.input_hash,m.task_id,m.active_version,m.created_at,t.status"

func scanMaterial(row interface{ Scan(...any) error }) (domain.Material, error) {
	var v domain.Material
	e := row.Scan(&v.ID, &v.ClassID, &v.CreatedBy, &v.Title, &v.Content, &v.InputHash, &v.TaskID, &v.ActiveVersion, &v.CreatedAt, &v.Status)
	if errors.Is(e, sql.ErrNoRows) {
		e = domain.ErrNotFound
	}
	return v, e
}
func (m *MySQL) Materials(ctx context.Context, u domain.Identity, limit, offset int) ([]domain.Material, error) {
	rows, e := m.db.QueryContext(ctx, "SELECT "+matFields+" FROM materials m JOIN tasks t ON t.id=m.task_id WHERE m.class_id=? ORDER BY m.created_at DESC,m.id DESC LIMIT ? OFFSET ?", u.ClassID, limit, offset)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []domain.Material{}
	for rows.Next() {
		v, e := scanMaterial(rows)
		if e != nil {
			return nil, e
		}
		out = append(out, v)
	}
	return out, rows.Err()
}
func (m *MySQL) Material(ctx context.Context, u domain.Identity, id string) (domain.Material, error) {
	return scanMaterial(m.db.QueryRowContext(ctx, "SELECT "+matFields+" FROM materials m JOIN tasks t ON t.id=m.task_id WHERE m.id=? AND m.class_id=?", id, u.ClassID))
}

const taskFields = "id,material_id,class_id,created_by,status,stage,attempts,delivery_count,fence,lease_owner,error_code,next_run_at,started_at,finished_at,input_hash"

func scanTask(row interface{ Scan(...any) error }) (domain.Task, error) {
	var v domain.Task
	e := row.Scan(&v.ID, &v.MaterialID, &v.ClassID, &v.CreatedBy, &v.Status, &v.Stage, &v.Attempts, &v.Deliveries, &v.Fence, &v.Owner, &v.ErrorCode, &v.NextRunAt, &v.StartedAt, &v.FinishedAt, &v.Hash)
	if errors.Is(e, sql.ErrNoRows) {
		e = domain.ErrNotFound
	}
	return v, e
}
func (m *MySQL) Task(ctx context.Context, u domain.Identity, id string) (domain.Task, error) {
	return scanTask(m.db.QueryRowContext(ctx, "SELECT "+taskFields+" FROM tasks WHERE id=? AND class_id=?", id, u.ClassID))
}
func (m *MySQL) ClaimOutbox(ctx context.Context, owner string) (domain.Outbox, error) {
	v := domain.Outbox{Owner: owner}
	tx, e := m.db.BeginTx(ctx, nil)
	if e != nil {
		return v, e
	}
	defer tx.Rollback()
	e = tx.QueryRowContext(ctx, "SELECT id,task_id,kind FROM task_outbox WHERE published_at IS NULL AND available_at<=UTC_TIMESTAMP(6) AND (lease_until IS NULL OR lease_until<UTC_TIMESTAMP(6)) ORDER BY available_at,id LIMIT 1 FOR UPDATE SKIP LOCKED").Scan(&v.ID, &v.TaskID, &v.Kind)
	if errors.Is(e, sql.ErrNoRows) {
		return v, domain.ErrNotFound
	}
	if e != nil {
		return v, e
	}
	_, e = tx.ExecContext(ctx, "UPDATE task_outbox SET lease_owner=?,lease_until=TIMESTAMPADD(SECOND,10,UTC_TIMESTAMP(6)),publish_attempts=publish_attempts+1 WHERE id=?", owner, v.ID)
	if e != nil {
		return v, e
	}
	return v, tx.Commit()
}
func (m *MySQL) Published(ctx context.Context, v domain.Outbox) error {
	_, e := m.db.ExecContext(ctx, "UPDATE task_outbox SET published_at=UTC_TIMESTAMP(6) WHERE id=? AND lease_owner=?", v.ID, v.Owner)
	return e
}
func (m *MySQL) ClaimTask(ctx context.Context, id, owner string) (domain.Task, domain.Material, error) {
	var mat domain.Material
	tx, e := m.db.BeginTx(ctx, nil)
	if e != nil {
		return domain.Task{}, mat, e
	}
	defer tx.Rollback()
	t, e := scanTask(tx.QueryRowContext(ctx, "SELECT "+taskFields+" FROM tasks WHERE id=? FOR UPDATE", id))
	if e != nil {
		return t, mat, e
	}
	if _, e = tx.ExecContext(ctx, "UPDATE tasks SET delivery_count=delivery_count+1 WHERE id=?", id); e != nil {
		return t, mat, e
	}
	if t.Status != "pending" {
		if e = tx.Commit(); e != nil {
			return t, mat, e
		}
		return t, mat, domain.ErrConflict
	}
	var ready bool
	if e = tx.QueryRowContext(ctx, "SELECT next_run_at<=UTC_TIMESTAMP(6) FROM tasks WHERE id=?", id).Scan(&ready); e != nil {
		return t, mat, e
	}
	if !ready {
		if e = tx.Commit(); e != nil {
			return t, mat, e
		}
		return t, mat, domain.ErrConflict
	}
	mat, e = scanMaterial(tx.QueryRowContext(ctx, "SELECT "+matFields+" FROM materials m JOIN tasks t ON t.id=m.task_id WHERE m.id=? AND m.class_id=?", t.MaterialID, t.ClassID))
	if e != nil {
		return t, mat, e
	}
	_, e = tx.ExecContext(ctx, "UPDATE tasks SET status='running',stage='embedding',attempts=attempts+1,fence=fence+1,lease_owner=?,lease_until=TIMESTAMPADD(SECOND,?,UTC_TIMESTAMP(6)),started_at=UTC_TIMESTAMP(6),updated_at=UTC_TIMESTAMP(6) WHERE id=?", owner, m.TaskLeaseSeconds, id)
	if e != nil {
		return t, mat, e
	}
	t.Attempts++
	t.Fence++
	t.Owner = owner
	t.Status = "running"
	if _, e = tx.ExecContext(ctx, "INSERT INTO task_attempts(task_id,attempt,started_at) VALUES(?,?,UTC_TIMESTAMP(6))", id, t.Attempts); e != nil {
		return t, mat, e
	}
	if e = tx.Commit(); e != nil {
		return t, mat, e
	}
	var permitted int
	if e = m.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users WHERE id=? AND class_id=? AND active=true AND role='teacher'", t.CreatedBy, t.ClassID).Scan(&permitted); e != nil {
		return t, mat, e
	}
	if permitted != 1 {
		return t, mat, domain.ErrForbidden
	}
	return t, mat, nil
}
func (m *MySQL) RenewTask(ctx context.Context, t domain.Task) error {
	r, e := m.db.ExecContext(ctx, "UPDATE tasks SET lease_until=TIMESTAMPADD(SECOND,?,UTC_TIMESTAMP(6)) WHERE id=? AND status='running' AND lease_owner=? AND fence=? AND lease_until>UTC_TIMESTAMP(6)", m.TaskLeaseSeconds, t.ID, t.Owner, t.Fence)
	if e != nil {
		return e
	}
	n, _ := r.RowsAffected()
	if n == 0 {
		return domain.ErrConflict
	}
	return nil
}
func guardTask(ctx context.Context, tx *sql.Tx, t domain.Task) error {
	var n int
	e := tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM tasks WHERE id=? AND status='running' AND lease_owner=? AND fence=? AND lease_until>UTC_TIMESTAMP(6) FOR UPDATE", t.ID, t.Owner, t.Fence).Scan(&n)
	if e != nil {
		return e
	}
	if n != 1 {
		return domain.ErrConflict
	}
	return nil
}
func (m *MySQL) CompleteTask(ctx context.Context, t domain.Task, chunks []domain.Chunk) error {
	tx, e := m.db.BeginTx(ctx, nil)
	if e != nil {
		return e
	}
	defer tx.Rollback()
	if e = guardTask(ctx, tx, t); e != nil {
		return e
	}
	var n int
	if e = tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM users WHERE id=? AND class_id=? AND role='teacher' AND active=true", t.CreatedBy, t.ClassID).Scan(&n); e != nil {
		return e
	}
	if n != 1 {
		return domain.ErrForbidden
	}
	for _, c := range chunks {
		_, e = tx.ExecContext(ctx, "INSERT INTO knowledge_chunks(id,class_id,material_id,version,chunk_index,start_char,end_char,body) VALUES(?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE body=VALUES(body)", c.ID, t.ClassID, t.MaterialID, t.Hash, c.Index, c.Start, c.End, c.Text)
		if e != nil {
			return e
		}
	}
	if _, e = tx.ExecContext(ctx, "INSERT INTO task_results VALUES(?,?,?,UTC_TIMESTAMP(6))", t.ID, t.Hash, len(chunks)); e != nil {
		return e
	}
	if _, e = tx.ExecContext(ctx, "UPDATE materials SET active_version=?,updated_at=UTC_TIMESTAMP(6) WHERE id=? AND class_id=?", t.Hash, t.MaterialID, t.ClassID); e != nil {
		return e
	}
	if _, e = tx.ExecContext(ctx, "UPDATE tasks SET status='succeeded',stage='publishing',finished_at=UTC_TIMESTAMP(6),updated_at=UTC_TIMESTAMP(6),lease_until=NULL,error_code='' WHERE id=?", t.ID); e != nil {
		return e
	}
	if _, e = tx.ExecContext(ctx, "UPDATE task_attempts SET finished_at=UTC_TIMESTAMP(6) WHERE task_id=? AND attempt=?", t.ID, t.Attempts); e != nil {
		return e
	}
	if e = appendAudit(ctx, tx, domain.Audit{ActorID: t.CreatedBy, ClassID: t.ClassID, Operation: "material.index", ObjectType: "material", ObjectID: t.MaterialID, HTTPStatus: 200, Outcome: "OK", TaskID: t.ID, RequestID: t.ID, Usage: "unknown"}); e != nil {
		return e
	}
	return tx.Commit()
}
func (m *MySQL) FailTask(ctx context.Context, t domain.Task, code string, retry bool) error {
	tx, e := m.db.BeginTx(ctx, nil)
	if e != nil {
		return e
	}
	defer tx.Rollback()
	if e = guardTask(ctx, tx, t); e != nil {
		return e
	}
	state, kind, delay := "failed", "dead", 0
	if retry && t.Attempts < m.TaskMaxAttempts {
		state, kind = "pending", "work"
		delay = m.TaskRetryBaseMS * (1 << (t.Attempts - 1))
		j, _ := strconv.ParseInt(NewID()[:3], 16, 32)
		delay += int(j % 1001)
	}
	if _, e = tx.ExecContext(ctx, "UPDATE tasks SET status=?,stage='queued',error_code=?,next_run_at=TIMESTAMPADD(MICROSECOND,?,UTC_TIMESTAMP(6)),lease_until=NULL,finished_at=IF(?='failed',UTC_TIMESTAMP(6),NULL),updated_at=UTC_TIMESTAMP(6) WHERE id=?", state, code, delay*1000, state, t.ID); e != nil {
		return e
	}
	if _, e = tx.ExecContext(ctx, "UPDATE task_attempts SET finished_at=UTC_TIMESTAMP(6),error_code=? WHERE task_id=? AND attempt=?", code, t.ID, t.Attempts); e != nil {
		return e
	}
	if e = enqueue(ctx, tx, t.ID, kind, t.Fence, delay); e != nil {
		return e
	}
	return tx.Commit()
}
func (m *MySQL) RecoverTasks(ctx context.Context) error {
	rows, e := m.db.QueryContext(ctx, "SELECT "+taskFields+" FROM tasks WHERE status='running' AND lease_until<UTC_TIMESTAMP(6) LIMIT 100")
	if e != nil {
		return e
	}
	var tasks []domain.Task
	for rows.Next() {
		t, e := scanTask(rows)
		if e != nil {
			rows.Close()
			return e
		}
		tasks = append(tasks, t)
	}
	e = rows.Err()
	rows.Close()
	if e != nil {
		return e
	}
	for _, t := range tasks {
		tx, e := m.db.BeginTx(ctx, nil)
		if e != nil {
			return e
		}
		state, kind := "pending", "work"
		if t.Attempts >= m.TaskMaxAttempts {
			state, kind = "failed", "dead"
		}
		res, e := tx.ExecContext(ctx, "UPDATE tasks SET status=?,stage='queued',fence=fence+1,error_code='WORKER_LOST',lease_until=NULL,next_run_at=UTC_TIMESTAMP(6),finished_at=IF(?='failed',UTC_TIMESTAMP(6),NULL) WHERE id=? AND fence=? AND status='running' AND lease_until<UTC_TIMESTAMP(6)", state, state, t.ID, t.Fence)
		if e == nil {
			n, _ := res.RowsAffected()
			if n == 1 {
				e = enqueue(ctx, tx, t.ID, kind, t.Fence+1, 0)
			}
		}
		if e == nil {
			_, e = tx.ExecContext(ctx, "UPDATE task_attempts SET finished_at=UTC_TIMESTAMP(6),error_code='WORKER_LOST' WHERE task_id=? AND attempt=? AND finished_at IS NULL", t.ID, t.Attempts)
		}
		if e != nil {
			tx.Rollback()
			return e
		}
		if e = tx.Commit(); e != nil {
			return e
		}
	}
	return nil
}
