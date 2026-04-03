package repository

import "github.com/freedom-music/catalog-service/internal/domain"

type CatalogRepository struct {
	tracks  []domain.Track
	albums  []domain.Album
	artists []domain.Artist
}

func NewCatalogRepository() *CatalogRepository {
	artists := []domain.Artist{
		{ID: "art_1", Name: "The Starters"},
		{ID: "art_2", Name: "Lo-Fi Coder"},
	}
	albums := []domain.Album{
		{ID: "alb_1", Title: "Launch Day", ArtistID: "art_1", Year: 2024},
		{ID: "alb_2", Title: "Night Commits", ArtistID: "art_2", Year: 2023},
	}
	tracks := []domain.Track{
		{ID: "trk_1", Title: "Boot Sequence", ArtistID: "art_1", AlbumID: "alb_1", Duration: 185},
		{ID: "trk_2", Title: "Merge", ArtistID: "art_1", AlbumID: "alb_1", Duration: 201},
		{ID: "trk_3", Title: "Code & Chill", ArtistID: "art_2", AlbumID: "alb_2", Duration: 230},
	}

	return &CatalogRepository{tracks: tracks, albums: albums, artists: artists}
}

func (r *CatalogRepository) ListTracks() []domain.Track   { return r.tracks }
func (r *CatalogRepository) ListAlbums() []domain.Album   { return r.albums }
func (r *CatalogRepository) ListArtists() []domain.Artist { return r.artists }

func (r *CatalogRepository) GetTrackByID(id string) *domain.Track {
	for _, t := range r.tracks {
		if t.ID == id {
			return &t
		}
	}
	return nil
}

func (r *CatalogRepository) GetAlbumByID(id string) *domain.Album {
	for _, a := range r.albums {
		if a.ID == id {
			return &a
		}
	}
	return nil
}

func (r *CatalogRepository) GetArtistByID(id string) *domain.Artist {
	for _, a := range r.artists {
		if a.ID == id {
			return &a
		}
	}
	return nil
}
