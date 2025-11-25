package tools

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/tmc/langchaingo/llms"
)

const (
	ServiceRotoReader  = "rotoreader"
	ServiceOddsTracker = "oddstracker"
)

// GetTools loads OpenAI function tools from OpenAPI specs of external services
// Falls back to hardcoded definitions if OpenAPI specs are unavailable
func GetTools() []llms.Tool {
	return GetToolsWithContext(context.Background())
}

// GetToolsWithContext loads tools with a custom context (useful for timeouts)
func GetToolsWithContext(ctx context.Context) []llms.Tool {
	// Get service URLs from environment or use defaults
	rotoreaderURL := os.Getenv("ROTOREADER_SERVICE_URL")
	if rotoreaderURL == "" {
		rotoreaderURL = "http://localhost:8081"
	}

	oddstrackerURL := os.Getenv("ODDSTRACKER_SERVICE_URL")
	if oddstrackerURL == "" {
		oddstrackerURL = "http://localhost:8080"
	}

	// Try to load OpenAPI specs from services
	sources := []SpecSource{
		{Service: ServiceRotoReader, URL: fmt.Sprintf("%s/openapi.json", rotoreaderURL)},
		{Service: ServiceOddsTracker, URL: fmt.Sprintf("%s/openapi.json", oddstrackerURL)},
	}

	specs, err := LoadMultipleSpecs(ctx, sources)
	if err != nil {
		log.Printf("Warning: Failed to load OpenAPI specs, falling back to hardcoded definitions: %v", err)
		return getFallbackTools()
	}

	// Convert OpenAPI specs to OpenAI function tools
	tools := ConvertOpenAPIToTools(specs)

	if len(tools) == 0 {
		log.Printf("Warning: No tools found in OpenAPI specs, falling back to hardcoded definitions")
		return getFallbackTools()
	}

	log.Printf("Successfully loaded %d tools from OpenAPI specs", len(tools))
	return tools
}

// getFallbackTools returns hardcoded tool definitions as a fallback
func getFallbackTools() []llms.Tool {
	resetToolMetadata()

	registerToolMetadata("get_roto_data", ToolMetadata{
		Service: ServiceRotoReader,
		Method:  http.MethodGet,
		Path:    "/feed",
	})
	registerToolMetadata("get_odds_data", ToolMetadata{
		Service: ServiceOddsTracker,
		Method:  http.MethodGet,
		Path:    "/changes",
	})

	return []llms.Tool{
		{
			Type: "function",
			Function: &llms.FunctionDefinition{
				Name:        "get_roto_data",
				Description: "Get the latest sports news feed from rotoreader",
				Parameters: map[string]any{
					"type":       "object",
					"properties": map[string]any{},
				},
			},
		},
		{
			Type: "function",
			Function: &llms.FunctionDefinition{
				Name:        "get_odds_data",
				Description: "Get recent betting odds changes from oddstracker",
				Parameters: map[string]any{
					"type":       "object",
					"properties": map[string]any{},
				},
			},
		},
	}
}
