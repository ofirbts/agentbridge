package run

import "testing"

func TestRunStorePersistAndLoad(t *testing.T) {
	dir := t.TempDir()
	store := NewRunStore(dir)
	run := &Run{
		ID:         "run_test123",
		Status:     "success",
		Steps:      []string{"search", "crawl"},
		Retries:    1,
		DurationMS: 100,
		Result:     map[string]any{"ok": true},
	}
	if err := store.Insert("run_test123", run); err != nil {
		t.Fatalf("insert: %v", err)
	}
	loaded, err := store.Get("run_test123")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if loaded.Status != "success" {
		t.Fatalf("expected success")
	}
}

func TestGenerateDeterministicIDStable(t *testing.T) {
	a := GenerateDeterministicID("same task")
	b := GenerateDeterministicID("same task")
	if a != b {
		t.Fatalf("ids differ")
	}
}
