package service

import (
	"context"
	"fmt"

	"github.com/freedom-music/media-service/internal/client"
	"github.com/freedom-music/media-service/internal/domain"
	"github.com/freedom-music/media-service/internal/storage"
)

type MediaService struct {
	storage *storage.MinioStorage
	catalog *client.CatalogClient
	objects map[string]domain.MediaObject
}

func NewMediaService(storage *storage.MinioStorage, catalog *client.CatalogClient) *MediaService {
	// Seed object mapping for Phase 2 demo
	objects := map[string]domain.MediaObject{
		"trk_1": {TrackID: "trk_1", ObjectKey: "tracks/trk_1.mp3", Bucket: "", ContentType: "audio/mpeg"},
		"trk_2": {TrackID: "trk_2", ObjectKey: "tracks/trk_2.mp3", Bucket: "", ContentType: "audio/mpeg"},
		"trk_3": {TrackID: "trk_3", ObjectKey: "tracks/trk_3.mp3", Bucket: "", ContentType: "audio/mpeg"},
	}
	return &MediaService{storage: storage, catalog: catalog, objects: objects}
}

func (s *MediaService) Metadata(ctx context.Context, trackID string) (*domain.MediaObject, error) {
	if err := s.catalog.ValidateTrack(ctx, trackID); err != nil {
		return nil, err
	}
	obj, ok := s.objects[trackID]
	if !ok {
		return nil, fmt.Errorf("media not found")
	}
	obj.Bucket = s.storage.Bucket()
	return &obj, nil
}

func (s *MediaService) SourceURL(ctx context.Context, trackID string) (string, error) {
	meta, err := s.Metadata(ctx, trackID)
	if err != nil {
		return "", err
	}
	exists, err := s.storage.ObjectExists(ctx, meta.ObjectKey)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", fmt.Errorf("media object missing")
	}
	// Always return direct URL for local dev (bucket is public). Avoid presign clock skew.
	return s.storage.PublicURL(meta.ObjectKey), nil
}
