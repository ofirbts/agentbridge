package workflow

type PlanStep struct {
	Name        string `json:"name"`
	Provider    string `json:"provider"`
	Description string `json:"description"`
}

type Plan struct {
	Task       string     `json:"task"`
	Mode       string     `json:"mode"`
	Crawler    string     `json:"crawler"`
	MCPEndpoint string    `json:"mcp_endpoint,omitempty"`
	Steps      []PlanStep `json:"steps"`
	Providers  []string   `json:"providers"`
	Behavior   string     `json:"expected_behavior"`
	RunIDHint  string     `json:"run_id_hint,omitempty"`
}

func DefaultSteps() []string {
	return []string{"search", "crawl", "extract", "normalize"}
}

func defaultPlanSteps(crawler string) []PlanStep {
	return []PlanStep{
		{Name: "search", Provider: "mock_search", Description: "resolve candidate URLs for the task"},
		{Name: "crawl", Provider: crawler, Description: "fetch page content with retries"},
		{Name: "extract", Provider: "extractor", Description: "strip HTML and extract text"},
		{Name: "normalize", Provider: "normalizer", Description: "structure output for downstream agents"},
	}
}

func ExplainPlan(task, mode, crawler, mcpEndpoint string) Plan {
	if crawler == "" {
		crawler = "mock"
	}
	behavior := "best-effort execution with retries"
	runHint := ""
	if mode == "deterministic" {
		behavior = "stable ordering, fixed run id seed, reproducible mock outputs"
		runHint = "deterministic run id derived from task hash"
	}
	providers := []string{"mock_search", crawler + "_crawler", "extractor", "normalizer"}
	if mcpEndpoint != "" {
		providers = append(providers, "mcp_stub")
		behavior += "; MCP stub available for future tool-backed providers"
	}
	return Plan{
		Task:        task,
		Mode:        mode,
		Crawler:     crawler,
		MCPEndpoint: mcpEndpoint,
		Steps:       defaultPlanSteps(crawler),
		Providers:   providers,
		Behavior:    behavior,
		RunIDHint:   runHint,
	}
}

func BuildPlan(task string) Plan {
	return ExplainPlan(task, "normal", "mock", "")
}
