package prompt

import(
	"fmt"
	"chak-server/internal/search"
)

type PromptManager struct {}

func NewPromptManager() *PromptManager {
	return &PromptManager{}
}

func (promptMgr *PromptManager) Build(szQuery string, searchResultData []search.SearchResultData) string {
	if(len(searchResultData) == 0) {
		return szQuery 
	}

	prompt := "You have access to the following web search results. Use this information to answer the user's question accurately: "
	prompt += "=== SEARCH RESULTS ===\n\n"

	for i, result := range searchResultData {
		prompt += fmt.Sprintf("Result %d:\n", i+1)
		prompt += fmt.Sprintf("Title: %s\n", result.SzTitle)
		prompt += fmt.Sprintf("Content: %s:\n", result.SzSnippet)
		prompt += fmt.Sprintf("URL: %s:\n\n", result.SzURL)
	}

	prompt += "=== END SEARCH RESULTS ===\n\n"
	prompt += "Instructions: Based ONLY on the search results above, answer the following question. If the search don't contain relevant information, say so.\n\n"
	prompt += fmt.Sprintf("User question: %s", szQuery)

	return prompt
}
