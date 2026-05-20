package crawl

import (
	"fmt"
)

type MockCrawler struct {
	Fail    bool
	Content []byte
}

func NewMockCrawler() *MockCrawler {
	return &MockCrawler{
		Content: []byte("<html><body><h1>AgentBridge</h1><p>mock page</p></body></html>"),
	}
}

func (m *MockCrawler) Name() string {
	return "mock"
}

func (m *MockCrawler) Fetch(url string) ([]byte, error) {
	if m.Fail {
		return nil, fmt.Errorf("mock fetch failed for %s", url)
	}
	body := make([]byte, len(m.Content))
	copy(body, m.Content)
	return body, nil
}
