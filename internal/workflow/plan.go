package workflow

import "github.com/ofirbts/agentbridge/pkg/mcp"

type PlanStep struct {
	Name        string `json:"name"`
	Provider    string `json:"provider"`
	Description string `json:"description"`
}

type ObservabilityPlan struct {
	InspectFields []string `json:"inspect_fields"`
	PerStep       []string `json:"per_step_metrics"`
	ErrorClasses  []string `json:"error_classes"`
}

type Plan struct {
	Task          string            `json:"task"`
	Mode          string            `json:"mode"`
	Search        string            `json:"search"`
	Crawler       string            `json:"crawler"`
	MCPEndpoint   string            `json:"mcp_endpoint,omitempty"`
	Steps         []PlanStep        `json:"steps"`
	Providers     []string          `json:"providers"`
	Behavior      string            `json:"expected_behavior"`
	RunIDHint     string            `json:"run_id_hint,omitempty"`
	Observability ObservabilityPlan `json:"observability"`
}

func DefaultSteps() []string {
	return []string{"search", "crawl", "extract", "normalize"}
}

func defaultPlanSteps(search, crawler string) []PlanStep {
	return []PlanStep{
		{Name: "search", Provider: search, Description: "resolve candidate URLs for the task"},
		{Name: "crawl", Provider: crawler, Description: "fetch page content with retries"},
		{Name: "extract", Provider: "extractor", Description: "strip HTML and extract text"},
		{Name: "normalize", Provider: "normalizer", Description: "structure output for downstream agents"},
	}
}

func defaultObservabilityPlan() ObservabilityPlan {
	return ObservabilityPlan{
		InspectFields: []string{
			"duration_ms",
			"retries",
			"step_logs",
			"classified_errors",
			"search",
			"crawler",
			"mode",
		},
		PerStep: []string{
			"duration_ms",
			"retries",
			"provider",
			"status",
			"error_class",
		},
		ErrorClasses: []string{
			"rate_limit",
			"timeout",
			"server_error",
			"client_error",
			"network",
			"provider",
			"unknown",
		},
	}
}

func ExplainPlan(task, mode, search, crawler, mcpEndpoint, mcpTransport string) Plan {
	if search == "" {
		search = "mock"
	}
	if crawler == "" {
		crawler = "mock"
	}
	effectiveSearch := search
	if mode == "deterministic" {
		effectiveSearch = "mock"
	}
	behavior := "best-effort execution with retries"
	runHint := ""
	if mode == "deterministic" {
		behavior = "stable ordering, fixed run id seed, reproducible mock outputs"
		if search != "mock" {
			behavior += " (search forced to mock in deterministic mode)"
		}
		runHint = "deterministic run id derived from task hash"
	}
	providers := []string{effectiveSearch, crawler, "extractor", "normalizer"}
	if mcpEndpoint != "" {
		transport := mcpTransport
		if transport == "" {
			transport = mcp.DefaultTransport(mcpEndpoint)
		}
		providers = append(providers, "mcp_"+transport)
		behavior += "; MCP " + transport + " transport for tool-backed providers"
	}
	behavior += "; inspect shows per-step duration, retries, providers, and error classification"
	return Plan{
		Task:          task,
		Mode:          mode,
		Search:        search,
		Crawler:       crawler,
		MCPEndpoint:   mcpEndpoint,
		Steps:         defaultPlanSteps(effectiveSearch, crawler),
		Providers:     providers,
		Behavior:      behavior,
		RunIDHint:     runHint,
		Observability: defaultObservabilityPlan(),
	}
}
