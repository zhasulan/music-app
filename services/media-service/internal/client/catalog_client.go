package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type CatalogClient struct {
	baseURL string
	http    *http.Client
}

type trackResponse struct {
	Track any `json:"track"`
}

func NewCatalogClient(baseURL string) *CatalogClient {
	return &CatalogClient{baseURL: baseURL, http: &http.Client{Timeout: 3 * time.Second}}
}

func (c *CatalogClient) ValidateTrack(ctx context.Context, trackID string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/api/v1/catalog/tracks/%s", c.baseURL, trackID), nil)
	if err != nil {
		return err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("track not found")
	}
	if resp.StatusCode >= 300 {
		return fmt.Errorf("catalog error %d", resp.StatusCode)
	}
	var body trackResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return fmt.Errorf("invalid catalog response: %w", err)
	}
	if body.Track == nil {
		return fmt.Errorf("track not found")
	}
	return nil
}
