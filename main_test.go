package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestQueryEndpoint(t *testing.T) {
	t.Skip("requires OPENAI_API_KEY")

	mux := setupServer()

	reqBody := map[string]string{"query": "test"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/query", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Errorf("unexpected status code: %d", w.Code)
	}
}
