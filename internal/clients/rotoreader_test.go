package clients

import (
	"context"
	"os"
	"testing"
)

func TestRotoReaderClient_Integration(t *testing.T) {
	// if os.Getenv("INTEGRATION_TESTS") == "" {
	// 	t.Skip("skipping integration test: set INTEGRATION_TESTS=1 to run")
	// }

	client := NewRotoReaderClient()
	ctx := context.Background()

	result, err := client.GetFeeds(ctx, map[string]interface{}{})

	if err != nil {
		t.Fatalf("failed to get data from rotoreader: %v", err)
	}

	if result == "" {
		t.Error("expected non-empty result from rotoreader")
	}

	t.Logf("received data: %s", result)
}

func TestRotoReaderClient_URLConfiguration(t *testing.T) {
	tests := []struct {
		name    string
		envURL  string
		wantURL string
	}{
		{
			name:    "uses env var when set",
			envURL:  "http://custom:9000",
			wantURL: "http://custom:9000",
		},
		{
			name:    "uses default when env var empty",
			envURL:  "",
			wantURL: "http://localhost:8001",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envURL != "" {
				os.Setenv("ROTOREADER_URL", tt.envURL)
				defer os.Unsetenv("ROTOREADER_URL")
			} else {
				os.Unsetenv("ROTOREADER_URL")
			}

			client := NewRotoReaderClient()

			if client.baseURL != tt.wantURL {
				t.Errorf("got baseURL %s, want %s", client.baseURL, tt.wantURL)
			}
		})
	}
}
