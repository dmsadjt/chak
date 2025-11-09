package memory

import (
	"chak-server/internal/embedding"
	"context"
	"fmt"
	"math"
	"sync"
	"time"
)

type MemoryManager struct {
	embedder embedding.EmbeddingInterface
	memories []MemoryEntry
	mu sync.RWMutex
}

type scoredMemory struct {
	memory MemoryEntry
	score float64
}

func NewMemoryManager(embedder embedding.EmbeddingInterface) *MemoryManager {
	return &MemoryManager{
		embedder: embedder,
		memories: []MemoryEntry{},
	}
}

func (memoryMgr *MemoryManager) SaveMemory(ctx context.Context, szText string, metadataMap map[string]string) error {
	vector, err := memoryMgr.embedder.EmbedText(ctx, szText)
	if err != nil {
		return err
	}

	memoryEntry := MemoryEntry{
		SzId: generateID(),
		SzContent: szText,
		FlVector: vector,
		MetadataMap: metadataMap,
	}

	memoryMgr.mu.Lock()
	memoryMgr.memories = append(memoryMgr.memories, memoryEntry)
	memoryMgr.mu.Unlock()

	return nil
}

func (memoryMgr *MemoryManager) RetrieveRelevantContext(ctx context.Context, szQuery string, iTopK int) ([]MemoryEntry, error) {
	queryVector, err := memoryMgr.embedder.EmbedText(ctx, szQuery)
	if err != nil {
		return nil, err
	}

	memoryMgr.mu.RLock()
	scores := make([]scoredMemory, 0, len(memoryMgr.memories))

	for _, mem := range memoryMgr.memories {
		similarity := cosineSimilarity(queryVector, mem.FlVector)
		scores = append(scores, scoredMemory{
			memory: mem,
			score: similarity,
		})
	}

	memoryMgr.mu.RUnlock()

	sortByScore(scores)

	topK := iTopK
	if topK > len(scores) {
		topK = len(scores)
	}

	results := make([]MemoryEntry, topK)
	for i := 0; i < topK; i++ {
		results[i] = scores[i].memory
	}

	return results, nil
}

func cosineSimilarity(a, b []float32) float64 {
	var dotProduct, normA, normB float64

	for i := range a {
		dotProduct += float64(a[i] * b[i])
		normA += float64(a[i] * a[i])
		normB += float64(b[i] * b[i])
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}

func sortByScore(scores []scoredMemory) {
	for i := 0; i < len(scores)-1; i++ {
		for j := 0; j < len(scores)-i-1; j++ {
			if scores[j].score < scores[j+1].score {
				scores[j], scores[j+1] = scores [j+1], scores[j]
			}
		}
	}
}



func generateID() string {
	return fmt.Sprintf("mem_%d", time.Now().UnixNano())
}
