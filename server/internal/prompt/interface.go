package prompt

import "chak-server/internal/search"

type PromptInterface interface {
	Build(szQuery string, searchResultData []search.SearchResultData) string
}

