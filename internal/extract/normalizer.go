package extract

import (
	"strings"
)

type Normalizer struct{}

func NewNormalizer() *Normalizer {
	return &Normalizer{}
}

func (n *Normalizer) Normalize(text string) map[string]any {
	words := strings.Fields(text)
	return map[string]any{
		"text":       text,
		"word_count": len(words),
		"tokens":     words,
	}
}
