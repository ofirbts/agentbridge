package mcp

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPClientListTools(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatal(err)
		}
		if req.Method != "tools/list" {
			t.Fatalf("method %q", req.Method)
		}
		_ = json.NewEncoder(w).Encode(Response{
			JSONRPC: JSONRPCVersion,
			ID:      req.ID,
			Result:  json.RawMessage(`{"tools":[{"name":"web_search"},{"name":"fetch_url"}]}`),
		})
	}))
	defer srv.Close()

	c := NewHTTPClient(srv.URL)
	tools, err := c.ListTools(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(tools) != 2 || tools[0] != "web_search" {
		t.Fatalf("unexpected tools: %v", tools)
	}
}

func TestHTTPClientCallTool(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatal(err)
		}
		if req.Method != "tools/call" {
			t.Fatalf("method %q", req.Method)
		}
		_ = json.NewEncoder(w).Encode(Response{
			JSONRPC: JSONRPCVersion,
			ID:      req.ID,
			Result:  json.RawMessage(`{"content":[{"type":"text","text":"ok"}]}`),
		})
	}))
	defer srv.Close()

	c := NewHTTPClient(srv.URL)
	res, err := c.CallTool(context.Background(), "web_search", map[string]any{"query": "ai"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Content) != 1 || res.Content[0].Text != "ok" {
		t.Fatalf("unexpected result: %#v", res)
	}
}

func TestHTTPClientRPCError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(Response{
			JSONRPC: JSONRPCVersion,
			ID:      1,
			Error:   &Error{Code: -32600, Message: "invalid request"},
		})
	}))
	defer srv.Close()

	_, err := NewHTTPClient(srv.URL).ListTools(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestNewClientHTTPDefault(t *testing.T) {
	c, err := NewClient("https://example.com/mcp", "")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := c.(*HTTPClient); !ok {
		t.Fatalf("expected HTTPClient, got %T", c)
	}
}

func TestNewClientStub(t *testing.T) {
	c, err := NewClient("http://localhost:8080/mcp", "stub")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := c.(*StubClient); !ok {
		t.Fatalf("expected StubClient, got %T", c)
	}
}
