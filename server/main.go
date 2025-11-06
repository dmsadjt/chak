package main

import (
	"fmt"
	"net/http"
	"chak-server/internal/search"
)

func main() {
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/search", handleSearch)

	fmt.Println("Server starting on :5000")
	http.ListenAndServe(":5000", nil)
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

	searchManager := search.NewDuckDuckGoManager()

	results, err := searchManager.Search(szQuery)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Search error: %v", err)))
		return 
	}

	response := fmt.Sprintf("Found %d results for: %s", len(results), szQuery)
	w.Write([]byte(response))
}

