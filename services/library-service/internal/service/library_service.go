package service

import (
	"context"

	"github.com/freedom-music/library-service/internal/client"
	"github.com/freedom-music/library-service/internal/domain"
	"github.com/freedom-music/library-service/internal/repository"
)

type LibraryService struct {
	repo    *repository.LibraryRepository
	catalog *client.CatalogClient
}

func NewLibraryService(repo *repository.LibraryRepository, catalog *client.CatalogClient) *LibraryService {
	return &LibraryService{repo: repo, catalog: catalog}
}

func (s *LibraryService) Like(ctx context.Context, userID int64, trackID string) (*domain.LikedTrack, error) {
	if !isExternal(trackID) {
		if err := s.catalog.ValidateTrack(ctx, trackID); err != nil {
			return nil, err
		}
	}
	return s.repo.Like(ctx, userID, trackID)
}

func (s *LibraryService) Unlike(ctx context.Context, userID int64, trackID string) error {
	return s.repo.Unlike(ctx, userID, trackID)
}

func (s *LibraryService) List(ctx context.Context, userID int64) ([]domain.LikedTrack, error) {
	return s.repo.List(ctx, userID)
}

func isExternal(trackID string) bool {
	return len(trackID) > 7 && trackID[:7] == "audius:"
}
