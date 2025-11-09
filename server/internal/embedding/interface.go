package embedding

import "context"

type EmbeddingInterface interface {
	EmbedText(ctx context.Context, szText string) ([]float32, error)
}
