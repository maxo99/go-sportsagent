package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestQueryEndpoint(t *testing.T) {

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
	fmt.Printf("w: %v\n", w.Body.String())
}

func TestHealthEndpoint(t *testing.T) {
	mux := setupServer()

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status code: %d", w.Code)
	}

	var payload map[string]string
	if err := json.NewDecoder(w.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if payload["status"] != "ok" {
		t.Fatalf("unexpected status value: %s", payload["status"])
	}
}
