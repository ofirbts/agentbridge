package mcp

import "context"

type Client interface {
	Endpoint() string
	CallTool(ctx context.Context, name string, args map[string]any) (*ToolResult, error)
	ListTools(ctx context.Context) ([]string, error)
	Close() error
}
