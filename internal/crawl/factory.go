package crawl

import "strings"

func NewCrawler(kind string) Crawler {
	switch strings.ToLower(kind) {
	case "http", "real":
		return NewHTTPCrawler(DefaultHTTPOptions())
	default:
		return NewMockCrawler()
	}
}
