package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/freedom-music/playlist-service/internal/client"
	"github.com/freedom-music/playlist-service/internal/domain"
	"github.com/freedom-music/playlist-service/internal/repository"
)

var (
	ErrPlaylistNotFound = errors.New("playlist not found")
)

type PlaylistService struct {
	repo    *repository.PlaylistRepository
	catalog *client.CatalogClient
}

func NewPlaylistService(repo *repository.PlaylistRepository, catalog *client.CatalogClient) *PlaylistService {
	return &PlaylistService{repo: repo, catalog: catalog}
}

func (s *PlaylistService) Create(ctx context.Context, userID int64, name, description string) (*domain.Playlist, error) {
	return s.repo.Create(ctx, userID, name, description)
}

func (s *PlaylistService) List(ctx context.Context, userID int64) ([]domain.Playlist, error) {
	return s.repo.List(ctx, userID)
}

func (s *PlaylistService) Get(ctx context.Context, userID, playlistID int64) (*domain.Playlist, []domain.PlaylistTrack, error) {
	pl, err := s.repo.Get(ctx, userID, playlistID)
	if err != nil {
		return nil, nil, err
	}
	if pl == nil {
		return nil, nil, ErrPlaylistNotFound
	}
	tracks, err := s.repo.ListTracks(ctx, playlistID)
	if err != nil {
		return nil, nil, err
	}
	return pl, tracks, nil
}

func (s *PlaylistService) Update(ctx context.Context, userID, playlistID int64, name, description string) (*domain.Playlist, error) {
	pl, err := s.repo.Update(ctx, userID, playlistID, name, description)
	if err != nil {
		return nil, err
	}
	if pl == nil {
		return nil, ErrPlaylistNotFound
	}
	return pl, nil
}

func (s *PlaylistService) Delete(ctx context.Context, userID, playlistID int64) error {
	if err := s.repo.Delete(ctx, userID, playlistID); err != nil {
		return err
	}
	return nil
}

func (s *PlaylistService) AddTrack(ctx context.Context, userID, playlistID int64, trackID string, position int) ([]domain.PlaylistTrack, error) {
	if err := s.catalog.ValidateTrack(ctx, trackID); err != nil {
		return nil, fmt.Errorf("validate track: %w", err)
	}
	pl, err := s.repo.Get(ctx, userID, playlistID)
	if err != nil {
		return nil, err
	}
	if pl == nil {
		return nil, ErrPlaylistNotFound
	}
	if err := s.repo.AddTrack(ctx, playlistID, trackID, position); err != nil {
		return nil, err
	}
	return s.repo.ListTracks(ctx, playlistID)
}

func (s *PlaylistService) RemoveTrack(ctx context.Context, userID, playlistID int64, trackID string) ([]domain.PlaylistTrack, error) {
	pl, err := s.repo.Get(ctx, userID, playlistID)
	if err != nil {
		return nil, err
	}
	if pl == nil {
		return nil, ErrPlaylistNotFound
	}
	removed, err := s.repo.RemoveTrack(ctx, playlistID, trackID)
	if err != nil {
		return nil, err
	}
	if !removed {
		return nil, fmt.Errorf("track not found in playlist")
	}
	return s.repo.ListTracks(ctx, playlistID)
}
