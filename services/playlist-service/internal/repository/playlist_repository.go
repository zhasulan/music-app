package repository

import (
	"context"

	"github.com/freedom-music/playlist-service/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PlaylistRepository struct {
	db *pgxpool.Pool
}

func NewPlaylistRepository(db *pgxpool.Pool) *PlaylistRepository {
	return &PlaylistRepository{db: db}
}

func (r *PlaylistRepository) Create(ctx context.Context, userID int64, name, description string) (*domain.Playlist, error) {
	row := r.db.QueryRow(ctx, `
        INSERT INTO playlists (user_id, name, description)
        VALUES ($1, $2, $3)
        RETURNING id, user_id, name, description, created_at, updated_at
    `, userID, name, description)
	return scanPlaylist(row)
}

func (r *PlaylistRepository) List(ctx context.Context, userID int64) ([]domain.Playlist, error) {
	rows, err := r.db.Query(ctx, `
        SELECT id, user_id, name, description, created_at, updated_at
        FROM playlists
        WHERE user_id=$1
        ORDER BY created_at DESC
    `, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.Playlist
	for rows.Next() {
		var p domain.Playlist
		if err := rows.Scan(&p.ID, &p.UserID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, p)
	}
	return items, rows.Err()
}

func (r *PlaylistRepository) Get(ctx context.Context, userID, playlistID int64) (*domain.Playlist, error) {
	row := r.db.QueryRow(ctx, `
        SELECT id, user_id, name, description, created_at, updated_at
        FROM playlists
        WHERE id=$1 AND user_id=$2
    `, playlistID, userID)
	return scanPlaylist(row)
}

func (r *PlaylistRepository) Update(ctx context.Context, userID, playlistID int64, name, description string) (*domain.Playlist, error) {
	row := r.db.QueryRow(ctx, `
        UPDATE playlists
        SET name=COALESCE(NULLIF($3, ''), name),
            description=COALESCE(NULLIF($4, ''), description),
            updated_at=NOW()
        WHERE id=$1 AND user_id=$2
        RETURNING id, user_id, name, description, created_at, updated_at
    `, playlistID, userID, name, description)
	return scanPlaylist(row)
}

func (r *PlaylistRepository) Delete(ctx context.Context, userID, playlistID int64) error {
	_, err := r.db.Exec(ctx, `DELETE FROM playlists WHERE id=$1 AND user_id=$2`, playlistID, userID)
	return err
}

func (r *PlaylistRepository) ListTracks(ctx context.Context, playlistID int64) ([]domain.PlaylistTrack, error) {
	rows, err := r.db.Query(ctx, `
        SELECT id, playlist_id, track_id, position, added_at
        FROM playlist_tracks
        WHERE playlist_id=$1
        ORDER BY position ASC
    `, playlistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.PlaylistTrack
	for rows.Next() {
		var t domain.PlaylistTrack
		if err := rows.Scan(&t.ID, &t.PlaylistID, &t.TrackID, &t.Position, &t.AddedAt); err != nil {
			return nil, err
		}
		items = append(items, t)
	}
	return items, rows.Err()
}

func (r *PlaylistRepository) AddTrack(ctx context.Context, playlistID int64, trackID string, position int) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var maxPos int
	if err := tx.QueryRow(ctx, `SELECT COALESCE(MAX(position),0) FROM playlist_tracks WHERE playlist_id=$1`, playlistID).Scan(&maxPos); err != nil {
		return err
	}
	if position <= 0 || position > maxPos+1 {
		position = maxPos + 1
	}

	if _, err := tx.Exec(ctx, `
        UPDATE playlist_tracks
        SET position = position + 1
        WHERE playlist_id=$1 AND position >= $2
    `, playlistID, position); err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, `
        INSERT INTO playlist_tracks (playlist_id, track_id, position)
        VALUES ($1, $2, $3)
    `, playlistID, trackID, position); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *PlaylistRepository) RemoveTrack(ctx context.Context, playlistID int64, trackID string) (bool, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return false, err
	}
	defer tx.Rollback(ctx)

	var rowID int64
	var pos int
	if err := tx.QueryRow(ctx, `
        SELECT id, position
        FROM playlist_tracks
        WHERE playlist_id=$1 AND track_id=$2
        ORDER BY position
        LIMIT 1
    `, playlistID, trackID).Scan(&rowID, &pos); err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	if _, err := tx.Exec(ctx, `DELETE FROM playlist_tracks WHERE id=$1`, rowID); err != nil {
		return false, err
	}
	if _, err := tx.Exec(ctx, `
        UPDATE playlist_tracks
        SET position = position - 1
        WHERE playlist_id=$1 AND position > $2
    `, playlistID, pos); err != nil {
		return false, err
	}

	return true, tx.Commit(ctx)
}

func scanPlaylist(row pgx.Row) (*domain.Playlist, error) {
	var p domain.Playlist
	if err := row.Scan(&p.ID, &p.UserID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}
