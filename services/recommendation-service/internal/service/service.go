package service

import (
	"context"
	"strconv"

	"github.com/freedom-music/recommendation-service/internal/client"
	"github.com/freedom-music/recommendation-service/internal/repository"
)

type Service struct {
	repo    *repository.Repository
	catalog *client.CatalogClient
}

type TrackView struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	DurationSec int      `json:"duration_sec"`
	AlbumID     string   `json:"album_id"`
	ArtistIDs   []string `json:"artist_ids"`
}

func New(repo *repository.Repository, catalog *client.CatalogClient) *Service {
	return &Service{repo: repo, catalog: catalog}
}

func (s *Service) Trending(ctx context.Context, limit int) ([]TrackView, error) {
	ids, err := s.repo.Trending(ctx, limit)
	// fall back to catalog seed on error or no events
	if err != nil || len(ids) == 0 {
		return s.fallbackCatalog(limit), nil
	}
	return s.fetchTracks(ids), nil
}

func (s *Service) RecentlyPlayed(ctx context.Context, userID int64, limit int) ([]TrackView, error) {
	ids, err := s.repo.RecentlyPlayed(ctx, userID, limit)
	if err != nil || len(ids) == 0 {
		return []TrackView{}, nil
	}
	return s.fetchTracks(ids), nil
}

func (s *Service) ForYou(ctx context.Context, userID int64, limit int) ([]TrackView, error) {
	liked, _ := s.repo.LikedTracks(ctx, userID, limit)
	trending, _ := s.repo.Trending(ctx, limit)
	ids := unique(liked, trending)
	if len(ids) > limit {
		ids = ids[:limit]
	}
	if len(ids) == 0 {
		return s.fallbackCatalog(limit), nil
	}
	return s.fetchTracks(ids), nil
}

func unique(a, b []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, id := range a {
		if !seen[id] {
			seen[id] = true
			out = append(out, id)
		}
	}
	for _, id := range b {
		if !seen[id] {
			seen[id] = true
			out = append(out, id)
		}
	}
	return out
}

func (s *Service) fetchTracks(ids []string) []TrackView {
	var out []TrackView
	for _, id := range ids {
		t, err := s.catalog.Track(id)
		if err != nil || t == nil {
			continue
		}
		out = append(out, TrackView{
			ID:          t.ID,
			Title:       t.Title,
			DurationSec: t.DurationSec,
			AlbumID:     t.AlbumID,
			ArtistIDs:   t.ArtistIDs,
		})
	}
	return out
}

func (s *Service) fallbackCatalog(limit int) []TrackView {
	tracks, err := s.catalog.Tracks()
	if err != nil || len(tracks) == 0 {
		return []TrackView{}
	}
	if len(tracks) > limit {
		tracks = tracks[:limit]
	}
	var out []TrackView
	for _, t := range tracks {
		out = append(out, TrackView{
			ID:          t.ID,
			Title:       t.Title,
			DurationSec: t.DurationSec,
			AlbumID:     t.AlbumID,
			ArtistIDs:   t.ArtistIDs,
		})
	}
	return out
}

// Parse user id string to int64 helper
func ParseUserID(val string) (int64, bool) {
	id, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, false
	}
	return id, true
}
