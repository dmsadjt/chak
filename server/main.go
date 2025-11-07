package main

import (
	"os"
	"encoding/json"
	"fmt"
	"net/http"
	"log"
	"chak-server/internal/search"
)


func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func main() {

	// VERBOSE LOG CONFIGURATION: This ensures you see detailed errors.
	log.SetOutput(os.Stdout) 
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	http.HandleFunc("/", corsMiddleware(handleHome))
	http.HandleFunc("/search", corsMiddleware(handleSearch))

	fmt.Println("Server starting on :5000")
	if err := http.ListenAndServe(":5000", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	w.Write([] byte("Chak backend API"))
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	szQuery := r.URL.Query().Get("q")

	if szQuery == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing query parameter 'q'"))
		return
	}

	log.Printf("Attempting search for query: %s", szQuery)

	searchManager := search.NewDuckDuckGoManager()

	results, err := searchManager.Search(szQuery)
	if err != nil {
		log.Printf("FATAL search error: %+v for query: %s", err, szQuery)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Search error: %v (Check server logs for root cause)", err)))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

