package service

import (
	"campusclaw-plus/internal/domain"
	"context"
	"encoding/json"
	"strings"
)

type Resources struct{ Store domain.ResourceStore }

func (s Resources) List(ctx context.Context, u domain.Identity, kind string) ([]domain.Resource, error) {
	return s.Store.Resources(ctx, u, kind)
}
func (s Resources) Read(ctx context.Context, u domain.Identity, kind, id string) (domain.Resource, error) {
	return s.Store.Resource(ctx, u, kind, id)
}
func (s Resources) Save(ctx context.Context, u domain.Identity, v domain.Resource, event domain.Audit) (domain.Resource, error) {
	if u.Role != "teacher" {
		return v, domain.ErrForbidden
	}
	if len(v.Body) > 65536 || !json.Valid(v.Body) {
		return v, domain.ErrInvalid
	}
	var b map[string]json.RawMessage
	if json.Unmarshal(v.Body, &b) != nil {
		return v, domain.ErrInvalid
	}
	allowed := map[string]map[string]bool{
		"agents":  {"name": true, "prompt": true, "skill_ids": true, "material_ids": true, "enabled": true},
		"skills":  {"name": true, "conditions": true, "instructions": true, "tools": true, "enabled": true},
		"prompts": {"name": true, "content": true},
	}
	fields, ok := allowed[v.Kind]
	if !ok {
		return v, domain.ErrInvalid
	}
	for k := range b {
		if !fields[k] {
			return v, domain.ErrInvalid
		}
	}
	var name string
	json.Unmarshal(b["name"], &name)
	if strings.TrimSpace(name) == "" || len([]rune(name)) > 100 {
		return v, domain.ErrInvalid
	}
	// Typed structural validation prevents clients smuggling malformed role/tool fields.
	for _, k := range []string{"prompt", "content", "instructions"} {
		if raw, ok := b[k]; ok {
			var text string
			if json.Unmarshal(raw, &text) != nil || len(text) > 32000 {
				return v, domain.ErrInvalid
			}
		}
	}
	for _, k := range []string{"skill_ids", "material_ids", "conditions", "tools"} {
		if raw, ok := b[k]; ok {
			var values []string
			if json.Unmarshal(raw, &values) != nil || len(values) > 100 {
				return v, domain.ErrInvalid
			}
		}
	}
	if raw, ok := b["enabled"]; ok {
		var enabled bool
		if json.Unmarshal(raw, &enabled) != nil {
			return v, domain.ErrInvalid
		}
	}
	return s.Store.SaveResource(ctx, u, v, event)
}
func (s Resources) Delete(ctx context.Context, u domain.Identity, kind, id string, event domain.Audit) error {
	if u.Role != "teacher" {
		return domain.ErrForbidden
	}
	return s.Store.DeleteResource(ctx, u, kind, id, event)
}
