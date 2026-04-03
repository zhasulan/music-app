package domain

type Track struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	ArtistID string `json:"artist_id"`
	AlbumID  string `json:"album_id"`
	Duration int    `json:"duration_sec"`
}

type Album struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	ArtistID string `json:"artist_id"`
	Year     int    `json:"year"`
}

type Artist struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
