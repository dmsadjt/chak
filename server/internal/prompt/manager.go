package prompt

import (
	"chak-server/internal/memory"
	"chak-server/internal/search"
	"chak-server/internal/types"
	"fmt"
	"time"
)

type PromptManager struct {}

func NewPromptManager() *PromptManager {
	return &PromptManager{}
}

func (promptMgr *PromptManager) Build(messageList []types.Message , searchResultData []search.SearchResultData, memories []memory.MemoryEntry) string {
	prompt := ""

	prompt += fmt.Sprintf("Current date and time: %s\n\n", time.Now().Format(time.RFC1123))

	if len(memories) > 0 {
		prompt += "=== RELEVANT PAST CONVERSATION ===\n\n"
		for i, mem := range memories {
			prompt += fmt.Sprintf("Memory %d:\n%s\n\n", i+1, mem.SzContent)
		}
		prompt += "=== END MEMORIES ===\n\n"

	}

	if len(searchResultData) > 0 {
		prompt = "You have access to the following web search results. Use this information to answer the user's question accurately: "
		prompt += "=== SEARCH RESULTS ===\n\n"

		for i, result := range searchResultData {
			prompt += fmt.Sprintf("Result %d:\n", i+1)
			prompt += fmt.Sprintf("Title: %s\n", result.SzTitle)
			prompt += fmt.Sprintf("Content: %s:\n", result.SzSnippet)
			prompt += fmt.Sprintf("URL: %s:\n\n", result.SzURL)
		}

		prompt += "=== END SEARCH RESULTS ===\n\n"
		prompt += "Instructions: Prioritize the search results. You may make simple inferences from them if needed, but do not add information that is not supported by the search results.\n\n"
	} else {
		prompt += "Instructions: Provide a clear and direct answer to the question.\n"
		prompt += "Use the conversation history only if it adds useful context.\n"
		prompt += "Do not overanalyze or reference the history unless necessary.\n\n"
	}

	prompt += "Conversation History: \n"
	for _, msg := range messageList{
		prompt += fmt.Sprintf("%s: %s\n", msg.SzRole, msg.SzContent)
	}

	szLastMessage := messageList[len(messageList)-1]
	prompt += fmt.Sprintf("User question: %s", szLastMessage.SzContent)

	return prompt
}
