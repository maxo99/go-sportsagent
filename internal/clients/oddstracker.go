package clients

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"sportsagent/internal/tools"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type OddsTrackerClient struct {
	baseURL string
	client  *http.Client
}

func NewOddsTrackerClient() *OddsTrackerClient {
	baseURL := os.Getenv("ODDSTRACKER_SERVICE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8082"
	}

	return &OddsTrackerClient{
		baseURL: baseURL,
		client: &http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		},
	}
}

func (c *OddsTrackerClient) CallOperation(ctx context.Context, metadata tools.ToolMetadata, params map[string]interface{}) (string, error) {
	req, err := tools.BuildHTTPRequest(ctx, c.baseURL, metadata, params)
	if err != nil {
		return "", err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (c *OddsTrackerClient) ExecuteOperation(ctx context.Context, operationID string, params map[string]interface{}) (string, error) {
	metadata, ok := tools.GetToolMetadata(operationID)
	if !ok {
		return "", fmt.Errorf("no metadata registered for operation %s", operationID)
	}

	return c.CallOperation(ctx, metadata, params)
}
