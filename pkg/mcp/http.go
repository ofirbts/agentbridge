package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
)

type HTTPClient struct {
	endpoint string
	client   *http.Client
	nextID   atomic.Int64
}

func NewHTTPClient(endpoint string) *HTTPClient {
	return &HTTPClient{
		endpoint: endpoint,
		client:   &http.Client{},
	}
}

func (c *HTTPClient) Endpoint() string {
	return c.endpoint
}

func (c *HTTPClient) Close() error {
	return nil
}

func (c *HTTPClient) ListTools(ctx context.Context) ([]string, error) {
	raw, err := c.call(ctx, "tools/list", map[string]any{})
	if err != nil {
		return nil, err
	}
	var out struct {
		Tools []struct {
			Name string `json:"name"`
		} `json:"tools"`
	}
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, fmt.Errorf("mcp tools/list decode: %w", err)
	}
	names := make([]string, 0, len(out.Tools))
	for _, t := range out.Tools {
		if t.Name != "" {
			names = append(names, t.Name)
		}
	}
	return names, nil
}

func (c *HTTPClient) CallTool(ctx context.Context, name string, args map[string]any) (*ToolResult, error) {
	raw, err := c.call(ctx, "tools/call", ToolCallParams{Name: name, Arguments: args})
	if err != nil {
		return nil, err
	}
	var result ToolResult
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("mcp tools/call decode: %w", err)
	}
	return &result, nil
}

func (c *HTTPClient) call(ctx context.Context, method string, params any) (json.RawMessage, error) {
	var rawParams json.RawMessage
	if params != nil {
		b, err := json.Marshal(params)
		if err != nil {
			return nil, err
		}
		rawParams = b
	}
	reqBody := Request{
		JSONRPC: JSONRPCVersion,
		ID:      c.nextID.Add(1),
		Method:  method,
		Params:  rawParams,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("mcp http %s: %w", method, err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("mcp http %s: status %d", method, resp.StatusCode)
	}
	var rpc Response
	if err := json.Unmarshal(respBody, &rpc); err != nil {
		return nil, fmt.Errorf("mcp http %s: %w", method, err)
	}
	if rpc.Error != nil {
		return nil, fmt.Errorf("mcp %s: %s", method, rpc.Error.Message)
	}
	return rpc.Result, nil
}
