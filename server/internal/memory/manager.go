package memory

import (
	"chak-server/internal/embedding"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"sync"
	"time"
)

type MemoryManager struct {
	embedder embedding.EmbeddingInterface
	memories []MemoryEntry
	mu sync.RWMutex
	szFilename string
}

type scoredMemory struct {
	memory MemoryEntry
	score float64
}

func NewMemoryManager(embedder embedding.EmbeddingInterface, szFilename string) *MemoryManager {
	manager := &MemoryManager{
		embedder: embedder,
		memories: []MemoryEntry{},
		szFilename: szFilename,
	}
	
	manager.LoadFromFile()
	
	return manager
}

func (memoryMgr *MemoryManager) Reload(szFilename string) error {
	memoryMgr.mu.Lock()
	defer memoryMgr.mu.Unlock()

	log.Printf("Reloading memory from: %s", szFilename)

	memoryMgr.szFilename = szFilename
	memoryMgr.memories = []MemoryEntry{}

	data, err := os.ReadFile(szFilename)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("No existing memory file, starting fresh\n")
			return nil
		}
		return err
	}

	err = json.Unmarshal(data, &memoryMgr.memories)
	if err != nil {
		return err
	}

	log.Printf("Reloaded %d memories from %s\n", len(memoryMgr.memories), szFilename)
	return nil
}

func (memoryMgr *MemoryManager) LoadFromFile() error {
	data, err := os.ReadFile(memoryMgr.szFilename)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("No existing memory file, starting fresh \n")
			return nil
		}
		return err
	}

	memoryMgr.mu.Lock()
	defer memoryMgr.mu.Unlock()

	err = json.Unmarshal(data, &memoryMgr.memories)	
	if err != nil {
		return err
	}

	fmt.Printf("Loaded %d memories from %s\n", len(memoryMgr.memories), memoryMgr.szFilename)
	return nil
}

func (memoryMgr *MemoryManager) SaveToFile() error {
	memoryMgr.mu.RLock()
	defer memoryMgr.mu.RUnlock()

	data, err := json.MarshalIndent(memoryMgr.memories, "", "  ")	
	if err != nil {
		return err
	}

	return os.WriteFile(memoryMgr.szFilename, data, 0644)
}

func (memoryMgr *MemoryManager) SaveMemory(ctx context.Context, szText string, metadataMap map[string]string) error {
	vector, err := memoryMgr.embedder.EmbedText(ctx, szText)
	if err != nil {
		fmt.Printf("ERROR Embedding: %v\n", err)
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

	err = memoryMgr.SaveToFile()
	if err != nil {
		fmt.Printf("Error saving to file: %v\n", err)
	} else {
		fmt.Printf("Saved to file %s\n\n", memoryMgr.szFilename)
	}

	return nil
}

func (memoryMgr *MemoryManager) RetrieveRelevantContext(ctx context.Context, szQuery string, iTopK int, szFilterType string) ([]MemoryEntry, error) {
	queryVector, err := memoryMgr.embedder.EmbedText(ctx, szQuery)
	if err != nil {
		return nil, err
	}

	memoryMgr.mu.RLock()
	scores := make([]scoredMemory, 0, len(memoryMgr.memories))

	for _, mem := range memoryMgr.memories {
		if szFilterType != "" && mem.MetadataMap["type"] != szFilterType {
			continue
		}

		similarity := cosineSimilarity(queryVector, mem.FlVector)

		preview := mem.SzContent
		if len(preview) > 60 {
			preview = preview[:60] + "..."
		}

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
		preview := results[i].SzContent
		if len(preview) > 80 {
			preview = preview[:80] + "..."
		}
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

func (memoryMgr *MemoryManager) DeleteMemoriesByMetadata(szKey string, szValue string) error {
	memoryMgr.mu.Lock()

	filteredMemoryList := make([]MemoryEntry, 0, len(memoryMgr.memories))
	for _, memory := range memoryMgr.memories {
		if memory.MetadataMap[szKey] != szValue {
			filteredMemoryList = append(filteredMemoryList, memory)
		}
	}

	memoryMgr.memories = filteredMemoryList
	memoryMgr.mu.Unlock()

	err := memoryMgr.SaveToFile()
	if err != nil {
		fmt.Printf("Error saving to file: %v\n", err)
	} else {
		fmt.Printf("Saved to file %s\n\n", memoryMgr.szFilename)
	}

	return nil
}
