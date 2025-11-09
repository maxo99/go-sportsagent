package handlers

import (
	"encoding/json"
	"net/http"

	"sportsagent/internal/version"
)

type healthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version,omitempty"`
}

// HandleHealth returns a static response so infrastructure can verify the process is running.
func HandleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(healthResponse{Status: "ok", Version: version.Version})
}
