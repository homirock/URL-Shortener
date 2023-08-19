package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func Test_getTopDomains(t *testing.T) {
	type args struct {
		domainCount map[string]int
		n           int
	}
	myArgs := args{
		domainCount: map[string]int{
			"example.com": 1,
			"google.com":  2,
			"github.com":  3,
			"golang.com":  4,
		},
		n: 3,
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{name: "Success Case", args: myArgs, want: []string{"golang.com", "github.com", "google.com"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTopDomains(tt.args.domainCount, tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getTopDomains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateShortURL(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{name: "success case", want: "Test123"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateShortURL(); len(got) != len(tt.want) {
				t.Errorf("generateShortURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateDomainStatistics(t *testing.T) {
	shortener := NewShortener()
	shortener.urlMap = map[string]string{
		"https://example.com/path":    "https://example.com",
		"https://example.com/page":    "https://example.com",
		"https://google.com/search":   "https://google.com",
		"https://github.com/repo":     "https://github.com",
		"https://example.org/about":   "https://example.org",
		"https://example.net/contact": "https://example.net",
	}

	domainStats := shortener.CalculateDomainStatistics()

	expectedCounts := map[string]int{
		"example.com": 2,
		"google.com":  1,
		"github.com":  1,
		"example.org": 1,
		"example.net": 1,
	}

	for domain, expectedCount := range expectedCounts {
		actualCount, exists := domainStats[domain]
		if !exists {
			t.Errorf("Domain %s not found in domainStats", domain)
		}
		if actualCount != expectedCount {
			t.Errorf("Domain %s: Got count %d, expected %d", domain, actualCount, expectedCount)
		}
	}

	if len(domainStats) != len(expectedCounts) {
		t.Errorf("Mismatch in the number of domains. Got %d, expected %d", len(domainStats), len(expectedCounts))
	}
}

func TestMetricsHandler(t *testing.T) {
	// Create a Shortener instance
	shortener := NewShortener()
	// Populate urlMap (you can modify this according to your use case)
	shortener.urlMap = map[string]string{
		"https://example.com/path":  "https://example.com",
		"https://google.com/search": "https://google.com",
		"https://github.com/repo":   "https://github.com",
		"https://example.org/about": "https://example.org",
	}

	// Create a mock request
	req, err := http.NewRequest("GET", "/metrics", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	handler := http.HandlerFunc(shortener.MetricsHandler)
	handler.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	// Check the response content type
	expectedContentType := "application/json"
	if ct := rr.Header().Get("Content-Type"); ct != expectedContentType {
		t.Errorf("Handler returned wrong content type: got %v, want %v", ct, expectedContentType)
	}

	// Decode the response body
	var responseData struct {
		TopDomains []string `json:"top_domains"`
	}
	err = json.NewDecoder(rr.Body).Decode(&responseData)
	if err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}

	// Define expected top domains (you can adjust this based on your use case)
	expectedTopDomains := []string{"example.com", "google.com", "github.com"}

	// Compare top domains
	if len(responseData.TopDomains) != len(expectedTopDomains) {
		t.Errorf("Mismatch in the number of top domains: got %d, want %d", len(responseData.TopDomains), len(expectedTopDomains))
	}

	for i, domain := range responseData.TopDomains {
		if domain != expectedTopDomains[i] {
			t.Errorf("Top domain mismatch at index %d: got %s, want %s", i, domain, expectedTopDomains[i])
		}
	}
}

func TestShortenHandler(t *testing.T) {
	// Create a Shortener instance
	shortener := NewShortener()
	// Mocked request data
	requestData := struct {
		URL string `json:"url"`
	}{
		URL: "https://example.com",
	}
	// Convert requestData to JSON
	requestBody, err := json.Marshal(requestData)
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock POST request
	req, err := http.NewRequest("POST", "/shorten", bytes.NewReader(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	handler := http.HandlerFunc(shortener.ShortenHandler)
	handler.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	// Check the response content type
	expectedContentType := "application/json"
	if ct := rr.Header().Get("Content-Type"); ct != expectedContentType {
		t.Errorf("Handler returned wrong content type: got %v, want %v", ct, expectedContentType)
	}

	// Decode the response body
	var responseData struct {
		ShortURL string `json:"short_url"`
	}
	err = json.NewDecoder(rr.Body).Decode(&responseData)
	if err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}
}
