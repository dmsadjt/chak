package main

import (
	"chak-server/internal/config"
	"chak-server/internal/embedding"
	"chak-server/internal/handler"
	"chak-server/internal/indexer"
	"chak-server/internal/memory"
	"chak-server/internal/middleware"
	"chak-server/internal/ollama"
	"chak-server/internal/prompt"
	"chak-server/internal/search"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

type AppManagers struct {
	configMgr config.ConfigInterface
	searchMgr search.SearchInterface
	promptMgr prompt.PromptInterface
	ollamaMgr ollama.OllamaInterface
	embedMgr embedding.EmbeddingInterface
	memoryMgr memory.MemoryInterface
	indexerMgr indexer.ManagerInterface
	chatMgr *handler.ChatHandlerManager
	mu sync.RWMutex
}

func (app *AppManagers) GetChatManager() *handler.ChatHandlerManager {
	app.mu.RLock()
	defer app.mu.RUnlock()
	return app.chatMgr
}

func (app *AppManagers) HotReloadProfile(szProfileName string) error {
	app.mu.Lock()
	defer app.mu.Unlock()

	log.Printf("Hot reloading profile: %s", szProfileName)
	newProfile, err := app.configMgr.GetProfile(szProfileName)
	if err != nil {
		return fmt.Errorf("failed to get profile: %w", err)
	}

	app.indexerMgr.StopWatcher()

	if err := app.memoryMgr.Reload(newProfile.SzMemoryFile); err != nil {
		return fmt.Errorf("failed to reload memory: %w", err)
	}

	newScanner := indexer.NewDirectoryScanner(
		newProfile.SzDirectories,
		newProfile.Extensions,
		newProfile.InMaxSizeFile,
	)
	app.indexerMgr = indexer.NewIndexerManager(newScanner, app.memoryMgr, newProfile.SzIndexFile)

	log.Println("Running initial indexing for new profile.")
	if err := app.indexerMgr.IndexAll(); err != nil {
		log.Printf("Indexing warning: %v", err)
	}

	app.chatMgr = handler.NewChatHandlerManager(
		app.searchMgr,
		app.promptMgr,
		app.ollamaMgr,
		app.memoryMgr,
	)

	log.Printf("Hot reload complete, current profile: %s", newProfile.SzName)

	return nil
}

func main() {
	log.SetOutput(os.Stdout) 
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	configManager := config.NewConfigManager("config.json")
	activeProfile := configManager.GetActiveProfile()

	log.Printf("Active profile: %s (%s)\n", activeProfile.SzName, activeProfile.SzDescription)

	szApiKey := os.Getenv("BRAVE_API_KEY")
	searchManager := search.NewBraveManager(szApiKey)
	promptManager := prompt.NewPromptManager()
	ollamaManager := ollama.NewDefaultOllamaManager("http://localhost:11434")
	embeddingManager := embedding.NewOllamaEmbedding("all-minilm:33m","http://localhost:11434")
	memoryManager := memory.NewMemoryManager(embeddingManager, activeProfile.SzMemoryFile)

	log.Println("Initalizing document indexer...")

	scannerMgr := indexer.NewDirectoryScanner(
		activeProfile.SzDirectories, 
		activeProfile.Extensions, 
		activeProfile.InMaxSizeFile)
	idxManager := indexer.NewIndexerManager(scannerMgr, memoryManager, activeProfile.SzIndexFile)

	log.Println("Running initial document indexing")
	if err := idxManager.IndexAll(); err != nil {
		log.Printf("Indexing failed: %v\n", err)
	}

	idxManager.StartWatcher(5 * time.Minute)

	chatManager := handler.NewChatHandlerManager(searchManager, promptManager, ollamaManager, memoryManager)

	appManagers := &AppManagers{
		configMgr: configManager,
		searchMgr: searchManager,
		promptMgr: promptManager,
		ollamaMgr: ollamaManager,
		embedMgr: embeddingManager,
		memoryMgr: memoryManager,
		indexerMgr: idxManager,
		chatMgr: chatManager,
	}


	logMiddleware := &middleware.LoggerMiddleware{}
	corsMiddleware := &middleware.CorsMiddleware{}

	http.Handle("/", http.HandlerFunc(handleHome))
	http.Handle("/chat", Chain(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			appManagers.GetChatManager().HandleChat(w, r)
		}), logMiddleware, corsMiddleware,
	))
	
	profileHandler := handler.NewProfileHandler(configManager)

	http.Handle("/profiles", Chain(
		http.HandlerFunc(profileHandler.HandleListProfile),
		logMiddleware, corsMiddleware,
	))


	http.Handle("/profile/active", Chain(
		http.HandlerFunc(profileHandler.GetActiveProfile),
		logMiddleware, corsMiddleware,
	))

	http.Handle("/profile/switch", Chain(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var req handler.SwitchProfileRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "Invalid request", http.StatusBadRequest)
				return
			}
			
			if err := configManager.SwitchProfile(req.SzProfileName); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if err := appManagers.HotReloadProfile(req.SzProfileName); err != nil {
				log.Printf("Hot reload failed: %v", err)
				http.Error(w, "Profile switched but reload failed", http.StatusInternalServerError)
				return
			}

			profile := configManager.GetActiveProfile()
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status": "success",
				"active_profile": profile,
				"hot_reloaded": true,
			})
		}),
		logMiddleware, corsMiddleware,
	))

	fmt.Println("Server starting on :5000")
	if err := http.ListenAndServe(":5000", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func Chain(handler http.Handler, mws ...middleware.Middleware) http.Handler {
	for i := len(mws) -1;  i >= 0; i-- {
		handler = mws[i].Handle(handler)
	}
	return handler
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	w.Write([] byte("Chak backend API"))
}

