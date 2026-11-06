package domain

import (
	"context"
	"time"
)

type MCPTool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Schema      any    `json:"inputSchema"`
}
type MCPServer struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	BaseURL   string    `json:"base_url"`
	Connected bool      `json:"connected"`
	Revision  int       `json:"revision"`
	Tools     []MCPTool `json:"tools"`
	ClassID   string    `json:"-"`
}
type MCPStore interface {
	MCPServers(context.Context, Identity) ([]MCPServer, error)
	MCPServer(context.Context, Identity, string) (MCPServer, error)
	RegisterMCP(context.Context, Identity, MCPServer, Audit) (MCPServer, error)
	MCPConnected(context.Context, Identity, MCPServer, bool, []MCPTool, Audit) error
	ToolServer(context.Context, Identity, string) (MCPServer, string, error)
}
type MCPRemote interface {
	Discover(context.Context, string) ([]MCPTool, error)
	Call(context.Context, string, string, string) (string, error)
	Allowed(string) bool
}
type ToolCall struct {
	Ref    string `json:"tool_ref"`
	Status string `json:"status"`
}
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type ChatModel interface {
	Stream(context.Context, []ChatMessage, func(string) error) error
}
type Conversation struct {
	ID        string    `json:"id"`
	AgentID   string    `json:"agent_id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}
type Turn struct {
	ID        string     `json:"id"`
	Question  string     `json:"question"`
	Answer    string     `json:"answer"`
	Status    string     `json:"status"`
	SkillName string     `json:"skill_name"`
	Citations []Chunk    `json:"citations"`
	ToolCalls []ToolCall `json:"tool_calls"`
	CreatedAt time.Time  `json:"created_at"`
}
type ConversationStore interface {
	CreateConversation(context.Context, Identity, string, Audit) (Conversation, error)
	Conversations(context.Context, Identity) ([]Conversation, error)
	Conversation(context.Context, Identity, string) (Conversation, []Turn, error)
	BeginTurn(context.Context, Identity, string, string, string, Audit) (Turn, error)
	FinishTurn(context.Context, Identity, string, string, string, []Chunk, []ToolCall, []ChatMessage, string, Audit) error
}
