package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadOpenAPISpecFromFile(t *testing.T) {
	fixturePath := filepath.Join("..", "clients", "oddstracker", "testdata", "openapi.json")
	if info, err := os.Stat(fixturePath); err != nil || info.Size() == 0 {
		t.Skip("oddstracker OpenAPI fixture not available")
	}

	absPath, err := filepath.Abs(fixturePath)
	if err != nil {
		t.Fatalf("failed to resolve fixture path: %v", err)
	}

	spec, err := LoadOpenAPISpec(context.Background(), "file://"+absPath)
	if err != nil {
		t.Fatalf("failed to load OpenAPI spec from file: %v", err)
	}

	if spec == nil {
		t.Fatal("expected non-nil spec")
	}

	if len(spec.Paths.Map()) == 0 {
		t.Fatal("expected spec to contain at least one path")
	}
}

func TestLoadOpenAPIFromServices(t *testing.T) {
	if os.Getenv("INTEGRATION_TESTS") != "1" {
		t.Skip("Skipping integration test - set INTEGRATION_TESTS=1 to run")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Test loading from oddstracker
	t.Run("LoadFromOddsTracker", func(t *testing.T) {
		url := "http://localhost:8082/openapi.json"
		spec, err := LoadOpenAPISpec(ctx, url)
		if err != nil {
			t.Fatalf("Failed to load OpenAPI spec from %s: %v", url, err)
		}

		if spec == nil {
			t.Fatal("Spec is nil")
		}

		t.Logf("✓ Successfully loaded OpenAPI spec from %s", url)
		t.Logf("  Title: %s", spec.Info.Title)
		t.Logf("  Version: %s", spec.Info.Version)
		t.Logf("  Description: %s", spec.Info.Description)
		t.Logf("  Paths: %d", len(spec.Paths.Map()))

		// List all operations
		t.Log("\n  Operations:")
		for path, pathItem := range spec.Paths.Map() {
			for method, operation := range pathItem.Operations() {
				if operation != nil {
					t.Logf("    %s %s (operationId: %s)", method, path, operation.OperationID)
					if operation.Summary != "" {
						t.Logf("      Summary: %s", operation.Summary)
					}
					if len(operation.Parameters) > 0 {
						t.Logf("      Parameters: %d", len(operation.Parameters))
						for _, param := range operation.Parameters {
							if param.Value != nil {
								t.Logf("        - %s (%s, required: %v)", param.Value.Name, param.Value.In, param.Value.Required)
							}
						}
					}
					if operation.RequestBody != nil && operation.RequestBody.Value != nil {
						t.Log("      Request Body: application/json")
					}
				}
			}
		}
	})

	// Test loading from rotoreader
	t.Run("LoadFromRotoReader", func(t *testing.T) {
		url := "http://localhost:8081/openapi.json"
		spec, err := LoadOpenAPISpec(ctx, url)
		if err != nil {
			t.Fatalf("Failed to load OpenAPI spec from %s: %v", url, err)
		}

		if spec == nil {
			t.Fatal("Spec is nil")
		}

		t.Logf("✓ Successfully loaded OpenAPI spec from %s", url)
		t.Logf("  Title: %s", spec.Info.Title)
		t.Logf("  Version: %s", spec.Info.Version)
		t.Logf("  Paths: %d", len(spec.Paths.Map()))

		// List all operations
		t.Log("\n  Operations:")
		for path, pathItem := range spec.Paths.Map() {
			for method, operation := range pathItem.Operations() {
				if operation != nil {
					t.Logf("    %s %s (operationId: %s)", method, path, operation.OperationID)
					if operation.Summary != "" {
						t.Logf("      Summary: %s", operation.Summary)
					}
					if len(operation.Parameters) > 0 {
						t.Logf("      Parameters: %d", len(operation.Parameters))
						for _, param := range operation.Parameters {
							if param.Value != nil {
								t.Logf("        - %s (%s, required: %v)", param.Value.Name, param.Value.In, param.Value.Required)
							}
						}
					}
					if operation.RequestBody != nil && operation.RequestBody.Value != nil {
						t.Log("      Request Body: application/json")
					}
				}
			}
		}
	})
}

func TestConvertOpenAPIToOpenAITools(t *testing.T) {
	if os.Getenv("INTEGRATION_TESTS") != "1" {
		t.Skip("Skipping integration test - set INTEGRATION_TESTS=1 to run")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Load specs from both services
	sources := []SpecSource{
		{Service: ServiceOddsTracker, URL: "http://localhost:8082/openapi.json"},
		{Service: ServiceRotoReader, URL: "http://localhost:8081/openapi.json"},
	}

	specs, err := LoadMultipleSpecs(ctx, sources)
	if err != nil {
		t.Fatalf("Failed to load OpenAPI specs: %v", err)
	}

	// Convert to OpenAI tools
	tools := ConvertOpenAPIToTools(specs)

	if len(tools) == 0 {
		t.Fatal("No tools generated from OpenAPI specs")
	}

	t.Logf("\n✓ Generated %d OpenAI function tools from OpenAPI specs\n", len(tools))

	// Display each tool in detail
	for i, tool := range tools {
		// Extract function details
		toolJSON, _ := json.MarshalIndent(tool, "", "  ")

		t.Logf("Tool %d:\n%s\n", i+1, string(toolJSON))
	}

	// Verify tool structure
	t.Run("VerifyToolStructure", func(t *testing.T) {
		for _, tool := range tools {
			// Basic validation - tools should have proper structure
			toolJSON, err := json.Marshal(tool)
			if err != nil {
				t.Errorf("Failed to marshal tool: %v", err)
				continue
			}

			var toolMap map[string]interface{}
			if err := json.Unmarshal(toolJSON, &toolMap); err != nil {
				t.Errorf("Failed to unmarshal tool: %v", err)
				continue
			}

			// Check for required fields
			if toolMap["type"] != "function" {
				t.Errorf("Tool type should be 'function', got: %v", toolMap["type"])
			}

			functionData, ok := toolMap["function"].(map[string]interface{})
			if !ok {
				t.Error("Tool should have 'function' field")
				continue
			}

			if functionData["name"] == nil || functionData["name"] == "" {
				t.Error("Function should have non-empty 'name'")
			}

			if functionData["parameters"] == nil {
				t.Error("Function should have 'parameters' field")
			}
		}
	})
}

func TestGetTools(t *testing.T) {
	if os.Getenv("INTEGRATION_TESTS") != "1" {
		t.Skip("Skipping integration test - set INTEGRATION_TESTS=1 to run")
	}

	// Test the main GetTools function
	tools := GetTools()

	if len(tools) == 0 {
		t.Fatal("GetTools returned no tools")
	}

	t.Logf("\n✓ GetTools() returned %d tools\n", len(tools))

	// Pretty print all tools
	for i, tool := range tools {
		toolJSON, _ := json.MarshalIndent(tool, "", "  ")
		fmt.Printf("\nTool %d:\n%s\n", i+1, string(toolJSON))
	}
}
