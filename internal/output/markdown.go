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
	if r.Mode != "" {
		fmt.Fprintf(&b, "- Mode: %s\n", r.Mode)
	}
	if r.Search != "" {
		fmt.Fprintf(&b, "- Search: %s\n", r.Search)
	}
	fmt.Fprintf(&b, "- Crawler: %s\n", r.Crawler)
	if r.MCP != "" {
		fmt.Fprintf(&b, "- MCP: %s\n", r.MCP)
	}
	fmt.Fprintf(&b, "- Duration: %dms\n", r.DurationMS)
	fmt.Fprintf(&b, "- Retries: %d\n\n", r.Retries)

	if len(r.StepLogs) > 0 {
		fmt.Fprintf(&b, "## Step timeline\n\n")
		fmt.Fprintf(&b, "| Step | Status | Provider | Duration | Retries | Error class |\n")
		fmt.Fprintf(&b, "|------|--------|----------|----------|---------|-------------|\n")
		for _, s := range r.StepLogs {
			ec := s.ErrorClass
			if ec == "" {
				ec = "-"
			}
			fmt.Fprintf(&b, "| %s | %s | %s | %dms | %d | %s |\n",
				s.Step, s.Status, s.Provider, s.DurationMS, s.Retries, ec)
		}
		b.WriteString("\n")
	} else {
		fmt.Fprintf(&b, "## Steps\n\n")
		for _, step := range r.Steps {
			fmt.Fprintf(&b, "- %s\n", step)
		}
		b.WriteString("\n")
	}

	if len(r.ClassifiedErrors) > 0 {
		fmt.Fprintf(&b, "## Classified errors\n\n")
		for _, e := range r.ClassifiedErrors {
			step := e.Step
			if step == "" {
				step = "-"
			}
			fmt.Fprintf(&b, "- [%s] **%s**: %s\n", e.Class, step, e.Message)
		}
		b.WriteString("\n")
	} else if len(r.Errors) > 0 {
		fmt.Fprintf(&b, "## Errors\n\n")
		for _, e := range r.Errors {
			fmt.Fprintf(&b, "- %s\n", e)
		}
	}
	return b.String()
}
