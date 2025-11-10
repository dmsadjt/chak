package memory

import "context"

type MemoryInterface interface {
	SaveMemory(ctx context.Context, szText string, metadataMap map[string]string) error
	RetrieveRelevantContext(ctx context.Context, szQuery string, iTopK int) ([]MemoryEntry, error)
	LoadFromFile() error
	SaveToFile() error
}

type MemoryEntry struct {
	SzId string
	SzContent string
	FlVector []float32
	MetadataMap map[string]string
}
