package search

import (
	"fmt"
	"sort"
	"strings"
)

type MockProvider struct {
	Fail    bool
	Results []Result
}

func NewMockProvider() *MockProvider {
	return &MockProvider{
		Results: []Result{
			{Title: "AgentBridge", URL: "https://example.com/agentbridge", Snippet: "execution layer for web agents"},
			{Title: "AI Infra Startups", URL: "https://example.com/ai-infra", Snippet: "reliable web workflows"},
		},
	}
}

func (m *MockProvider) Name() string {
	return "mock"
}

func (m *MockProvider) Search(query string) ([]Result, error) {
	if m.Fail {
		return nil, fmt.Errorf("mock search failure for query %q", query)
	}
	out := make([]Result, len(m.Results))
	copy(out, m.Results)
	for i := range out {
		out[i].Snippet = strings.TrimSpace(out[i].Snippet + " | " + query)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].URL < out[j].URL
	})
	return out, nil
}
