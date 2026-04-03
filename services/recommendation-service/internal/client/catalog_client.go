package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Track struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	DurationSec int      `json:"duration_sec"`
	AlbumID     string   `json:"album_id"`
	ArtistIDs   []string `json:"artist_ids"`
}

type CatalogClient struct {
	baseURL string
	http    *http.Client
}

func NewCatalogClient(baseURL string) *CatalogClient {
	return &CatalogClient{
		baseURL: baseURL,
		http:    &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *CatalogClient) Track(id string) (*Track, error) {
	var resp struct {
		Track *Track `json:"track"`
	}
	if err := c.get(fmt.Sprintf("/api/v1/catalog/tracks/%s", id), &resp); err != nil {
		return nil, err
	}
	return resp.Track, nil
}

func (c *CatalogClient) get(path string, out any) error {
	url := c.baseURL + path
	r, err := c.http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("catalog status %d", r.StatusCode)
	}
	return json.NewDecoder(r.Body).Decode(out)
}

func (c *CatalogClient) Tracks() ([]Track, error) {
	var resp struct {
		Tracks []Track `json:"tracks"`
	}
	if err := c.get("/api/v1/catalog/tracks", &resp); err != nil {
		return nil, err
	}
	return resp.Tracks, nil
}
