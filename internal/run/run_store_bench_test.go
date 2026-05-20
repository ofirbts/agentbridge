package run

import "testing"

func BenchmarkRunStoreInsertGet(b *testing.B) {
	dir := b.TempDir()
	store := NewRunStore(dir)
	record := &Run{
		ID:         "run_bench",
		Status:     "success",
		Steps:      []string{"search", "crawl", "extract", "normalize"},
		DurationMS: 10,
		StepLogs: []StepLog{
			{Step: "search", Status: "ok", Provider: "mock", DurationMS: 1},
			{Step: "crawl", Status: "ok", Provider: "mock", DurationMS: 2},
		},
	}
	if err := store.Insert("run_bench", record); err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := store.Get("run_bench"); err != nil {
			b.Fatal(err)
		}
	}
}
