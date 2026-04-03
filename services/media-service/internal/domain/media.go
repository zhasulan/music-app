package domain

type MediaObject struct {
	TrackID     string `json:"track_id"`
	ObjectKey   string `json:"object_key"`
	Bucket      string `json:"bucket"`
	ContentType string `json:"content_type"`
}
