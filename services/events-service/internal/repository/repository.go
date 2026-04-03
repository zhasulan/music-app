package repository

import (
	"context"

	"github.com/freedom-music/events-service/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Insert(ctx context.Context, e *domain.Event) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO events (user_id, event_type, track_id, playlist_id, query)
		VALUES ($1, $2, $3, $4, $5)
	`, e.UserID, e.EventType, e.TrackID, e.PlaylistID, e.Query)
	return err
}
