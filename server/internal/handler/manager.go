package handler

import (
	"chak-server/internal/memory"
	"chak-server/internal/ollama"
	"chak-server/internal/prompt"
	"chak-server/internal/search"
	"chak-server/internal/types"
	"encoding/json"
	"log"
	"net/http"
)

const MaxRememberedMessages = 10

type ChatRequest struct {
    MessageList []types.Message `json:"messages"`
    Search bool   `json:"search"`
    Model  string `json:"model"`
}

type ChatResponse struct {
    Response string                `json:"response"`
    Sources  []search.SearchResultData `json:"sources,omitempty"`
    Tokens   int                   `json:"tokens"`
    Time     float64               `json:"time"`
}

type ChatHandlerManager struct {
	searchManager search.SearchInterface
	promptManager prompt.PromptInterface
	ollamaManager ollama.OllamaInterface
	memoryManager memory.MemoryInterface
}

func NewChatHandlerManager(sm search.SearchInterface, pm prompt.PromptInterface, om ollama.OllamaInterface, mm memory.MemoryInterface) *ChatHandlerManager {
	return &ChatHandlerManager{
		searchManager: sm,
		promptManager: pm,
		ollamaManager: om,
		memoryManager: mm,
	}
}

func (chatManager *ChatHandlerManager) buildContext(messages []types.Message) []types.Message {
	if len(messages) - 1 > MaxRememberedMessages {
		messages = messages[len(messages)-MaxRememberedMessages:len(messages)-1]
	}
	return messages
}

func (chatManager *ChatHandlerManager) HandleChat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	messages := req.MessageList
	if len(messages) == 0 {
		http.Error(w, "Empty conversation", http.StatusBadRequest)
		return
	}

	messages = chatManager.buildContext(messages)

	szLastMessage := messages[len(messages)-1].SzContent
	ctx := r.Context()

	relevantMemories, err := chatManager.memoryManager.RetrieveRelevantContext(ctx, szLastMessage, 3)
	if err != nil {
		log.Printf("Memory retrieval error: %v", err)
	}

	var searchResultData []search.SearchResultData
	
	if req.Search && len(messages) > 0 {
		lastMsg := messages[len(messages)-1].SzContent

		if result, err := chatManager.searchManager.Search(lastMsg); err == nil {
			searchResultData = result
		} else {
			http.Error(w, "Search error", http.StatusInternalServerError)
			return
		}
	}

	szFinalPrompt := chatManager.promptManager.Build(messages, searchResultData, relevantMemories)

	ollamaResp, err := chatManager.ollamaManager.Generate(req.Model, szFinalPrompt)	
	if err != nil {
		http.Error(w, "Ollama error", http.StatusInternalServerError)
		return
	}

	resp := ChatResponse {
		Response: ollamaResp.SzResponse,
		Sources: searchResultData,
		Tokens: ollamaResp.ITotalTokens,
		Time: ollamaResp.FTotalTime,
	}

	json.NewEncoder(w).Encode(resp)
}
	

