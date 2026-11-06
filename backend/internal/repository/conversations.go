package repository

import (
	"campusclaw-plus/internal/domain"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
)

func (m *MySQL) CreateConversation(ctx context.Context, u domain.Identity, agent string, event domain.Audit) (domain.Conversation, error) {
	v := domain.Conversation{ID: NewID(), AgentID: agent, Title: "新会话"}
	tx, e := m.db.BeginTx(ctx, nil)
	if e != nil {
		return v, e
	}
	defer tx.Rollback()
	var n int
	if e = tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM resources WHERE id=? AND class_id=? AND kind='agents' AND JSON_EXTRACT(body,'$.enabled')=true", agent, u.ClassID).Scan(&n); e != nil {
		return v, e
	}
	if n != 1 {
		return v, domain.ErrNotFound
	}
	if _, e = tx.ExecContext(ctx, "INSERT INTO conversations(id,class_id,created_by,agent_id,title,created_at) VALUES(?,?,?,?,?,UTC_TIMESTAMP(6))", v.ID, u.ClassID, u.ID, agent, v.Title); e != nil {
		return v, e
	}
	event.ObjectType = "conversation"
	event.ObjectID = v.ID
	if e = appendAudit(ctx, tx, event); e != nil {
		return v, e
	}
	if e = tx.Commit(); e != nil {
		return v, e
	}
	v, _, e = m.Conversation(ctx, u, v.ID)
	return v, e
}
func (m *MySQL) Conversations(ctx context.Context, u domain.Identity) ([]domain.Conversation, error) {
	rows, e := m.db.QueryContext(ctx, "SELECT id,agent_id,title,created_at FROM conversations WHERE class_id=? AND created_by=? ORDER BY created_at DESC,id DESC LIMIT 100", u.ClassID, u.ID)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []domain.Conversation{}
	for rows.Next() {
		var v domain.Conversation
		if e = rows.Scan(&v.ID, &v.AgentID, &v.Title, &v.CreatedAt); e != nil {
			return nil, e
		}
		out = append(out, v)
	}
	return out, rows.Err()
}
func (m *MySQL) Conversation(ctx context.Context, u domain.Identity, id string) (domain.Conversation, []domain.Turn, error) {
	var v domain.Conversation
	e := m.db.QueryRowContext(ctx, "SELECT id,agent_id,title,created_at FROM conversations WHERE id=? AND class_id=? AND created_by=?", id, u.ClassID, u.ID).Scan(&v.ID, &v.AgentID, &v.Title, &v.CreatedAt)
	if errors.Is(e, sql.ErrNoRows) {
		e = domain.ErrNotFound
	}
	if e != nil {
		return v, nil, e
	}
	rows, e := m.db.QueryContext(ctx, "SELECT id,question,answer,status,skill_name,citations,tool_calls,created_at FROM conversation_turns WHERE conversation_id=? ORDER BY created_at,id LIMIT 100", id)
	if e != nil {
		return v, nil, e
	}
	defer rows.Close()
	turns := []domain.Turn{}
	for rows.Next() {
		var t domain.Turn
		var c, tools []byte
		if e = rows.Scan(&t.ID, &t.Question, &t.Answer, &t.Status, &t.SkillName, &c, &tools, &t.CreatedAt); e != nil {
			return v, nil, e
		}
		if e = json.Unmarshal(c, &t.Citations); e != nil {
			return v, nil, e
		}
		if e = json.Unmarshal(tools, &t.ToolCalls); e != nil {
			return v, nil, e
		}
		turns = append(turns, t)
	}
	return v, turns, rows.Err()
}
func (m *MySQL) BeginTurn(ctx context.Context, u domain.Identity, conv, question, skill string, event domain.Audit) (domain.Turn, error) {
	v := domain.Turn{ID: NewID(), Question: question, Status: "running", SkillName: skill}
	tx, e := m.db.BeginTx(ctx, nil)
	if e != nil {
		return v, e
	}
	defer tx.Rollback()
	var busy string
	var active bool
	e = tx.QueryRowContext(ctx, "SELECT busy_turn,COALESCE(lease_until>UTC_TIMESTAMP(6),false) FROM conversations WHERE id=? AND class_id=? AND created_by=? FOR UPDATE", conv, u.ClassID, u.ID).Scan(&busy, &active)
	if errors.Is(e, sql.ErrNoRows) {
		e = domain.ErrNotFound
	}
	if e != nil {
		return v, e
	}
	if busy != "" && active {
		return v, domain.ErrConflict
	}
	if busy != "" {
		if _, e = tx.ExecContext(ctx, "UPDATE conversation_turns SET status='failed' WHERE id=? AND status='running'", busy); e != nil {
			return v, e
		}
	}
	if _, e = tx.ExecContext(ctx, "INSERT INTO conversation_turns(id,conversation_id,question,answer,status,skill_name,citations,tool_calls,snapshot,created_at) VALUES(?,?,?,'','running',?,'[]','[]','[]',UTC_TIMESTAMP(6))", v.ID, conv, question, skill); e != nil {
		return v, e
	}
	title := []rune(question)
	if len(title) > 80 {
		title = title[:80]
	}
	if _, e = tx.ExecContext(ctx, "UPDATE conversations SET busy_turn=?,lease_until=TIMESTAMPADD(SECOND,90,UTC_TIMESTAMP(6)),title=? WHERE id=?", v.ID, string(title), conv); e != nil {
		return v, e
	}
	event.ObjectID = conv
	event.ObjectType = "conversation"
	event.TurnID = v.ID
	event.Outcome = "pending"
	if e = appendAudit(ctx, tx, event); e != nil {
		return v, e
	}
	return v, tx.Commit()
}
func (m *MySQL) FinishTurn(ctx context.Context, u domain.Identity, conv, turn, answer string, citations []domain.Chunk, tools []domain.ToolCall, snapshot []domain.ChatMessage, skill string, event domain.Audit) error {
	tx, e := m.db.BeginTx(ctx, nil)
	if e != nil {
		return e
	}
	defer tx.Rollback()
	var active int
	if e = tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM conversations WHERE id=? AND class_id=? AND created_by=? AND busy_turn=? AND lease_until>UTC_TIMESTAMP(6) FOR UPDATE", conv, u.ClassID, u.ID, turn).Scan(&active); e != nil {
		return e
	}
	if active != 1 {
		return domain.ErrConflict
	}
	if citations == nil {
		citations = []domain.Chunk{}
	}
	if tools == nil {
		tools = []domain.ToolCall{}
	}
	c, _ := json.Marshal(citations)
	t, _ := json.Marshal(tools)
	s, _ := json.Marshal(snapshot)
	state := "succeeded"
	if event.HTTPStatus >= 400 {
		state = "failed"
	}
	if _, e = tx.ExecContext(ctx, "UPDATE conversation_turns SET answer=?,status=?,citations=?,tool_calls=?,snapshot=?,skill_name=? WHERE id=? AND conversation_id=?", answer, state, c, t, s, skill, turn, conv); e != nil {
		return e
	}
	if _, e = tx.ExecContext(ctx, "UPDATE conversations SET busy_turn='',lease_until=NULL WHERE id=? AND busy_turn=?", conv, turn); e != nil {
		return e
	}
	event.ObjectID = conv
	event.ObjectType = "conversation"
	event.TurnID = turn
	if e = appendAudit(ctx, tx, event); e != nil {
		return e
	}
	return tx.Commit()
}
