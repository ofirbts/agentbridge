package extract

import "testing"

func TestExtractStripsHTML(t *testing.T) {
	e := NewExtractor()
	text, err := e.Extract([]byte("<html><body><h1>Hi</h1></body></html>"))
	if err != nil {
		t.Fatal(err)
	}
	if text != "Hi" {
		t.Fatalf("expected Hi, got %q", text)
	}
}

func TestNormalizeStructured(t *testing.T) {
	n := NewNormalizer()
	out := n.Normalize("one two three")
	if out["word_count"] != 3 {
		t.Fatalf("expected 3 words")
	}
}
