package run

import (
	"fmt"

	"github.com/google/uuid"
)

type Run struct {
	ID         string   `json:"run_id"`
	Status     string   `json:"status"`
	Steps      []string `json:"steps"`
	Errors     []string `json:"errors"`
	Retries    int      `json:"retries"`
	DurationMS int64    `json:"duration_ms"`
	Crawler    string   `json:"crawler,omitempty"`
	MCP        string   `json:"mcp_endpoint,omitempty"`
	Result     any      `json:"result"`
	Logs       any      `json:"logs,omitempty"`
}

func GenerateID() string {
	return "run_" + uuid.NewString()[:8]
}

func GenerateDeterministicID(task string) string {
	var sum uint64
	for _, c := range task {
		sum = sum*31 + uint64(c)
	}
	return fmt.Sprintf("run_det_%08x", sum&0xffffffff)
}
