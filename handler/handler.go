package handler

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
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

func (s *Shortener) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	domainCount := s.CalculateDomainStatistics()
	topDomains := getTopDomains(domainCount, 3)

	responseData := struct {
		TopDomains []string `json:"top_domains"`
	}{
		TopDomains: topDomains,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)
}

func (s *Shortener) CalculateDomainStatistics() map[string]int {
	domainCount := make(map[string]int)
	for _, originalURL := range s.urlMap {
		u, err := url.Parse(originalURL)
		if err != nil {
			continue
		}
		host := u.Host
		domainCount[host]++
	}
	return domainCount
}

func getTopDomains(domainCount map[string]int, n int) []string {
	var topDomains []string
	type domainStat struct {
		domain string
		count  int
	}
	var domainStats []domainStat
	for domain, count := range domainCount {
		domainStats = append(domainStats, domainStat{domain, count})
	}
	sort.Slice(domainStats, func(i, j int) bool {
		return domainStats[i].count > domainStats[j].count
	})
	for i := 0; i < n && i < len(domainStats); i++ {
		topDomains = append(topDomains, domainStats[i].domain)
	}
	return topDomains
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
