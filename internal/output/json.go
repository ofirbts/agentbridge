package output

import "github.com/ofirbts/agentbridge/internal/run"

func FormatRunResult(r *run.Run) map[string]any {
	if r == nil {
		return map[string]any{"status": "failed"}
	}
	out := map[string]any{
		"status":      r.Status,
		"steps":       r.Steps,
		"errors":      r.Errors,
		"retries":     r.Retries,
		"duration_ms": r.DurationMS,
		"run_id":      r.ID,
		"crawler":     r.Crawler,
		"result":      r.Result,
	}
	if r.Mode != "" {
		out["mode"] = r.Mode
	}
	if r.Search != "" {
		out["search"] = r.Search
	}
	if r.MCP != "" {
		out["mcp_endpoint"] = r.MCP
	}
	if len(r.StepLogs) > 0 {
		out["step_logs"] = r.StepLogs
	}
	if len(r.ClassifiedErrors) > 0 {
		out["classified_errors"] = r.ClassifiedErrors
	}
	return out
}

func FormatInspect(r *run.Run) map[string]any {
	out := FormatRunResult(r)
	out["providers"] = map[string]string{
		"search":  r.Search,
		"crawler": r.Crawler,
	}
	return out
}
