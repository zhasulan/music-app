package domain

import "time"

type Playlist struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PlaylistTrack struct {
	ID         int64     `json:"id"`
	PlaylistID int64     `json:"playlist_id"`
	TrackID    string    `json:"track_id"`
	Position   int       `json:"position"`
	AddedAt    time.Time `json:"added_at"`
}
