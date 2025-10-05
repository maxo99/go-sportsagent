package clients

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
)

type RotoReaderClient struct {
	baseURL string
	client  *http.Client
}

func NewRotoReaderClient() *RotoReaderClient {
	baseURL := os.Getenv("ROTOREADER_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8001"
	}

	return &RotoReaderClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func (c *RotoReaderClient) GetData(ctx context.Context, params map[string]interface{}) (string, error) {
	resp, err := c.client.Get(fmt.Sprintf("%s/api/data", c.baseURL))
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
