package tools

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
)

type SpecSource struct {
	Service string
	URL     string
}

type ServiceSpec struct {
	Service string
	Spec    *openapi3.T
}

// LoadOpenAPISpec fetches and parses an OpenAPI spec from the given URL
func LoadOpenAPISpec(ctx context.Context, location string) (*openapi3.T, error) {
	var data []byte
	var err error

	if strings.HasPrefix(location, "file://") || !strings.Contains(location, "://") {
		path := strings.TrimPrefix(location, "file://")
		data, err = os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read OpenAPI spec file %s: %w", path, err)
		}
	} else {
		client := &http.Client{
			Timeout: 5 * time.Second,
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, location, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch OpenAPI spec from %s: %w", location, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unexpected status code %d from %s", resp.StatusCode, location)
		}

		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}
	}

	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromData(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse OpenAPI spec: %w", err)
	}

	return doc, nil
}

// LoadMultipleSpecs loads OpenAPI specs from multiple service URLs
func LoadMultipleSpecs(ctx context.Context, sources []SpecSource) ([]ServiceSpec, error) {
	specs := make([]ServiceSpec, 0, len(sources))

	for _, source := range sources {
		spec, err := LoadOpenAPISpec(ctx, source.URL)
		if err != nil {
			return nil, fmt.Errorf("failed to load spec from %s: %w", source.URL, err)
		}
		specs = append(specs, ServiceSpec{Service: source.Service, Spec: spec})
	}

	return specs, nil
}
