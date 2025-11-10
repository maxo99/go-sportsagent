package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type ParameterLocation string

const (
	ParameterInPath   ParameterLocation = "path"
	ParameterInQuery  ParameterLocation = "query"
	ParameterInHeader ParameterLocation = "header"
)

type ParameterDefinition struct {
	Name     string
	In       ParameterLocation
	Required bool
}

type ToolMetadata struct {
	Service     string
	Method      string
	Path        string
	PathParams  []ParameterDefinition
	QueryParams []ParameterDefinition
	HasJSONBody bool
}

var (
	toolMetadataMu  sync.RWMutex
	toolMetadataMap = map[string]ToolMetadata{}
)

func resetToolMetadata() {
	toolMetadataMu.Lock()
	defer toolMetadataMu.Unlock()

	toolMetadataMap = map[string]ToolMetadata{}
}

func registerToolMetadata(operationID string, metadata ToolMetadata) {
	if operationID == "" {
		return
	}

	toolMetadataMu.Lock()
	toolMetadataMap[operationID] = metadata
	toolMetadataMu.Unlock()
}

func GetToolMetadata(operationID string) (ToolMetadata, bool) {
	toolMetadataMu.RLock()
	metadata, ok := toolMetadataMap[operationID]
	toolMetadataMu.RUnlock()
	return metadata, ok
}

func GetToolService(operationID string) (string, bool) {
	metadata, ok := GetToolMetadata(operationID)
	if !ok {
		return "", false
	}
	return metadata.Service, true
}

func BuildHTTPRequest(ctx context.Context, baseURL string, metadata ToolMetadata, args map[string]interface{}) (*http.Request, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("baseURL cannot be empty")
	}

	remaining := map[string]interface{}{}
	for k, v := range args {
		remaining[k] = v
	}

	path := metadata.Path
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	for _, param := range metadata.PathParams {
		raw, ok := remaining[param.Name]
		if !ok {
			if param.Required {
				return nil, fmt.Errorf("missing required path parameter %s", param.Name)
			}
			continue
		}

		value := fmt.Sprintf("%v", raw)
		placeholder := fmt.Sprintf("{%s}", param.Name)
		path = strings.ReplaceAll(path, placeholder, url.PathEscape(value))
		delete(remaining, param.Name)
	}

	if strings.Contains(path, "{") {
		return nil, fmt.Errorf("unresolved path parameters in %s", path)
	}

	trimmedBase := strings.TrimRight(baseURL, "/")
	fullURL := trimmedBase + path

	parsedURL, err := url.Parse(fullURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL %s: %w", fullURL, err)
	}

	query := parsedURL.Query()
	for _, param := range metadata.QueryParams {
		raw, ok := remaining[param.Name]
		if !ok {
			if param.Required {
				return nil, fmt.Errorf("missing required query parameter %s", param.Name)
			}
			continue
		}

		switch v := raw.(type) {
		case []string:
			for _, item := range v {
				query.Add(param.Name, item)
			}
		case []interface{}:
			for _, item := range v {
				query.Add(param.Name, fmt.Sprintf("%v", item))
			}
		default:
			query.Set(param.Name, fmt.Sprintf("%v", raw))
		}

		delete(remaining, param.Name)
	}
	parsedURL.RawQuery = query.Encode()

	var body io.Reader
	if metadata.HasJSONBody && len(remaining) > 0 {
		bodyBytes, err := json.Marshal(remaining)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		body = bytes.NewReader(bodyBytes)
	}

	method := strings.ToUpper(metadata.Method)
	req, err := http.NewRequestWithContext(ctx, method, parsedURL.String(), body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}
