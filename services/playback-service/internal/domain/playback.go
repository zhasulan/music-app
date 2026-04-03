package domain

import "time"

type PlaybackState string

const (
	PlaybackStatePlaying PlaybackState = "playing"
	PlaybackStatePaused  PlaybackState = "paused"
)

type PlaybackSession struct {
	UserID     int64         `json:"user_id"`
	TrackID    string        `json:"track_id"`
	PositionMs int           `json:"position_ms"`
	State      PlaybackState `json:"state"`
	UpdatedAt  time.Time     `json:"updated_at"`
}
