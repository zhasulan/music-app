package repository

import (
	"context"

	"github.com/freedom-music/library-service/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LibraryRepository struct {
	db *pgxpool.Pool
}

func NewLibraryRepository(db *pgxpool.Pool) *LibraryRepository {
	return &LibraryRepository{db: db}
}

func (r *LibraryRepository) Like(ctx context.Context, userID int64, trackID string) (*domain.LikedTrack, error) {
	row := r.db.QueryRow(ctx, `
        INSERT INTO liked_tracks (user_id, track_id)
        VALUES ($1, $2)
        ON CONFLICT (user_id, track_id) DO UPDATE SET track_id=EXCLUDED.track_id
        RETURNING user_id, track_id, created_at
    `, userID, trackID)
	return scanLiked(row)
}

func (r *LibraryRepository) Unlike(ctx context.Context, userID int64, trackID string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM liked_tracks WHERE user_id=$1 AND track_id=$2`, userID, trackID)
	return err
}

func (r *LibraryRepository) List(ctx context.Context, userID int64) ([]domain.LikedTrack, error) {
	rows, err := r.db.Query(ctx, `
        SELECT user_id, track_id, created_at
        FROM liked_tracks
        WHERE user_id=$1
        ORDER BY created_at DESC
    `, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.LikedTrack
	for rows.Next() {
		var lt domain.LikedTrack
		if err := rows.Scan(&lt.UserID, &lt.TrackID, &lt.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, lt)
	}
	return items, rows.Err()
}

func scanLiked(row pgx.Row) (*domain.LikedTrack, error) {
	var lt domain.LikedTrack
	if err := row.Scan(&lt.UserID, &lt.TrackID, &lt.CreatedAt); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &lt, nil
}
