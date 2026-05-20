package agentbridge_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestCLIHelp(t *testing.T) {
	root := projectRoot(t)
	cmd := exec.Command("go", "run", ".", "--help")
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("help failed: %v\n%s", err, out)
	}
	if !contains(string(out), "agentbridge") {
		t.Fatalf("expected help output")
	}
}

func TestCLIRunInspectFlow(t *testing.T) {
	root := projectRoot(t)
	store := filepath.Join(t.TempDir(), "runs")
	run := exec.Command("go", "run", ".", "run", "e2e flow task", "--mode", "deterministic", "--store", store)
	run.Dir = root
	out, err := run.CombinedOutput()
	if err != nil {
		t.Fatalf("run failed: %v\n%s", err, out)
	}
	id := extractJSONField(string(out), "run_id")
	if id == "" {
		t.Fatalf("missing run_id in output: %s", out)
	}
	inspect := exec.Command("go", "run", ".", "inspect", id, "--store", store)
	inspect.Dir = root
	inspectOut, err := inspect.CombinedOutput()
	if err != nil {
		t.Fatalf("inspect failed: %v\n%s", err, inspectOut)
	}
	if !contains(string(inspectOut), id) {
		t.Fatalf("inspect output missing run id")
	}
}

func projectRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Clean(filepath.Join(wd, ".."))
}

func contains(s, sub string) bool {
	return len(sub) == 0 || (len(s) >= len(sub) && (s == sub || len(s) > 0 && stringIndex(s, sub) >= 0))
}

func stringIndex(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

func extractJSONField(body, field string) string {
	key := `"` + field + `": "`
	start := stringIndex(body, key)
	if start < 0 {
		return ""
	}
	start += len(key)
	end := start
	for end < len(body) && body[end] != '"' {
		end++
	}
	return body[start:end]
}
