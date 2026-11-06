package repository

import (
	"campusclaw-plus/internal/domain"
	"context"
	"database/sql"
	"errors"
	"strings"
)

const chunkFields = "c.id,c.material_id,m.title,c.chunk_index,c.start_char,c.end_char,c.body,c.version"

func chunkScope(ids []string) (string, []any) {
	if ids == nil {
		return "", nil
	}
	if len(ids) == 0 {
		return " AND 1=0", nil
	}
	args := []any{}
	for _, v := range ids {
		args = append(args, v)
	}
	return " AND c.material_id IN (" + strings.TrimSuffix(strings.Repeat("?,", len(ids)), ",") + ")", args
}
func scanChunk(row interface{ Scan(...any) error }) (domain.Chunk, error) {
	var c domain.Chunk
	e := row.Scan(&c.ID, &c.MaterialID, &c.Title, &c.Index, &c.Start, &c.End, &c.Text, &c.Version)
	if errors.Is(e, sql.ErrNoRows) {
		e = domain.ErrNotFound
	}
	return c, e
}
func (m *MySQL) Keyword(ctx context.Context, u domain.Identity, query string, ids []string, limit int) ([]domain.Chunk, error) {
	scope, args := chunkScope(ids)
	params := []any{u.ClassID, query}
	params = append(params, args...)
	params = append(params, query, limit)
	rows, e := m.db.QueryContext(ctx, "SELECT "+chunkFields+" FROM knowledge_chunks c JOIN materials m ON m.id=c.material_id AND m.class_id=c.class_id AND m.active_version=c.version WHERE c.class_id=? AND MATCH(c.body) AGAINST (? IN NATURAL LANGUAGE MODE)>0"+scope+" ORDER BY MATCH(c.body) AGAINST (? IN NATURAL LANGUAGE MODE) DESC,c.id LIMIT ?", params...)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []domain.Chunk{}
	for rows.Next() {
		v, e := scanChunk(rows)
		if e != nil {
			return nil, e
		}
		out = append(out, v)
	}
	return out, rows.Err()
}
func (m *MySQL) ChunksByID(ctx context.Context, u domain.Identity, ids, materials []string) ([]domain.Chunk, error) {
	if len(ids) == 0 {
		return []domain.Chunk{}, nil
	}
	params := []any{u.ClassID}
	for _, v := range ids {
		params = append(params, v)
	}
	scope, args := chunkScope(materials)
	params = append(params, args...)
	rows, e := m.db.QueryContext(ctx, "SELECT "+chunkFields+" FROM knowledge_chunks c JOIN materials m ON m.id=c.material_id AND m.class_id=c.class_id AND m.active_version=c.version WHERE c.class_id=? AND c.id IN ("+strings.TrimSuffix(strings.Repeat("?,", len(ids)), ",")+")"+scope, params...)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []domain.Chunk{}
	for rows.Next() {
		v, e := scanChunk(rows)
		if e != nil {
			return nil, e
		}
		out = append(out, v)
	}
	return out, rows.Err()
}
func (m *MySQL) Chunk(ctx context.Context, u domain.Identity, id string) (domain.Chunk, error) {
	return scanChunk(m.db.QueryRowContext(ctx, "SELECT "+chunkFields+" FROM knowledge_chunks c JOIN materials m ON m.id=c.material_id AND m.class_id=c.class_id AND m.active_version=c.version WHERE c.class_id=? AND c.id=?", u.ClassID, id))
}
