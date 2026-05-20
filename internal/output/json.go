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
	if r.MCP != "" {
		out["mcp_endpoint"] = r.MCP
	}
	return out
}
