package search

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const sampleDDG = `{
  "AbstractURL": "https://example.com/abstract",
  "AbstractText": "Sample abstract",
  "RelatedTopics": [
    {"FirstURL": "https://example.com/a", "Text": "Topic A - example.com"},
    {"Topics": [
      {"FirstURL": "https://example.com/b", "Text": "Topic B"}
    ]}
  ]
}`

func TestHTTPProviderSearch(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("q") != "ai infra" {
			t.Fatalf("query %q", r.URL.Query().Get("q"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(sampleDDG))
	}))
	defer srv.Close()

	p := NewHTTPProvider(HTTPOptions{
		Timeout:  DefaultHTTPOptions().Timeout,
		Endpoint: srv.URL,
	})
	results, err := p.Search("ai infra")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) < 3 {
		t.Fatalf("expected at least 3 results, got %d", len(results))
	}
	if results[0].URL == "" {
		t.Fatal("expected URLs")
	}
}

func TestHTTPProviderRateLimit(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer srv.Close()

	p := NewHTTPProvider(HTTPOptions{Endpoint: srv.URL, Timeout: DefaultHTTPOptions().Timeout})
	_, err := p.Search("test")
	if err == nil || !strings.Contains(err.Error(), "429") {
		t.Fatalf("expected 429 error, got %v", err)
	}
}

func TestHTTPProviderNoResults(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"RelatedTopics":[]}`))
	}))
	defer srv.Close()

	p := NewHTTPProvider(HTTPOptions{Endpoint: srv.URL, Timeout: DefaultHTTPOptions().Timeout})
	_, err := p.Search("empty")
	if err == nil {
		t.Fatal("expected error for empty results")
	}
}

func TestNewSearchProviderFactory(t *testing.T) {
	if NewSearchProvider("mock").Name() != "mock" {
		t.Fatal("expected mock")
	}
	if NewSearchProvider("http").Name() != "http" {
		t.Fatal("expected http")
	}
}
