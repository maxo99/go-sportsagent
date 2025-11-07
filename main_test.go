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

func TestToolsEndpoint(t *testing.T) {
	mux := setupServer()

	req := httptest.NewRequest(http.MethodGet, "/tools", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status code: %d", w.Code)
	}

	var payload map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if _, ok := payload["tools"]; !ok {
		t.Fatalf("missing 'tools' field in response")
	}

	if _, ok := payload["count"]; !ok {
		t.Fatalf("missing 'count' field in response")
	}

	count, ok := payload["count"].(float64)
	if !ok {
		t.Fatalf("count is not a number")
	}

	if count < 0 {
		t.Fatalf("count should be non-negative, got: %v", count)
	}

	fmt.Printf("Tools endpoint returned %v tools\n", count)
}
