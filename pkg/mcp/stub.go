package mcp

import (
	"context"
	"fmt"
)

type StubClient struct {
	endpoint string
	tools    []string
}

func NewStubClient(endpoint string) *StubClient {
	return &StubClient{
		endpoint: endpoint,
		tools:    []string{"web_search", "fetch_url", "extract_text"},
	}
}

func (s *StubClient) Endpoint() string {
	return s.endpoint
}

func (s *StubClient) ListTools(ctx context.Context) ([]string, error) {
	if s.endpoint == "" {
		return nil, fmt.Errorf("mcp stub: endpoint not configured")
	}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return append([]string(nil), s.tools...), nil
	}
}

func (s *StubClient) CallTool(ctx context.Context, name string, args map[string]any) (*ToolResult, error) {
	if s.endpoint == "" {
		return nil, fmt.Errorf("mcp stub: endpoint not configured")
	}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	query, _ := args["query"].(string)
	return &ToolResult{
		Content: []ContentBlock{{
			Type: "text",
			Text: fmt.Sprintf("stub mcp tool %q at %s query=%q", name, s.endpoint, query),
		}},
	}, nil
}

func (s *StubClient) Close() error {
	return nil
}
