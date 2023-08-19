package api

import (
	"fmt"
	"net/http"

	url "github.com/homirock/URL-Shortener/handler"
)
func Start(){
	shortener := url.NewShortener()
	http.HandleFunc("/shorten", shortener.ShortenHandler)
	http.HandleFunc("/r/", shortener.RedirectionHandler)
	port := 8084
	fmt.Printf("Starting server on port %d...\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}