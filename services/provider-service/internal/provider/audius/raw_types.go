package audius

import "encoding/json"

type Artwork struct {
	Small  string `json:"150x150"`
	Medium string `json:"480x480"`
	Large  string `json:"1000x1000"`
}

type AudiusTrack struct {
	ID           string     `json:"id"`
	Title        string     `json:"title"`
	Duration     int        `json:"duration"`
	Genre        string     `json:"genre"`
	Artwork      Artwork    `json:"artwork"`
	ArtworkURL   string     `json:"artwork_url"`
	Downloadable bool       `json:"downloadable"`
	DownloadURL  string     `json:"download_url"`
	Permalink    string     `json:"permalink"`
	StreamURL    string     `json:"stream_url"`
	User         AudiusUser `json:"user"`
}

type AudiusUser struct {
	ID                 string          `json:"id"`
	Handle             string          `json:"handle"`
	Name               string          `json:"name"`
	ProfilePicRaw      json.RawMessage `json:"profile_picture"`
	ProfilePicSizesRaw json.RawMessage `json:"profile_picture_sizes"`
	Followers          int             `json:"follower_count"`
}

type AudiusPlaylist struct {
	ID           string  `json:"id"`
	PlaylistName string  `json:"playlist_name"`
	Artwork      Artwork `json:"artwork"`
	ArtworkURL   string  `json:"artwork_url"`
	Permalink    string  `json:"permalink"`
}

func normalizeTrack(t AudiusTrack) Track {
	art := firstNonEmpty(t.Artwork.Large, t.ArtworkURL, t.Artwork.Medium, t.Artwork.Small)
	stream := t.StreamURL
	if stream == "" {
		// Fallback constructed stream endpoint; actual baseURL appended later where needed.
		// The caller may replace this placeholder with an upstream stream endpoint via Provider.StreamURL.
		stream = ""
	}
	return Track{
		ID:            "audius:" + t.ID,
		Provider:      "audius",
		ProviderTrack: t.ID,
		Title:         t.Title,
		ArtistName:    t.User.Name,
		ArtistHandle:  t.User.Handle,
		ArtworkURL:    art,
		Duration:      t.Duration,
		Genre:         t.Genre,
		StreamURL:     stream,
		Downloadable:  t.Downloadable,
		DownloadURL:   t.DownloadURL,
		Permalink:     t.Permalink,
	}
}

func normalizeArtist(u AudiusUser) Artist {
	pic := pickPicture(u.ProfilePicRaw, u.ProfilePicSizesRaw)
	return Artist{
		ID:           u.ID,
		Handle:       u.Handle,
		Name:         u.Name,
		ProfileImage: pic,
		Followers:    u.Followers,
	}
}

func normalizePlaylist(p AudiusPlaylist) Playlist {
	return Playlist{
		ID:        p.ID,
		Title:     p.PlaylistName,
		Artwork:   firstNonEmpty(p.Artwork.Large, p.ArtworkURL, p.Artwork.Medium, p.Artwork.Small),
		Permalink: p.Permalink,
	}
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

func pickPicture(raw json.RawMessage, sizesRaw json.RawMessage) string {
	if len(raw) > 0 {
		var s string
		if err := json.Unmarshal(raw, &s); err == nil && s != "" {
			return s
		}
		var art Artwork
		if err := json.Unmarshal(raw, &art); err == nil {
			if v := firstNonEmpty(art.Large, art.Medium, art.Small); v != "" {
				return v
			}
		}
		var m map[string]string
		if err := json.Unmarshal(raw, &m); err == nil {
			if v := firstNonEmpty(m["1000x1000"], m["480x480"], m["150x150"]); v != "" {
				return v
			}
		}
	}
	if len(sizesRaw) > 0 {
		var s string
		if err := json.Unmarshal(sizesRaw, &s); err == nil {
			// if it's a CID string, not directly usable, fall through
		} else {
			var art Artwork
			if err := json.Unmarshal(sizesRaw, &art); err == nil {
				if v := firstNonEmpty(art.Large, art.Medium, art.Small); v != "" {
					return v
				}
			}
			var m map[string]string
			if err := json.Unmarshal(sizesRaw, &m); err == nil {
				if v := firstNonEmpty(m["1000x1000"], m["480x480"], m["150x150"]); v != "" {
					return v
				}
			}
		}
	}
	return ""
}
