package indexer

import (
	"chak-server/internal/document"
	"chak-server/internal/memory"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type IndexedFile struct {
	SzPath string `json:"path"`
	SzHash string `json:"hash"`
	TmIndexedTime time.Time `json:"indexed_at`
}

type IndexerManager struct {
	scannerMgr ScannerInterface
	memoryMgr memory.MemoryInterface
	indexedFilesMap map[string]IndexedFile
	szIndexFilePath string
}

func NewIndexerManager(scannerMgr ScannerInterface, memoryMgr memory.MemoryInterface, szIndexFile string) *IndexerManager {
	idxMgr := &IndexerManager{
		scannerMgr: scannerMgr,
		memoryMgr: memoryMgr,
		indexedFilesMap: make(map[string]IndexedFile),
		szIndexFilePath: szIndexFile,
	}

	idxMgr.loadIndexState()

	return idxMgr
}

func (idxMgr *IndexerManager) StartWatcher(tmInterval time.Duration) {
	ticker := time.NewTicker(tmInterval)

	go func() {
		for range ticker.C {
			log.Println("Auto indexing check")
			if err := idxMgr.IndexAll(); err != nil {
				log.Printf("Auto indexing error: %v \n", err)
			}
		}
	}()

	log.Printf("Watcher watching every %v\n", tmInterval)
}

func (idxMgr *IndexerManager) IndexAll() error {
	log.Println("Scanning files...")

	files, err := idxMgr.scannerMgr.ScanDirectories()
	if err != nil {
		return fmt.Errorf("Error scanning files: %w", err)
	}

	log.Printf("Found %d files\n", len(files))

	if len(files) == 0 {
		log.Println("No files to index")
		return nil 
	}

	inIndexed := 0
	inSkipped := 0

	ctx := context.Background()

	for _, file := range files {
		if idxMgr.shouldIndex(file) {
			if err := idxMgr.indexFile(ctx, file); err != nil {
				log.Printf("Error indexing %s: %v", file.SzName, err)
				continue
			}
			inIndexed++
		} else {
			inSkipped++
		}
	}

	log.Printf("✅ Indexed: %d files, Skipped: %d files (unchanged)\n", inIndexed, inSkipped)
	
	if err := idxMgr.saveIndexState(); err != nil {
		log.Printf("⚠️  Warning: Failed to save index state: %v\n", err)
	}

	return nil
}

func (idxMgr *IndexerManager) saveIndexState() error {
	data, err := json.MarshalIndent(idxMgr.indexedFilesMap, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(idxMgr.szIndexFilePath, data, 0644)
}

func (idxMgr *IndexerManager) indexFile(ctx context.Context, file FileInfo) error {
	log.Printf("Indexing %s\n", file.SzName)

	content, err := os.ReadFile(file.SzPath)
	if err != nil {
		return fmt.Errorf("Failed to read file: %w", err)
	}

	text := string(content)
	log.Printf("    Read %d bytes from file\n", len(text))


	chunks := document.ChunkText(text, 500)
	log.Printf("     Chunked into %d pieces\n", len(chunks))

	if len(chunks) == 0 {
		log.Printf("No chunks created")
		return nil
	}

	for i, chunk := range chunks {
		log.Printf("   💾 Saving chunk %d/%d (length: %d)\n", i+1, len(chunks), len(chunk))
		metadata := map[string]string {
			"type":         "document",
			"source":       "filesystem",
			"filepath":     file.SzPath,
			"filename":     file.SzName,
			"extension":    file.SzExtension,
			"chunk_id":     fmt.Sprintf("%d", i),
			"total_chunks": fmt.Sprintf("%d", len(chunks)),
			"indexed_at":   time.Now().Format(time.RFC3339),
		}

		err := idxMgr.memoryMgr.SaveMemory(ctx, chunk, metadata)
		if err != nil {
			log.Printf("Error saving memory %d: %v\n", i, err)
			return fmt.Errorf("Error saving memory %d: %w", i, err)
		}
		log.Println("Successfully saved memory")
	}

	idxMgr.indexedFilesMap[file.SzPath] = IndexedFile{
		SzPath: file.SzPath,
		SzHash: file.SzHash,
		TmIndexedTime: time.Now(),
	}
	
	return nil 
}

func (idxMgr *IndexerManager) shouldIndex(file FileInfo) bool {
	existing, exists := idxMgr.indexedFilesMap[file.SzPath]

	if !exists {
		return true
	}

	if existing.SzHash != file.SzHash {
		return true
	}

	return false
}

func (idxMgr *IndexerManager) loadIndexState() error {
	data, err := os.ReadFile(idxMgr.szIndexFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("No Index Found, starting fresh")
			return nil
		}
		return err
	}

	if err := json.Unmarshal(data, &idxMgr.indexedFilesMap); err != nil {
		return err
	}

	log.Printf("Index loaded: %d files were previously indexed.", len(idxMgr.indexedFilesMap))
	return nil
}


