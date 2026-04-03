package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type CatalogClient struct {
	baseURL string
	client  *http.Client
}

type Track struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	DurationSec int      `json:"duration_sec"`
	AlbumID     string   `json:"album_id"`
	ArtistIDs   []string `json:"artist_ids"`
}

type Album struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Artist string `json:"artist"`
}

type Artist struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewCatalogClient(baseURL string) *CatalogClient {
	return &CatalogClient{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 5 * time.Second},
	}
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

func (c *CatalogClient) Albums() ([]Album, error) {
	var resp struct {
		Albums []Album `json:"albums"`
	}
	if err := c.get("/api/v1/catalog/albums", &resp); err != nil {
		return nil, err
	}
	return resp.Albums, nil
}

func (c *CatalogClient) Artists() ([]Artist, error) {
	var resp struct {
		Artists []Artist `json:"artists"`
	}
	if err := c.get("/api/v1/catalog/artists", &resp); err != nil {
		return nil, err
	}
	return resp.Artists, nil
}

func (c *CatalogClient) get(path string, out any) error {
	url := c.baseURL + path
	r, err := c.client.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("catalog status %d", r.StatusCode)
	}
	return json.NewDecoder(r.Body).Decode(out)
}
