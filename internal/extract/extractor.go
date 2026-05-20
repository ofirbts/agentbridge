package extract

import (
	"regexp"
	"strings"
)

var tagRe = regexp.MustCompile(`<[^>]+>`)

type Extractor struct{}

func NewExtractor() *Extractor {
	return &Extractor{}
}

func (e *Extractor) Extract(content []byte) (string, error) {
	text := tagRe.ReplaceAllString(string(content), " ")
	text = strings.Join(strings.Fields(text), " ")
	return text, nil
}
