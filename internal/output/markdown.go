package output

import (
	"fmt"
	"strings"

	"github.com/ofirbts/agentbridge/internal/run"
)

func FormatMarkdown(r *run.Run) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# Run %s\n\n", r.ID)
	fmt.Fprintf(&b, "- Status: %s\n", r.Status)
	fmt.Fprintf(&b, "- Crawler: %s\n", r.Crawler)
	if r.MCP != "" {
		fmt.Fprintf(&b, "- MCP: %s\n", r.MCP)
	}
	fmt.Fprintf(&b, "- Duration: %dms\n", r.DurationMS)
	fmt.Fprintf(&b, "- Retries: %d\n\n", r.Retries)
	fmt.Fprintf(&b, "## Steps\n\n")
	for _, step := range r.Steps {
		fmt.Fprintf(&b, "- %s\n", step)
	}
	if len(r.Errors) > 0 {
		fmt.Fprintf(&b, "\n## Errors\n\n")
		for _, e := range r.Errors {
			fmt.Fprintf(&b, "- %s\n", e)
		}
	}
	return b.String()
}
