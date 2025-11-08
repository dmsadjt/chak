package ollama

import (
	"fmt"
	"bytes"
	"encoding/json"
	"net/http"
)

type OllamaManager struct {
	szApiURL string
}

func NewDefaultOllamaManager(szApiURL string) *OllamaManager {
	return &OllamaManager{
		szApiURL: szApiURL,
	}
}

type OllamaRequest struct {
	SzModel string `json:"model"`
	SzPrompt string `json:"prompt"`
	BStream bool `json:"stream"`
}

type OllamaResponse struct {
	SzResponse string `json:"response"`
	ITotalDuration int64 `json:"total_duration"`
	IEvalCount int `json:"eval_count"`
	IPromptEvalCount int `json:"prompt_eval_count"`
}

func (ollamaMgr *OllamaManager) Generate(szModel string, szPrompt string) (GenerateResponse, error) {
	reqBody := OllamaRequest {
		SzModel: szModel,
		SzPrompt: szPrompt,
		BStream: false,
	}

	jsonData, _ := json.Marshal(reqBody)

	resp, err := http.Post(
		fmt.Sprintf("%s/api/generate", ollamaMgr.szApiURL),
		"application/json",
		bytes.NewBuffer(jsonData),	
	)
	if err != nil {
		return GenerateResponse{}, err
	}
	defer resp.Body.Close()

	var ollamaResp OllamaResponse
	json.NewDecoder(resp.Body).Decode(&ollamaResp)

	return GenerateResponse{
		SzResponse: ollamaResp.SzResponse,
		ITotalTokens: ollamaResp.IEvalCount + ollamaResp.IPromptEvalCount,
		FTotalTime: float64(ollamaResp.ITotalDuration) / 1e9,
	}, nil
}

