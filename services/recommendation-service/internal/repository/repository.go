package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Trending(ctx context.Context, limit int) ([]string, error) {
	rows, err := r.db.Query(ctx, `
		SELECT track_id
		FROM events
		WHERE track_id IS NOT NULL AND event_type IN ('track_played','track_liked')
		GROUP BY track_id
		ORDER BY COUNT(*) DESC, MAX(created_at) DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (r *Repository) RecentlyPlayed(ctx context.Context, userID int64, limit int) ([]string, error) {
	rows, err := r.db.Query(ctx, `
		SELECT track_id
		FROM events
		WHERE user_id = $1 AND event_type = 'track_played' AND track_id IS NOT NULL
		ORDER BY created_at DESC
		LIMIT $2
	`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (r *Repository) LikedTracks(ctx context.Context, userID int64, limit int) ([]string, error) {
	rows, err := r.db.Query(ctx, `
		SELECT track_id
		FROM events
		WHERE user_id = $1 AND event_type = 'track_liked' AND track_id IS NOT NULL
		ORDER BY created_at DESC
		LIMIT $2
	`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}
