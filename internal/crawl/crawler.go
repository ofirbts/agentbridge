package crawl

type Crawler interface {
	Fetch(url string) ([]byte, error)
	Name() string
}
