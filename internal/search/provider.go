package search

type Result struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Snippet string `json:"snippet"`
}

type SearchProvider interface {
	Search(query string) ([]Result, error)
	Name() string
}
