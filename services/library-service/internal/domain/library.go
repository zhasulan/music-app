package domain

import "time"

type LikedTrack struct {
	UserID    int64     `json:"user_id"`
	TrackID   string    `json:"track_id"`
	CreatedAt time.Time `json:"created_at"`
}
