package run

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunStorePersistAndLoad(t *testing.T) {
	dir := t.TempDir()
	store := NewRunStore(dir)
	r := &Run{
		ID:         "run_test123",
		Status:     "success",
		Mode:       "normal",
		Search:     "mock",
		Steps:      []string{"search", "crawl"},
		Retries:    1,
		DurationMS: 100,
		StepLogs: []StepLog{
			{Step: "search", Status: "ok", Provider: "mock", DurationMS: 10, Retries: 0},
		},
		Result: map[string]any{"ok": true},
	}
	if err := store.Insert("run_test123", r); err != nil {
		t.Fatalf("insert: %v", err)
	}
	loaded, err := store.Get("run_test123")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if loaded.Status != "success" || len(loaded.StepLogs) != 1 {
		t.Fatalf("expected step logs on load")
	}
}

func TestRunStoreLegacyLogsField(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "run_legacy.json")
	legacy := `{
  "run_id": "run_legacy",
  "status": "success",
  "steps": ["search"],
  "errors": [],
  "retries": 0,
  "duration_ms": 5,
  "logs": [
    {"step": "search", "status": "ok", "provider": "mock", "duration_ms": 5, "retries": 0, "query": "task"}
  ]
}`
	if err := os.WriteFile(path, []byte(legacy), 0o644); err != nil {
		t.Fatal(err)
	}
	loaded, err := decodeRun([]byte(legacy))
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded.StepLogs) != 1 || loaded.StepLogs[0].Step != "search" {
		t.Fatalf("legacy logs not migrated: %#v", loaded.StepLogs)
	}
}

func TestGenerateDeterministicIDStable(t *testing.T) {
	a := GenerateDeterministicID("same task")
	b := GenerateDeterministicID("same task")
	if a != b {
		t.Fatalf("ids differ")
	}
}

func TestNewClassifiedError(t *testing.T) {
	ce := NewClassifiedError("crawl", errSample("http crawl: rate limited (429)"))
	if ce.Class != "rate_limit" || ce.Step != "crawl" {
		t.Fatalf("unexpected classified error: %#v", ce)
	}
}

type errSample string

func (e errSample) Error() string { return string(e) }
