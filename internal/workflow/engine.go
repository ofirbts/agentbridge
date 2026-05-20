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
		Mode:    e.cfg.Mode,
		Search:  e.search.Name(),
		Steps:   DefaultSteps(),
		Errors:  []string{},
		Crawler: e.crawler.Name(),
		MCP:     e.cfg.MCPEndpoint,
	}

	recorder.AddStepWithMeta("plan", "ok", "workflow", task, 0, 0)

	var searchResults []search.Result
	var page []byte
	var extractErr, crawlErr, searchErr error
	var mcpPreview string

	if e.mcp != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		tools, err := e.mcp.ListTools(ctx)
		cancel()
		if err != nil {
			appendRunError(record, "mcp", err)
		} else if len(tools) > 0 {
			mcpPreview = tools[0]
			recorder.AddStepWithMeta("mcp", "ok", "mcp", tools[0], 0, 0)
		}
	}

	searchStart := time.Now()
	retriesBefore := e.policy.RetriesDone()
	searchErr = failure.RetryWithPolicy(e.policy, func() error {
		var err error
		searchResults, err = e.search.Search(task)
		return err
	})
	searchRetries := int(e.policy.RetriesDone() - retriesBefore)
	recorder.AddStepWithMeta("search", statusFor(searchErr), e.search.Name(), task, searchRetries, time.Since(searchStart).Milliseconds())
	appendRunError(record, "search", searchErr)

	url := "https://example.com"
	if len(searchResults) > 0 {
		url = searchResults[0].URL
	}

	crawlStart := time.Now()
	retriesBefore = e.policy.RetriesDone()
	crawlErr = failure.RetryWithPolicy(e.policy, func() error {
		var err error
		page, err = e.crawler.Fetch(url)
		return err
	})
	crawlRetries := int(e.policy.RetriesDone() - retriesBefore)
	recorder.AddStepWithMeta("crawl", statusFor(crawlErr), e.crawler.Name(), url, crawlRetries, time.Since(crawlStart).Milliseconds())
	appendRunError(record, "crawl", crawlErr)

	extractStart := time.Now()
	text, extractErr := e.extractor.Extract(page)
	recorder.AddStepWithMeta("extract", statusFor(extractErr), "extractor", task, 0, time.Since(extractStart).Milliseconds())
	appendRunError(record, "extract", extractErr)

	normalizeStart := time.Now()
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
	record.StepLogs = run.StepLogsFromTracer(e.tracer.Logs())

	_ = e.runs.Insert(runID, record)

	return record, firstError(searchErr, crawlErr, extractErr)
}

func appendRunError(record *run.Run, step string, err error) {
	if err == nil {
		return
	}
	record.Errors = append(record.Errors, err.Error())
	record.ClassifiedErrors = append(record.ClassifiedErrors, run.NewClassifiedError(step, err))
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
