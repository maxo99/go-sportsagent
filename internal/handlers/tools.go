package handlers

import (
	"encoding/json"
	"net/http"
	"sportsagent/internal/tools"
)

type ToolsHandler struct {
	tools []interface{}
}

func NewToolsHandler() *ToolsHandler {
	// Get the tools and convert to a serializable format
	rawTools := tools.GetTools()
	serializedTools := make([]interface{}, len(rawTools))

	for i, tool := range rawTools {
		// Convert to map for JSON serialization
		toolBytes, _ := json.Marshal(tool)
		var toolMap map[string]interface{}
		json.Unmarshal(toolBytes, &toolMap)
		serializedTools[i] = toolMap
	}

	return &ToolsHandler{
		tools: serializedTools,
	}
}

type ToolsResponse struct {
	Tools []interface{} `json:"tools"`
	Count int           `json:"count"`
}

func (h *ToolsHandler) HandleGetTools(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := ToolsResponse{
		Tools: h.tools,
		Count: len(h.tools),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
