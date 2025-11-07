package main

import (
	"log"
	"net/http"
	"sportsagent/internal/handlers"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func setupServer() *http.ServeMux {
	mux := http.NewServeMux()
	handler := handlers.NewAgentHandler()
	toolsHandler := handlers.NewToolsHandler()
	mux.HandleFunc("/query", handler.HandleQuery)
	mux.HandleFunc("/tools", toolsHandler.HandleGetTools)
	mux.HandleFunc("/healthz", handlers.HandleHealth)
	mux.Handle("/metrics", promhttp.Handler())
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
