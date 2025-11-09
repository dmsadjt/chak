package prompt

import (
	"chak-server/internal/search"
	"chak-server/internal/types"
	"chak-server/internal/memory"
)

type PromptInterface interface {
	Build(messageList []types.Message, searchResultData []search.SearchResultData, memories []memory.MemoryEntry) string
}
