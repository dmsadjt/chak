package search

type SearchResultData struct {
	SzTitle string `json:"title"`
	SzSnippet string `json:"snippet"`
	SzURL string `json:"url"`
}

type SearchInterface interface {
	Search(SzQuery string) ([]SearchResultData, error)
}
