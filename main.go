package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Shortener struct {
	urlMap map[string]string
}

func NewShortener() *Shortener {
	return &Shortener{
		urlMap: make(map[string]string),
	}
}

func (s *Shortener) ShortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestData struct {
		URL string `json:"url"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	// Check if the URL already exists in the map
	shortURL, exists := s.urlMap[requestData.URL]
	if !exists {
		// Generate a short URL and store in the map
		shortURL = generateShortURL()
		s.urlMap[shortURL] = requestData.URL
	}

	responseData := struct {
		ShortURL string `json:"short_url"`
	}{
		ShortURL: shortURL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)
}

func main() {
	shortener := NewShortener()

	http.HandleFunc("/shorten", shortener.ShortenHandler)

	port := 8084
	fmt.Printf("Starting server on port %d...\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
