package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestLoadOpenAPIFromServices(t *testing.T) {
	// Only run if INTEGRATION_TESTS is set
	if os.Getenv("INTEGRATION_TESTS") != "1" {
		t.Skip("Skipping integration test - set INTEGRATION_TESTS=1 to run")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Test loading from oddstracker
	t.Run("LoadFromOddsTracker", func(t *testing.T) {
		url := "http://localhost:8000/openapi.json"
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
		url := "http://localhost:8001/openapi.json"
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
	// Only run if INTEGRATION_TESTS is set
	if os.Getenv("INTEGRATION_TESTS") != "1" {
		t.Skip("Skipping integration test - set INTEGRATION_TESTS=1 to run")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Load specs from both services
	urls := []string{
		"http://localhost:8000/openapi.json",
		"http://localhost:8001/openapi.json",
	}

	specs, err := LoadMultipleSpecs(ctx, urls)
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
	// Only run if INTEGRATION_TESTS is set
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
