package search

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

type HTTPOptions struct {
	Timeout   time.Duration
	UserAgent string
	Endpoint  string
}

func DefaultHTTPOptions() HTTPOptions {
	return HTTPOptions{
		Timeout:   10 * time.Second,
		UserAgent: "AgentBridge/0.1 (+https://github.com/ofirbts/agentbridge)",
		Endpoint:  "https://api.duckduckgo.com/",
	}
}

type HTTPProvider struct {
	client *http.Client
	opts   HTTPOptions
}

func NewHTTPProvider(opts HTTPOptions) *HTTPProvider {
	if opts.Timeout <= 0 {
		opts = DefaultHTTPOptions()
	}
	if opts.UserAgent == "" {
		opts.UserAgent = DefaultHTTPOptions().UserAgent
	}
	if opts.Endpoint == "" {
		opts.Endpoint = DefaultHTTPOptions().Endpoint
	}
	return &HTTPProvider{
		client: &http.Client{Timeout: opts.Timeout},
		opts:   opts,
	}
}

func (h *HTTPProvider) Name() string {
	return "http"
}

func (h *HTTPProvider) Search(query string) ([]Result, error) {
	u, err := url.Parse(h.opts.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("search http: endpoint: %w", err)
	}
	q := u.Query()
	q.Set("q", query)
	q.Set("format", "json")
	q.Set("no_redirect", "1")
	q.Set("no_html", "1")
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("search http: build request: %w", err)
	}
	req.Header.Set("User-Agent", h.opts.UserAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("search http: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("search http: read body: %w", err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return parseDDGResponse(body, query)
	case http.StatusTooManyRequests:
		return nil, fmt.Errorf("search http: rate limited (429)")
	case http.StatusRequestTimeout, http.StatusGatewayTimeout:
		return nil, fmt.Errorf("search http: timeout status (%d)", resp.StatusCode)
	default:
		if resp.StatusCode >= 500 {
			return nil, fmt.Errorf("search http: server error (%d)", resp.StatusCode)
		}
		return nil, fmt.Errorf("search http: unexpected status (%d)", resp.StatusCode)
	}
}

type ddgResponse struct {
	AbstractURL   string     `json:"AbstractURL"`
	AbstractText  string     `json:"AbstractText"`
	RelatedTopics []ddgTopic `json:"RelatedTopics"`
}

type ddgTopic struct {
	FirstURL string     `json:"FirstURL"`
	Text     string     `json:"Text"`
	Topics   []ddgTopic `json:"Topics"`
}

func parseDDGResponse(body []byte, query string) ([]Result, error) {
	var payload ddgResponse
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("search http: decode: %w", err)
	}
	var out []Result
	if payload.AbstractURL != "" {
		out = append(out, Result{
			Title:   truncate(payload.AbstractText, 80),
			URL:     payload.AbstractURL,
			Snippet: strings.TrimSpace(payload.AbstractText + " | " + query),
		})
	}
	collectDDGTopics(payload.RelatedTopics, &out, query)
	if len(out) == 0 {
		return nil, fmt.Errorf("search http: no results for query %q", query)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].URL < out[j].URL
	})
	if len(out) > 10 {
		out = out[:10]
	}
	return out, nil
}

func collectDDGTopics(topics []ddgTopic, out *[]Result, query string) {
	for _, t := range topics {
		if len(t.Topics) > 0 {
			collectDDGTopics(t.Topics, out, query)
			continue
		}
		if t.FirstURL == "" {
			continue
		}
		title := t.Text
		if idx := strings.Index(title, " - "); idx > 0 {
			title = title[:idx]
		}
		*out = append(*out, Result{
			Title:   truncate(title, 80),
			URL:     t.FirstURL,
			Snippet: strings.TrimSpace(t.Text + " | " + query),
		})
	}
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
