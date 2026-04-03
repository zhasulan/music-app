package audius

type Track struct {
	ID            string `json:"id"`
	Provider      string `json:"provider"`
	ProviderTrack string `json:"providerTrackId"`
	Title         string `json:"title"`
	ArtistName    string `json:"artistName"`
	ArtistHandle  string `json:"artistHandle"`
	ArtworkURL    string `json:"artworkUrl"`
	Duration      int    `json:"duration"`
	Genre         string `json:"genre"`
	StreamURL     string `json:"streamUrl"`
	Downloadable  bool   `json:"downloadable"`
	DownloadURL   string `json:"downloadUrl"`
	Permalink     string `json:"permalink"`
}

type Artist struct {
	ID           string `json:"id"`
	Handle       string `json:"handle"`
	Name         string `json:"name"`
	ProfileImage string `json:"profileImage"`
	Followers    int    `json:"followers"`
}

type Playlist struct {
	ID        string  `json:"id"`
	Title     string  `json:"title"`
	Artwork   string  `json:"artworkUrl"`
	Tracks    []Track `json:"tracks,omitempty"`
	Permalink string  `json:"permalink"`
}

type Paged[T any] struct {
	Items   []T  `json:"items"`
	Limit   int  `json:"limit"`
	Offset  int  `json:"offset"`
	HasMore bool `json:"hasMore"`
	Total   *int `json:"total,omitempty"`
}
