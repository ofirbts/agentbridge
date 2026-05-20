package crawl

import "strings"

func NewCrawler(kind string) Crawler {
	switch strings.ToLower(kind) {
	case "http":
		return NewHTTPCrawler(DefaultHTTPOptions())
	default:
		return NewMockCrawler()
	}
}
