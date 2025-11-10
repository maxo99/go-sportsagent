package tools

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestConvertOpenAPIToToolsRegistersServices(t *testing.T) {
	oddsSpecPath := filepath.Join("..", "clients", "oddstracker", "testdata", "openapi.json")
	if info, err := os.Stat(oddsSpecPath); err != nil || info.Size() == 0 {
		t.Skip("oddstracker OpenAPI fixture not available")
	}

	absOddsPath, err := filepath.Abs(oddsSpecPath)
	if err != nil {
		t.Fatalf("failed to resolve odds spec path: %v", err)
	}

	ctx := context.Background()
	oddsSpec, err := LoadOpenAPISpec(ctx, "file://"+absOddsPath)
	if err != nil {
		t.Fatalf("failed to load odds spec: %v", err)
	}

	specs := []ServiceSpec{{Service: ServiceOddsTracker, Spec: oddsSpec}}

	tools := ConvertOpenAPIToTools(specs)
	if len(tools) == 0 {
		t.Fatal("expected tools to be generated")
	}

	metadata, ok := GetToolMetadata("collect_sportevents")
	if !ok {
		t.Fatal("expected metadata for collect_sportevents")
	}
	if metadata.Service != ServiceOddsTracker {
		t.Fatalf("expected collect_sportevents to map to odds tracker, got %s", metadata.Service)
	}
	if metadata.Path != "/collect" {
		t.Fatalf("expected collect_sportevents path to be /collect, got %s", metadata.Path)
	}
	if metadata.Method != "POST" {
		t.Fatalf("expected collect_sportevents method POST, got %s", metadata.Method)
	}

	for _, tool := range tools {
		fn := tool.GetFunction()
		if fn == nil {
			continue
		}

		switch fn.Name {
		case "metrics_metrics_get", "health_check":
			t.Fatalf("unexpected function %s returned in tool list", fn.Name)
		}

		if meta, ok := GetToolMetadata(fn.Name); !ok || meta.Service == "" {
			t.Fatalf("missing service mapping for function %s", fn.Name)
		}
	}
}
