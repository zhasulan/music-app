package domain

import "time"

type Event struct {
	ID         int64     `json:"id"`
	UserID     *int64    `json:"user_id,omitempty"`
	EventType  string    `json:"event_type"`
	TrackID    *string   `json:"track_id,omitempty"`
	PlaylistID *int64    `json:"playlist_id,omitempty"`
	Query      *string   `json:"query,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}
