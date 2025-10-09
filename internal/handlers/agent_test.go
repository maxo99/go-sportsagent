package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load("../../.env")
}

func TestHandleQuery_Success(t *testing.T) {
	handler := NewAgentHandler()

	reqBody := QueryRequest{Query: "What's the latest sports news?"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/query", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleQuery(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp QueryResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Response == "" {
		t.Error("expected non-empty response")
	}

	t.Logf("Response: %s", resp.Response)
}
