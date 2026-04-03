package service

import "github.com/freedom-music/catalog-service/internal/repository"

// CatalogService is currently a thin wrapper over the in-memory repository.
type CatalogService struct {
	repo *repository.CatalogRepository
}

func NewCatalogService(repo *repository.CatalogRepository) *CatalogService {
	return &CatalogService{repo: repo}
}

func (s *CatalogService) ListTracks() interface{}      { return s.repo.ListTracks() }
func (s *CatalogService) ListAlbums() interface{}      { return s.repo.ListAlbums() }
func (s *CatalogService) ListArtists() interface{}     { return s.repo.ListArtists() }
func (s *CatalogService) Track(id string) interface{}  { return s.repo.GetTrackByID(id) }
func (s *CatalogService) Album(id string) interface{}  { return s.repo.GetAlbumByID(id) }
func (s *CatalogService) Artist(id string) interface{} { return s.repo.GetArtistByID(id) }
