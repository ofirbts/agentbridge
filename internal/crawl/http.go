package crawl

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type HTTPOptions struct {
	Timeout   time.Duration
	UserAgent string
}

func DefaultHTTPOptions() HTTPOptions {
	return HTTPOptions{
		Timeout:   15 * time.Second,
		UserAgent: "AgentBridge/0.2 (+https://github.com/yourname/agentbridge)",
	}
}

type HTTPCrawler struct {
	client *http.Client
	opts   HTTPOptions
}

func NewHTTPCrawler(opts HTTPOptions) *HTTPCrawler {
	if opts.Timeout <= 0 {
		opts = DefaultHTTPOptions()
	}
	if opts.UserAgent == "" {
		opts.UserAgent = DefaultHTTPOptions().UserAgent
	}
	return &HTTPCrawler{
		client: &http.Client{Timeout: opts.Timeout},
		opts:   opts,
	}
}

func (h *HTTPCrawler) Name() string {
	return "http"
}

func (h *HTTPCrawler) Fetch(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("http crawl: build request: %w", err)
	}
	req.Header.Set("User-Agent", h.opts.UserAgent)
	req.Header.Set("Accept", "text/html,application/json")

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http crawl: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	if err != nil {
		return nil, fmt.Errorf("http crawl: read body: %w", err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return body, nil
	case http.StatusTooManyRequests:
		return nil, fmt.Errorf("http crawl: rate limited (429) for %s", url)
	case http.StatusRequestTimeout, http.StatusGatewayTimeout:
		return nil, fmt.Errorf("http crawl: timeout status (%d) for %s", resp.StatusCode, url)
	default:
		if resp.StatusCode >= 500 {
			return nil, fmt.Errorf("http crawl: server error (%d) for %s", resp.StatusCode, url)
		}
		return nil, fmt.Errorf("http crawl: unexpected status (%d) for %s", resp.StatusCode, url)
	}
}
