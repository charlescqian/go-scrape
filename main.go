package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/charlescqian/go-scrape/scraper"
)

type scrapeRequest struct {
	URL string `json:"url"`
}

type scrapeResponse struct {
	Content string `json:"content"`
}

// Handler for POST /scrape
func scrapeHandler(w http.ResponseWriter, r *http.Request) {
	// Check that it's a POST request
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode the JSON body into a Go struct
	var req scrapeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	content, err := scraper.ExtractTextFromURL((req.URL))
	if err != nil {
		http.Error(w, "Scraping failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(scrapeResponse{Content: content})
}

func main() {
	http.HandleFunc("/scrape", scrapeHandler)
	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
