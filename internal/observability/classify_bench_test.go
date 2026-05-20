package observability

import "testing"

func BenchmarkClassifyError(b *testing.B) {
	msg := "http crawl: rate limited (429) for https://example.com"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ClassifyError(msg)
	}
}
