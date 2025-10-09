package tools

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
)

// LoadOpenAPISpec fetches and parses an OpenAPI spec from the given URL
func LoadOpenAPISpec(ctx context.Context, url string) (*openapi3.T, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch OpenAPI spec from %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d from %s", resp.StatusCode, url)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromData(body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse OpenAPI spec: %w", err)
	}

	return doc, nil
}

// LoadMultipleSpecs loads OpenAPI specs from multiple service URLs
func LoadMultipleSpecs(ctx context.Context, urls []string) ([]*openapi3.T, error) {
	specs := make([]*openapi3.T, 0, len(urls))

	for _, url := range urls {
		spec, err := LoadOpenAPISpec(ctx, url)
		if err != nil {
			return nil, fmt.Errorf("failed to load spec from %s: %w", url, err)
		}
		specs = append(specs, spec)
	}

	return specs, nil
}
