package clients

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
)

type OddsTrackerClient struct {
	baseURL string
	client  *http.Client
}

func NewOddsTrackerClient() *OddsTrackerClient {
	baseURL := os.Getenv("ODDSTRACKER_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8000"
	}

	return &OddsTrackerClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func (c *OddsTrackerClient) GetData(ctx context.Context, params map[string]interface{}) (string, error) {
	resp, err := c.client.Get(fmt.Sprintf("%s/api/odds", c.baseURL))
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
