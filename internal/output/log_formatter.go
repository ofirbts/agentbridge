package output

import (
	"fmt"
	"strings"

	"github.com/yourname/agentbridge/internal/observability"
)

func FormatStepLogs(logs []observability.StepLog) string {
	lines := make([]string, 0, len(logs))
	for _, l := range logs {
		lines = append(lines, fmt.Sprintf("[%s] %s status=%s retries=%d duration_ms=%d provider=%s query=%q",
			l.RunID, l.Step, l.Status, l.Retries, l.DurationMS, l.Provider, l.Query))
	}
	return strings.Join(lines, "\n")
}
