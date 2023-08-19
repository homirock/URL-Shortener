package handler

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
)

const (
	base62Chars    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	shortURLLength = 7
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

func (s *Shortener) RedirectionHandler(w http.ResponseWriter, r *http.Request) {
	shortURL := r.URL.Path[len("/r/"):]
	originalURL, exists := s.urlMap[shortURL]
	if !exists {
		http.Error(w, "Short URL not found", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, originalURL, http.StatusSeeOther)
}

func generateShortURL() string {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	shortURL := make([]byte, shortURLLength)
	for i := 0; i < shortURLLength; i++ {
		shortURL[i] = base62Chars[r.Intn(len(base62Chars))]
	}
	return string(shortURL)
}
