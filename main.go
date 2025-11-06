package main

import (
	"log"
	"net/http"
	"sportsagent/internal/handlers"

	"github.com/joho/godotenv"
)

func setupServer() *http.ServeMux {
	mux := http.NewServeMux()
	handler := handlers.NewAgentHandler()
	mux.HandleFunc("/query", handler.HandleQuery)
	mux.HandleFunc("/healthz", handlers.HandleHealth)
	return mux
}

func main() {
	godotenv.Load()

	mux := setupServer()

	log.Println("Starting server on :8082")
	if err := http.ListenAndServe(":8082", mux); err != nil {
		log.Fatal(err)
	}
}
