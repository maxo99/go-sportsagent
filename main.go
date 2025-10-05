package main

import (
	"log"
	"net/http"
	"sportsagent/internal/handlers"
	
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	
	mux := http.NewServeMux()
	
	handler := handlers.NewAgentHandler()
	mux.HandleFunc("/query", handler.HandleQuery)
	
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
