package service

import (
	"campusclaw-plus/internal/domain"
	"context"
	"regexp"
)

var serverName = regexp.MustCompile(`^[a-z][a-z0-9-]{0,31}$`)

type MCP struct {
	Store  domain.MCPStore
	Remote domain.MCPRemote
}

func (s MCP) Register(ctx context.Context, u domain.Identity, v domain.MCPServer, event domain.Audit) (domain.MCPServer, error) {
	if u.Role != "teacher" {
		return v, domain.ErrForbidden
	}
	if !serverName.MatchString(v.Name) || !s.Remote.Allowed(v.BaseURL) {
		return v, domain.ErrInvalid
	}
	return s.Store.RegisterMCP(ctx, u, v, event)
}
func (s MCP) Connect(ctx context.Context, u domain.Identity, id string, connect bool, event domain.Audit) ([]domain.MCPTool, error) {
	if u.Role != "teacher" {
		return nil, domain.ErrForbidden
	}
	v, e := s.Store.MCPServer(ctx, u, id)
	if e != nil {
		return nil, e
	}
	tools := []domain.MCPTool{}
	if connect {
		tools, e = s.Remote.Discover(ctx, v.BaseURL)
		if e != nil {
			return nil, e
		}
	}
	e = s.Store.MCPConnected(ctx, u, v, connect, tools, event)
	return tools, e
}
