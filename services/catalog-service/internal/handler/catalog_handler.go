package handler

import (
	"net/http"

	"github.com/freedom-music/catalog-service/internal/service"
	sharedErrors "github.com/freedom-music/shared/errors"
	"github.com/freedom-music/shared/httpx"
	"github.com/gin-gonic/gin"
)

type CatalogHandler struct {
	svc *service.CatalogService
}

func NewCatalogHandler(svc *service.CatalogService) *CatalogHandler {
	return &CatalogHandler{svc: svc}
}

func (h *CatalogHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/health", httpx.Health)

	api := r.Group("/api/v1/catalog")
	{
		api.GET("/tracks", h.listTracks)
		api.GET("/tracks/:id", h.track)
		api.GET("/albums", h.listAlbums)
		api.GET("/albums/:id", h.album)
		api.GET("/artists", h.listArtists)
		api.GET("/artists/:id", h.artist)
	}
}

func (h *CatalogHandler) listTracks(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"tracks": h.svc.ListTracks()})
}

func (h *CatalogHandler) track(c *gin.Context) {
	tr := h.svc.Track(c.Param("id"))
	if tr == nil {
		httpx.JSONError(c, http.StatusNotFound, sharedErrors.CodeNotFound, "track not found")
		return
	}
	c.JSON(http.StatusOK, gin.H{"track": tr})
}

func (h *CatalogHandler) listAlbums(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"albums": h.svc.ListAlbums()})
}

func (h *CatalogHandler) album(c *gin.Context) {
	alb := h.svc.Album(c.Param("id"))
	if alb == nil {
		httpx.JSONError(c, http.StatusNotFound, sharedErrors.CodeNotFound, "album not found")
		return
	}
	c.JSON(http.StatusOK, gin.H{"album": alb})
}

func (h *CatalogHandler) listArtists(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"artists": h.svc.ListArtists()})
}

func (h *CatalogHandler) artist(c *gin.Context) {
	art := h.svc.Artist(c.Param("id"))
	if art == nil {
		httpx.JSONError(c, http.StatusNotFound, sharedErrors.CodeNotFound, "artist not found")
		return
	}
	c.JSON(http.StatusOK, gin.H{"artist": art})
}
