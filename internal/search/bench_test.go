package search

import "testing"

func BenchmarkMockSearch(b *testing.B) {
	p := NewMockProvider()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := p.Search("benchmark query"); err != nil {
			b.Fatal(err)
		}
	}
}
