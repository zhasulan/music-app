package service

import (
	"context"
	"errors"
	"time"

	"github.com/freedom-music/playback-service/internal/client"
	"github.com/freedom-music/playback-service/internal/domain"
	"github.com/freedom-music/playback-service/internal/store"
)

var (
	ErrSessionNotFound = errors.New("playback session not found")
)

type PlaybackService struct {
	store   *store.RedisStore
	catalog *client.CatalogClient
}

func NewPlaybackService(store *store.RedisStore, catalog *client.CatalogClient) *PlaybackService {
	return &PlaybackService{store: store, catalog: catalog}
}

func (s *PlaybackService) Start(ctx context.Context, userID int64, trackID string, positionMs int) (*domain.PlaybackSession, error) {
	// Audius (and other external providers) are not present in local catalog.
	// Allow them to pass through so playback works for streamed tracks.
	if !isExternal(trackID) {
		if err := s.catalog.ValidateTrack(ctx, trackID); err != nil {
			return nil, err
		}
	}
	sess := &domain.PlaybackSession{
		UserID:     userID,
		TrackID:    trackID,
		PositionMs: positionMs,
		State:      domain.PlaybackStatePlaying,
		UpdatedAt:  time.Now().UTC(),
	}
	if err := s.store.Save(ctx, sess); err != nil {
		return nil, err
	}
	return sess, nil
}

func (s *PlaybackService) Pause(ctx context.Context, userID int64, positionMs int) (*domain.PlaybackSession, error) {
	sess, err := s.store.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	if sess == nil {
		return nil, ErrSessionNotFound
	}
	sess.PositionMs = positionMs
	sess.State = domain.PlaybackStatePaused
	sess.UpdatedAt = time.Now().UTC()
	if err := s.store.Save(ctx, sess); err != nil {
		return nil, err
	}
	return sess, nil
}

func (s *PlaybackService) Resume(ctx context.Context, userID int64) (*domain.PlaybackSession, error) {
	sess, err := s.store.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	if sess == nil {
		return nil, ErrSessionNotFound
	}
	sess.State = domain.PlaybackStatePlaying
	sess.UpdatedAt = time.Now().UTC()
	if err := s.store.Save(ctx, sess); err != nil {
		return nil, err
	}
	return sess, nil
}

func (s *PlaybackService) Seek(ctx context.Context, userID int64, positionMs int) (*domain.PlaybackSession, error) {
	sess, err := s.store.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	if sess == nil {
		return nil, ErrSessionNotFound
	}
	sess.PositionMs = positionMs
	sess.UpdatedAt = time.Now().UTC()
	if err := s.store.Save(ctx, sess); err != nil {
		return nil, err
	}
	return sess, nil
}

func (s *PlaybackService) Current(ctx context.Context, userID int64) (*domain.PlaybackSession, error) {
	sess, err := s.store.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	if sess == nil {
		return nil, ErrSessionNotFound
	}
	return sess, nil
}

// isExternal returns true for provider-scoped IDs (e.g., "audius:*").
func isExternal(trackID string) bool {
	return len(trackID) > 7 && trackID[:7] == "audius:"
}
