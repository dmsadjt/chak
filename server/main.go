package main

import (
	"chak-server/internal/embedding"
	"chak-server/internal/handler"
	"chak-server/internal/memory"
	"chak-server/internal/middleware"
	"chak-server/internal/ollama"
	"chak-server/internal/prompt"
	"chak-server/internal/search"
	"fmt"
	"log"
	"net/http"
	"os"

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

	szApiKey := os.Getenv("BRAVE_API_KEY")
	searchManager := search.NewBraveManager(szApiKey)
	promptManager := prompt.NewPromptManager()
	ollamaManager := ollama.NewDefaultOllamaManager("http://localhost:11434")
	embeddingManager := embedding.NewOllamaEmbedding("all-minilm:33m","http://localhost:11434")
	memoryManager := memory.NewMemoryManager(embeddingManager, "memory.json")

	chatManager := handler.NewChatHandlerManager(searchManager, promptManager, ollamaManager, memoryManager)

	chatHandler := http.HandlerFunc(chatManager.HandleChat)
	logMiddleware := &middleware.LoggerMiddleware{}
	corsMiddleware := &middleware.CorsMiddleware{}

	chainedChatHandler := Chain(chatHandler, logMiddleware, corsMiddleware)

	http.Handle("/", http.HandlerFunc(handleHome))
	http.Handle("/chat", chainedChatHandler)

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

