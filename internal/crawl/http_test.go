package crawl

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPCrawlerFetchSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("<html>ok</html>"))
	}))
	defer srv.Close()

	c := NewHTTPCrawler(DefaultHTTPOptions())
	body, err := c.Fetch(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "<html>ok</html>" {
		t.Fatalf("unexpected body: %s", body)
	}
}

func TestHTTPCrawlerRateLimit(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer srv.Close()

	c := NewHTTPCrawler(DefaultHTTPOptions())
	_, err := c.Fetch(srv.URL)
	if err == nil {
		t.Fatal("expected rate limit error")
	}
}

func TestNewCrawlerFactory(t *testing.T) {
	if NewCrawler("mock").Name() != "mock_http" {
		t.Fatal("expected mock")
	}
	if NewCrawler("http").Name() != "http" {
		t.Fatal("expected http")
	}
}
