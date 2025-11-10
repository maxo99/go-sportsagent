package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sportsagent/internal/services"
)

type AgentHandler struct {
	agentService *services.AgentService
}

func NewAgentHandler() *AgentHandler {
	return &AgentHandler{
		agentService: services.NewAgentService(),
	}
}

type QueryRequest struct {
	Query string `json:"query"`
}

type QueryResponse struct {
	Response string `json:"response"`
}

func (h *AgentHandler) HandleQuery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req QueryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("Received query:", req.Query)
	response, err := h.agentService.ProcessQuery(r.Context(), req.Query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("Sending response:", response)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(QueryResponse{Response: response})
}
