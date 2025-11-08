package handler

import (
	"chak-server/internal/prompt"
	"chak-server/internal/ollama"
	"chak-server/internal/search"
	"encoding/json"
	"net/http"
)

type ChatRequest struct {
    Query  string `json:"query"`
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
}

func NewChatHandlerManager(sm search.SearchInterface, pm prompt.PromptInterface, om ollama.OllamaInterface) *ChatHandlerManager {
	return &ChatHandlerManager{
		searchManager: sm,
		promptManager: pm,
		ollamaManager: om,
	}
}

func (chatManager *ChatHandlerManager) HandleChat(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var searchResultData []search.SearchResultData
	if req.Search {
		results, err := chatManager.searchManager.Search(req.Query)
		if err == nil {
			searchResultData = results
		}
	}

	szFinalPrompt := chatManager.promptManager.Build(req.Query, searchResultData)

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
	

