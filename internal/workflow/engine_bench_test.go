package workflow

import "testing"

func BenchmarkEngineRunTaskDeterministic(b *testing.B) {
	cfg := Config{
		Mode:      "deterministic",
		StorePath: b.TempDir(),
		Search:    "mock",
		Crawler:   "mock",
	}
	engine, err := NewEngine(cfg)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := engine.RunTask("benchmark task for agentbridge"); err != nil {
			b.Fatal(err)
		}
	}
}
