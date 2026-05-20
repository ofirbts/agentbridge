package workflow

import (
	"context"
	"time"

	"github.com/ofirbts/agentbridge/internal/crawl"
	"github.com/ofirbts/agentbridge/internal/extract"
	"github.com/ofirbts/agentbridge/internal/failure"
	"github.com/ofirbts/agentbridge/internal/observability"
	"github.com/ofirbts/agentbridge/internal/run"
	"github.com/ofirbts/agentbridge/internal/search"
	"github.com/ofirbts/agentbridge/pkg/mcp"
)

type Config struct {
	Mode         string
	ConfigFile   string
	StorePath    string
	Crawler      string
	MCPEndpoint  string
	MCPTransport string
}

type Engine struct {
	cfg        Config
	policy     *failure.RetryPolicy
	tracer     *observability.Tracer
	runs       *run.RunStore
	search     search.SearchProvider
	crawler    crawl.Crawler
	extractor  *extract.Extractor
	normalizer *extract.Normalizer
	mcp        mcp.Client
}

func NewEngine(cfg Config) *Engine {
	var mcpClient mcp.Client
	if cfg.MCPEndpoint != "" {
		transport := cfg.MCPTransport
		if transport == "" {
			transport = mcp.DefaultTransport(cfg.MCPEndpoint)
		}
		client, err := mcp.NewClient(cfg.MCPEndpoint, transport)
		if err == nil {
			mcpClient = client
		}
	}
	return &Engine{
		cfg:        cfg,
		policy:     ptr(failure.NewRetryPolicy(3, 100)),
		tracer:     observability.NewTracer(),
		runs:       run.NewRunStore(cfg.StorePath),
		search:     search.NewMockProvider(),
		crawler:    crawl.NewCrawler(cfg.Crawler),
		extractor:  extract.NewExtractor(),
		normalizer: extract.NewNormalizer(),
		mcp:        mcpClient,
	}
}

type Result struct {
	Status     string   `json:"status"`
	Steps      []string `json:"steps"`
	Errors     []string `json:"errors"`
	Retries    int      `json:"retries"`
	DurationMS int64    `json:"duration_ms"`
	RunID      string   `json:"run_id"`
	Crawler    string   `json:"crawler"`
	MCP        string   `json:"mcp_endpoint,omitempty"`
	Result     any      `json:"result"`
}

func (e *Engine) RunTask(task string) (*Result, error) {
	start := time.Now()
	e.policy.Reset()

	runID := run.GenerateID()
	if e.cfg.Mode == "deterministic" {
		runID = run.GenerateDeterministicID(task)
	}

	recorder := e.tracer.StartRun(runID)
	result := &Result{
		Status:  "running",
		RunID:   runID,
		Steps:   DefaultSteps(),
		Errors:  []string{},
		Crawler: e.crawler.Name(),
		MCP:     e.cfg.MCPEndpoint,
	}

	recorder.AddStep("plan", "started planning search and crawl steps")

	var searchResults []search.Result
	var page []byte
	var extractErr, crawlErr, searchErr error
	var mcpPreview string

	if e.mcp != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		tools, err := e.mcp.ListTools(ctx)
		cancel()
		if err != nil {
			result.Errors = append(result.Errors, err.Error())
		} else if len(tools) > 0 {
			mcpPreview = tools[0]
			recorder.AddStep("mcp", "listed tools: "+tools[0])
		}
	}

	searchStart := time.Now()
	searchErr = failure.RetryWithPolicy(e.policy, func() error {
		var err error
		searchResults, err = e.search.Search(task)
		return err
	})
	recorder.AddStepWithMeta("search", statusFor(searchErr), e.search.Name(), task, int(e.policy.RetriesDone()), time.Since(searchStart).Milliseconds())
	if searchErr != nil {
		result.Errors = append(result.Errors, searchErr.Error())
	}

	url := "https://example.com"
	if len(searchResults) > 0 {
		url = searchResults[0].URL
	}

	crawlStart := time.Now()
	crawlErr = failure.RetryWithPolicy(e.policy, func() error {
		var err error
		page, err = e.crawler.Fetch(url)
		return err
	})
	recorder.AddStepWithMeta("crawl", statusFor(crawlErr), e.crawler.Name(), url, int(e.policy.RetriesDone()), time.Since(crawlStart).Milliseconds())
	if crawlErr != nil {
		result.Errors = append(result.Errors, crawlErr.Error())
	}

	extractStart := time.Now()
	recorder.AddStep("extract", "parsing content")
	text, extractErr := e.extractor.Extract(page)
	recorder.AddStepWithMeta("extract", statusFor(extractErr), "extractor", task, 0, time.Since(extractStart).Milliseconds())
	if extractErr != nil {
		result.Errors = append(result.Errors, extractErr.Error())
	}

	normalizeStart := time.Now()
	recorder.AddStep("normalize", "structuring output")
	normalized := e.normalizer.Normalize(text)
	if e.cfg.Mode == "deterministic" {
		normalized = deterministicNormalize(normalized)
	}
	recorder.AddStepWithMeta("normalize", "ok", "normalizer", task, 0, time.Since(normalizeStart).Milliseconds())

	result.Result = map[string]any{
		"query":       task,
		"url":         url,
		"normalized":  normalized,
		"hits":        len(searchResults),
		"mcp_preview": mcpPreview,
	}

	if searchErr != nil || crawlErr != nil || extractErr != nil {
		result.Status = "failed"
	} else {
		result.Status = "success"
	}

	result.DurationMS = time.Since(start).Milliseconds()
	result.Retries = int(e.policy.RetriesDone())

	stored := &run.Run{
		ID:         runID,
		Status:     result.Status,
		Steps:      result.Steps,
		Errors:     result.Errors,
		Retries:    result.Retries,
		DurationMS: result.DurationMS,
		Result:     result.Result,
		Logs:       e.tracer.Logs(),
		Crawler:    result.Crawler,
		MCP:        result.MCP,
	}
	_ = e.runs.Insert(runID, stored)

	return result, firstError(searchErr, crawlErr, extractErr)
}

func statusFor(err error) string {
	if err != nil {
		return "error"
	}
	return "ok"
}

func deterministicNormalize(in map[string]any) map[string]any {
	out := map[string]any{
		"text":       in["text"],
		"word_count": in["word_count"],
	}
	if tokens, ok := in["tokens"].([]string); ok {
		out["tokens"] = tokens
	}
	return out
}

func firstError(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

func ptr(p failure.RetryPolicy) *failure.RetryPolicy {
	return &p
}
