package output

import (
	"testing"

	"github.com/yourname/agentbridge/internal/workflow"
)

func TestFormatRunResult(t *testing.T) {
	out := FormatRunResult(&workflow.Result{
		Status: "success",
		RunID:  "run_test",
		Steps:  []string{"search"},
	})
	if out["status"] != "success" {
		t.Fatal("expected success status")
	}
}

func TestFormatMarkdown(t *testing.T) {
	md := FormatMarkdown(&workflow.Result{
		RunID:  "run_test",
		Status: "success",
		Steps:  []string{"search", "crawl"},
	})
	if md == "" {
		t.Fatal("expected markdown")
	}
}
