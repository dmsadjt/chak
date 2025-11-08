package ollama

type OllamaInterface interface {
	Generate(szModel string, szPrompt string) (GenerateResponse, error)
}

type GenerateResponse struct {
	SzResponse string `json:"response"`
	ITotalTokens int `json:"total_tokens"`
	FTotalTime float64 `json:"total_time"`
}

