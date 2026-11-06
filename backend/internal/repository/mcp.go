package repository

import (
	"campusclaw-plus/internal/domain"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
)

func scanMCP(row interface{ Scan(...any) error }) (domain.MCPServer, error) {
	var s domain.MCPServer
	var raw []byte
	e := row.Scan(&s.ID, &s.ClassID, &s.Name, &s.BaseURL, &s.Connected, &s.Revision, &raw)
	if errors.Is(e, sql.ErrNoRows) {
		e = domain.ErrNotFound
	}
	if e == nil {
		e = json.Unmarshal(raw, &s.Tools)
	}
	return s, e
}
func (m *MySQL) MCPServers(ctx context.Context, u domain.Identity) ([]domain.MCPServer, error) {
	if u.Role != "teacher" {
		return nil, domain.ErrForbidden
	}
	rows, e := m.db.QueryContext(ctx, "SELECT id,class_id,name,base_url,connected,revision,tools FROM mcp_servers WHERE class_id=? ORDER BY name,id LIMIT 100", u.ClassID)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []domain.MCPServer{}
	for rows.Next() {
		v, e := scanMCP(rows)
		if e != nil {
			return nil, e
		}
		out = append(out, v)
	}
	return out, rows.Err()
}
func (m *MySQL) MCPServer(ctx context.Context, u domain.Identity, id string) (domain.MCPServer, error) {
	return scanMCP(m.db.QueryRowContext(ctx, "SELECT id,class_id,name,base_url,connected,revision,tools FROM mcp_servers WHERE class_id=? AND id=?", u.ClassID, id))
}
func (m *MySQL) RegisterMCP(ctx context.Context, u domain.Identity, s domain.MCPServer, event domain.Audit) (domain.MCPServer, error) {
	if u.Role != "teacher" {
		return s, domain.ErrForbidden
	}
	tx, e := m.db.BeginTx(ctx, nil)
	if e != nil {
		return s, e
	}
	defer tx.Rollback()
	s.ID = NewID()
	_, e = tx.ExecContext(ctx, "INSERT INTO mcp_servers(id,class_id,name,base_url,tools) VALUES(?,?,?,?,'[]')", s.ID, u.ClassID, s.Name, s.BaseURL)
	if e != nil {
		return s, domain.ErrConflict
	}
	event.ObjectID = s.ID
	event.ObjectType = "mcp_server"
	if e = appendAudit(ctx, tx, event); e != nil {
		return s, e
	}
	if e = tx.Commit(); e != nil {
		return s, e
	}
	return m.MCPServer(ctx, u, s.ID)
}
func (m *MySQL) MCPConnected(ctx context.Context, u domain.Identity, s domain.MCPServer, connected bool, tools []domain.MCPTool, event domain.Audit) error {
	if u.Role != "teacher" {
		return domain.ErrForbidden
	}
	tx, e := m.db.BeginTx(ctx, nil)
	if e != nil {
		return e
	}
	defer tx.Rollback()
	if tools == nil {
		tools = []domain.MCPTool{}
	}
	raw, e := json.Marshal(tools)
	if e != nil {
		return e
	}
	res, e := tx.ExecContext(ctx, "UPDATE mcp_servers SET connected=?,tools=?,revision=revision+1 WHERE id=? AND class_id=? AND revision=?", connected, raw, s.ID, u.ClassID, s.Revision)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrConflict
	}
	event.ObjectID = s.ID
	event.ObjectType = "mcp_server"
	if e = appendAudit(ctx, tx, event); e != nil {
		return e
	}
	return tx.Commit()
}
func toolServer(ctx context.Context, q interface {
	QueryRowContext(context.Context, string, ...any) *sql.Row
}, u domain.Identity, ref string, lock bool) (domain.MCPServer, string, error) {
	p := strings.Split(ref, "@")
	if len(p) != 2 {
		return domain.MCPServer{}, "", domain.ErrForbidden
	}
	query := "SELECT id,class_id,name,base_url,connected,revision,tools FROM mcp_servers WHERE class_id=? AND name=?"
	if lock {
		query += " FOR UPDATE"
	}
	s, e := scanMCP(q.QueryRowContext(ctx, query, u.ClassID, p[1]))
	if e != nil {
		return s, p[0], domain.ErrForbidden
	}
	if !s.Connected {
		return s, p[0], domain.ErrForbidden
	}
	for _, t := range s.Tools {
		if t.Name == p[0] {
			return s, p[0], nil
		}
	}
	return s, p[0], domain.ErrForbidden
}
func (m *MySQL) ToolServer(ctx context.Context, u domain.Identity, ref string) (domain.MCPServer, string, error) {
	return toolServer(ctx, m.db, u, ref, false)
}
