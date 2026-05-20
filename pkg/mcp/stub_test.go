package mcp

import (
	"context"
	"testing"
)

func TestStubClientListTools(t *testing.T) {
	c := NewStubClient("http://localhost:8080/mcp")
	tools, err := c.ListTools(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(tools) != 3 {
		t.Fatalf("expected 3 tools, got %d", len(tools))
	}
}

func TestStubClientCallTool(t *testing.T) {
	c := NewStubClient("http://localhost:8080/mcp")
	res, err := c.CallTool(context.Background(), "web_search", map[string]any{"query": "ai infra"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Content) != 1 || res.Content[0].Text == "" {
		t.Fatal("expected stub content")
	}
}

func TestStubClientNoEndpoint(t *testing.T) {
	c := NewStubClient("")
	_, err := c.ListTools(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}
