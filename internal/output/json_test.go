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
	})
	if out["status"] != "success" {
		t.Fatal("expected success status")
	}
}

func TestFormatMarkdown(t *testing.T) {
	md := FormatMarkdown(&run.Run{
		ID:     "run_test",
		Status: "success",
		Steps:  []string{"search", "crawl"},
	})
	if md == "" {
		t.Fatal("expected markdown")
	}
}
