package main

import (
	"net/http"
	"wordsearch/internal/handler"
)

func main() {
	router := http.NewServeMux()

	router.HandleFunc("GET /files/search", handler.Search) // Слово для поиска передается в query параметре "word"

	server := http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: router,
	}

	server.ListenAndServe()
}
