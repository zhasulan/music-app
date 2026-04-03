package service

import (
	"strings"

	"github.com/freedom-music/search-service/internal/client"
)

type Result struct {
	Type   string         `json:"type"`
	Track  *client.Track  `json:"track,omitempty"`
	Album  *client.Album  `json:"album,omitempty"`
	Artist *client.Artist `json:"artist,omitempty"`
	Score  int            `json:"score"`
}

type Suggestion struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

type Service struct {
	tracks  []client.Track
	albums  []client.Album
	artists []client.Artist
}

func New(tracks []client.Track, albums []client.Album, artists []client.Artist) *Service {
	return &Service{tracks: tracks, albums: albums, artists: artists}
}

func (s *Service) Search(query string) []Result {
	q := strings.TrimSpace(strings.ToLower(query))
	if q == "" {
		return []Result{}
	}
	var results []Result
	for _, t := range s.tracks {
		if contains(t.Title, q) {
			results = append(results, Result{Type: "track", Track: &t, Score: scoreContains(t.Title, q)})
		}
	}
	for _, a := range s.albums {
		if contains(a.Title, q) {
			results = append(results, Result{Type: "album", Album: &a, Score: scoreContains(a.Title, q)})
		}
	}
	for _, art := range s.artists {
		if contains(art.Name, q) {
			results = append(results, Result{Type: "artist", Artist: &art, Score: scoreContains(art.Name, q)})
		}
	}
	// simple score sort descending
	for i := 0; i < len(results); i++ {
		for j := i + 1; j < len(results); j++ {
			if results[j].Score > results[i].Score {
				results[i], results[j] = results[j], results[i]
			}
		}
	}
	return results
}

func (s *Service) Suggest(query string) []Suggestion {
	q := strings.TrimSpace(strings.ToLower(query))
	if q == "" {
		return []Suggestion{}
	}
	seen := map[string]bool{}
	var out []Suggestion
	for _, t := range s.tracks {
		if contains(t.Title, q) && !seen[t.Title] {
			out = append(out, Suggestion{Text: t.Title, Type: "track"})
			seen[t.Title] = true
		}
	}
	for _, a := range s.albums {
		if contains(a.Title, q) && !seen[a.Title] {
			out = append(out, Suggestion{Text: a.Title, Type: "album"})
			seen[a.Title] = true
		}
	}
	for _, art := range s.artists {
		if contains(art.Name, q) && !seen[art.Name] {
			out = append(out, Suggestion{Text: art.Name, Type: "artist"})
			seen[art.Name] = true
		}
	}
	if len(out) > 5 {
		out = out[:5]
	}
	return out
}

func contains(text, q string) bool {
	return strings.Contains(strings.ToLower(text), q)
}

func scoreContains(text, q string) int {
	textLower := strings.ToLower(text)
	if textLower == q {
		return 100
	}
	if strings.HasPrefix(textLower, q) {
		return 80
	}
	if strings.Contains(textLower, q) {
		return 60
	}
	return 50
}
