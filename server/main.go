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
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
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
	profileHandler := handler.NewProfileHandler(configManager)

	chatHandler := http.HandlerFunc(chatManager.HandleChat)
	logMiddleware := &middleware.LoggerMiddleware{}
	corsMiddleware := &middleware.CorsMiddleware{}

	chainedChatHandler := Chain(chatHandler, logMiddleware, corsMiddleware)

	http.Handle("/", http.HandlerFunc(handleHome))
	http.Handle("/chat", chainedChatHandler)

	http.Handle("/profiles", Chain(
		http.HandlerFunc(profileHandler.HandleListProfile),
		logMiddleware, corsMiddleware,
	))

	http.Handle("/profile/active", Chain(
		http.HandlerFunc(profileHandler.GetActiveProfile),
		logMiddleware, corsMiddleware,
	))

	http.Handle("/profile/switch", Chain(
		http.HandlerFunc(profileHandler.HandleSwitchProfile),
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

