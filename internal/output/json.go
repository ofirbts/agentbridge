package output

import "github.com/ofirbts/agentbridge/internal/workflow"

func FormatRunResult(r *workflow.Result) map[string]any {
	if r == nil {
		return map[string]any{"status": "failed"}
	}
	out := map[string]any{
		"status":      r.Status,
		"steps":       r.Steps,
		"errors":      r.Errors,
		"retries":     r.Retries,
		"duration_ms": r.DurationMS,
		"run_id":      r.RunID,
		"crawler":     r.Crawler,
		"result":      r.Result,
	}
	if r.MCP != "" {
		out["mcp_endpoint"] = r.MCP
	}
	return out
}
