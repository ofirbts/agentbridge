package output

import (
	"testing"

	"github.com/ofirbts/agentbridge/internal/run"
)

func TestFormatRunResult(t *testing.T) {
	out := FormatRunResult(&run.Run{
		Status: "success",
		ID:     "run_test",
		Steps:  []string{"search"},
		StepLogs: []run.StepLog{{
			Step: "search", Status: "ok", Provider: "mock", DurationMS: 1,
		}},
		ClassifiedErrors: []run.ClassifiedError{{
			Step: "crawl", Message: "timeout", Class: "timeout",
		}},
	})
	if out["status"] != "success" {
		t.Fatal("expected success status")
	}
	if _, ok := out["step_logs"]; !ok {
		t.Fatal("expected step_logs")
	}
	if _, ok := out["classified_errors"]; !ok {
		t.Fatal("expected classified_errors")
	}
}

func TestFormatMarkdownStepLogs(t *testing.T) {
	md := FormatMarkdown(&run.Run{
		ID:     "run_test",
		Status: "failed",
		Search: "mock",
		Steps:  []string{"search", "crawl"},
		StepLogs: []run.StepLog{
			{Step: "search", Status: "ok", Provider: "mock", DurationMS: 2, Retries: 0},
			{Step: "crawl", Status: "error", Provider: "http", DurationMS: 100, Retries: 3, ErrorClass: "rate_limit"},
		},
		ClassifiedErrors: []run.ClassifiedError{
			{Step: "crawl", Message: "rate limited", Class: "rate_limit"},
		},
	})
	if md == "" {
		t.Fatal("expected markdown")
	}
	if !contains(md, "Step timeline") || !contains(md, "rate_limit") {
		t.Fatalf("expected step table and error class: %s", md)
	}
}

func contains(s, sub string) bool {
	return len(sub) == 0 || (len(s) >= len(sub) && stringIndex(s, sub) >= 0)
}

func stringIndex(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
