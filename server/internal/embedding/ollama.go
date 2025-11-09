package embedding

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type OllamaEmbedding struct {
	SzModel string
	SzHost string
	Client *http.Client
}

func NewOllamaEmbedding(szModel, szHost string) *OllamaEmbedding {
	return &OllamaEmbedding{
		SzModel: szModel,
		SzHost: szHost,
		Client: &http.Client{},
	}
}

type ollamaEmbedRequest struct {
	SzModel string `json:"model"`
	SzInput string `json:"input"`
}

type ollamaEmbedResponse struct {
	SzEmbedding []float32 `json:"embedding"`
}

func (ollamaEmbed *OllamaEmbedding) EmbedText(ctx context.Context, szText string) ([]float32, error) {
	reqBody, _ := json.Marshal(ollamaEmbedRequest{
		SzModel: ollamaEmbed.SzModel,
		SzInput: szText,
	})

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/api/embed", ollamaEmbed.SzHost), bytes.NewBuffer(reqBody))

	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := ollamaEmbed.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var res ollamaEmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return res.SzEmbedding, nil
}
