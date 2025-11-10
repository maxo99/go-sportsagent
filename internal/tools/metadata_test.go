package tools

import (
	"context"
	"net/http"
	"testing"
)

func TestBuildHTTPRequest_PathAndQuery(t *testing.T) {
	metadata := ToolMetadata{
		Service: ServiceOddsTracker,
		Method:  http.MethodGet,
		Path:    "/event/{event_id}",
		PathParams: []ParameterDefinition{
			{Name: "event_id", In: ParameterInPath, Required: true},
		},
		QueryParams: []ParameterDefinition{
			{Name: "range", In: ParameterInQuery},
		},
	}

	args := map[string]interface{}{
		"event_id": "123",
		"range":    "week",
	}

	req, err := BuildHTTPRequest(context.Background(), "http://example.com", metadata, args)
	if err != nil {
		t.Fatalf("BuildHTTPRequest returned error: %v", err)
	}

	if req.Method != http.MethodGet {
		t.Fatalf("expected method GET, got %s", req.Method)
	}

	expectedURL := "http://example.com/event/123?range=week"
	if req.URL.String() != expectedURL {
		t.Fatalf("expected URL %s, got %s", expectedURL, req.URL.String())
	}
}

func TestBuildHTTPRequest_MissingPathParam(t *testing.T) {
	metadata := ToolMetadata{
		Service: ServiceOddsTracker,
		Method:  http.MethodGet,
		Path:    "/event/{event_id}",
		PathParams: []ParameterDefinition{
			{Name: "event_id", In: ParameterInPath, Required: true},
		},
	}

	_, err := BuildHTTPRequest(context.Background(), "http://example.com", metadata, map[string]interface{}{})
	if err == nil {
		t.Fatal("expected error for missing path parameter, got nil")
	}
}
