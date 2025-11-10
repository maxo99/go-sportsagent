package clients

import (
	"context"
	"os"
	"testing"

	"sportsagent/internal/tools"
)

func TestOddsTrackerClient_Integration(t *testing.T) {
	if os.Getenv("INTEGRATION_TESTS") == "" {
		t.Skip("skipping integration test: set INTEGRATION_TESTS=1 to run")
	}

	client := NewOddsTrackerClient()
	ctx := context.Background()
	tools.GetTools()

	result, err := client.ExecuteOperation(ctx, "get_odds_data", map[string]any{})

	if err != nil {
		t.Fatalf("failed to get data from oddstracker: %v", err)
	}

	if result == "" {
		t.Error("expected non-empty result from oddstracker")
	}

	t.Logf("received data: %s", result)
}

func TestOddsTrackerClient_URLConfiguration(t *testing.T) {
	tests := []struct {
		name    string
		envURL  string
		wantURL string
	}{
		{
			name:    "uses env var when set",
			envURL:  "http://custom:9001",
			wantURL: "http://custom:9001",
		},
		{
			name:    "uses default when env var empty",
			envURL:  "",
			wantURL: "http://localhost:8082",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envURL != "" {
				os.Setenv("ODDSTRACKER_SERVICE_URL", tt.envURL)
				defer os.Unsetenv("ODDSTRACKER_SERVICE_URL")
			} else {
				os.Unsetenv("ODDSTRACKER_SERVICE_URL")
			}

			client := NewOddsTrackerClient()

			if client.baseURL != tt.wantURL {
				t.Errorf("got baseURL %s, want %s", client.baseURL, tt.wantURL)
			}
		})
	}
}
