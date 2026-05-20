package mcp

import (
	"fmt"
	"strings"
)

func NewClient(endpoint, transport string) (Client, error) {
	if endpoint == "" {
		return nil, fmt.Errorf("mcp: endpoint not configured")
	}
	switch strings.ToLower(transport) {
	case "stub":
		return NewStubClient(endpoint), nil
	case "http", "":
		return NewHTTPClient(endpoint), nil
	default:
		return nil, fmt.Errorf("mcp: unknown transport %q", transport)
	}
}

func DefaultTransport(endpoint string) string {
	if strings.HasPrefix(endpoint, "http://") || strings.HasPrefix(endpoint, "https://") {
		return "http"
	}
	return "stub"
}
