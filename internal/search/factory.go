package search

import "strings"

func NewSearchProvider(kind string) SearchProvider {
	switch strings.ToLower(kind) {
	case "http", "duckduckgo":
		return NewHTTPProvider(DefaultHTTPOptions())
	default:
		return NewMockProvider()
	}
}
