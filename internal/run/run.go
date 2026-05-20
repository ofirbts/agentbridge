package run

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/ofirbts/agentbridge/internal/observability"
)

type StepLog struct {
	Step       string `json:"step"`
	Status     string `json:"status"`
	Provider   string `json:"provider"`
	DurationMS int64  `json:"duration_ms"`
	Retries    int    `json:"retries"`
	Target     string `json:"target,omitempty"`
	Message    string `json:"message,omitempty"`
	ErrorClass string `json:"error_class,omitempty"`
}

type ClassifiedError struct {
	Step    string `json:"step,omitempty"`
	Message string `json:"message"`
	Class   string `json:"class"`
}

type Run struct {
	ID               string            `json:"run_id"`
	Status           string            `json:"status"`
	Mode             string            `json:"mode,omitempty"`
	Search           string            `json:"search,omitempty"`
	Steps            []string          `json:"steps"`
	Errors           []string          `json:"errors"`
	ClassifiedErrors []ClassifiedError `json:"classified_errors,omitempty"`
	Retries          int               `json:"retries"`
	DurationMS       int64             `json:"duration_ms"`
	Crawler          string            `json:"crawler,omitempty"`
	MCP              string            `json:"mcp_endpoint,omitempty"`
	Result           any               `json:"result"`
	StepLogs         []StepLog         `json:"step_logs,omitempty"`
}

func NewClassifiedError(step string, err error) ClassifiedError {
	msg := ""
	if err != nil {
		msg = err.Error()
	}
	return ClassifiedError{
		Step:    step,
		Message: msg,
		Class:   observability.ClassifyError(msg),
	}
}

func StepLogsFromTracer(logs []observability.StepLog) []StepLog {
	out := make([]StepLog, 0, len(logs))
	for _, l := range logs {
		sl := StepLog{
			Step:       l.Step,
			Status:     l.Status,
			Provider:   l.Provider,
			DurationMS: l.DurationMS,
			Retries:    l.Retries,
			Target:     l.Query,
			Message:    l.Message,
		}
		if l.Status == "error" && sl.Message != "" {
			sl.ErrorClass = observability.ClassifyError(sl.Message)
		}
		out = append(out, sl)
	}
	return out
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
