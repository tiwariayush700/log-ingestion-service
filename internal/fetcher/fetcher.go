package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/tiwariayush700/log-ingestion-service/internal/models"
)

// Fetcher is responsible for retrieving data from external APIs
type Fetcher struct {
	endpoint string
	client   *http.Client
	timeout  time.Duration
}

// New creates a new Fetcher instance
func New(endpoint string) *Fetcher {
	return &Fetcher{
		endpoint: endpoint,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		timeout: 30 * time.Second,
	}
}

// FetchPosts retrieves posts from the API
func (f *Fetcher) FetchPosts(ctx context.Context) ([]models.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, f.timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, f.endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var posts []models.Post
	if err := json.Unmarshal(body, &posts); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return posts, nil
}
