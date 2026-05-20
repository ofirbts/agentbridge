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
	StorePath    string
	Search       string
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

func NewEngine(cfg Config) (*Engine, error) {
	var mcpClient mcp.Client
	if cfg.MCPEndpoint != "" {
		transport := cfg.MCPTransport
		if transport == "" {
			transport = mcp.DefaultTransport(cfg.MCPEndpoint)
		}
		client, err := mcp.NewClient(cfg.MCPEndpoint, transport)
		if err != nil {
			return nil, err
		}
		mcpClient = client
	}
	searchKind := cfg.Search
	if cfg.Mode == "deterministic" {
		searchKind = "mock"
	}
	return &Engine{
		cfg:        cfg,
		policy:     ptr(failure.NewRetryPolicy(3, 100)),
		tracer:     observability.NewTracer(),
		runs:       run.NewRunStore(cfg.StorePath),
		search:     search.NewSearchProvider(searchKind),
		crawler:    crawl.NewCrawler(cfg.Crawler),
		extractor:  extract.NewExtractor(),
		normalizer: extract.NewNormalizer(),
		mcp:        mcpClient,
	}, nil
}

func (e *Engine) RunTask(task string) (*run.Run, error) {
	start := time.Now()
	e.policy.Reset()

	runID := run.GenerateID()
	if e.cfg.Mode == "deterministic" {
		runID = run.GenerateDeterministicID(task)
	}

	recorder := e.tracer.StartRun(runID)
	record := &run.Run{
		Status:  "running",
		ID:      runID,
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
			record.Errors = append(record.Errors, err.Error())
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
		record.Errors = append(record.Errors, searchErr.Error())
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
		record.Errors = append(record.Errors, crawlErr.Error())
	}

	extractStart := time.Now()
	recorder.AddStep("extract", "parsing content")
	text, extractErr := e.extractor.Extract(page)
	recorder.AddStepWithMeta("extract", statusFor(extractErr), "extractor", task, 0, time.Since(extractStart).Milliseconds())
	if extractErr != nil {
		record.Errors = append(record.Errors, extractErr.Error())
	}

	normalizeStart := time.Now()
	recorder.AddStep("normalize", "structuring output")
	normalized := e.normalizer.Normalize(text)
	if e.cfg.Mode == "deterministic" {
		normalized = deterministicNormalize(normalized)
	}
	recorder.AddStepWithMeta("normalize", "ok", "normalizer", task, 0, time.Since(normalizeStart).Milliseconds())

	record.Result = map[string]any{
		"query":       task,
		"url":         url,
		"normalized":  normalized,
		"hits":        len(searchResults),
		"mcp_preview": mcpPreview,
	}

	if searchErr != nil || crawlErr != nil || extractErr != nil {
		record.Status = "failed"
	} else {
		record.Status = "success"
	}

	record.DurationMS = time.Since(start).Milliseconds()
	record.Retries = int(e.policy.RetriesDone())
	record.Logs = e.tracer.Logs()

	_ = e.runs.Insert(runID, record)

	return record, firstError(searchErr, crawlErr, extractErr)
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
