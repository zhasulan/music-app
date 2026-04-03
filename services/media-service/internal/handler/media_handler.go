package handler

import (
	"net/http"

	"github.com/freedom-music/media-service/internal/service"
	sharedErrors "github.com/freedom-music/shared/errors"
	"github.com/freedom-music/shared/httpx"
	"github.com/gin-gonic/gin"
)

type MediaHandler struct {
	svc *service.MediaService
}

func NewMediaHandler(svc *service.MediaService) *MediaHandler {
	return &MediaHandler{svc: svc}
}

func (h *MediaHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/health", httpx.Health)

	api := r.Group("/api/v1/media")
	{
		api.GET("/tracks/:trackId", h.metadata)
		api.GET("/tracks/:trackId/source", h.source)
	}
}

func (h *MediaHandler) metadata(c *gin.Context) {
	trackID := c.Param("trackId")
	meta, err := h.svc.Metadata(c.Request.Context(), trackID)
	if err != nil {
		httpx.JSONError(c, http.StatusNotFound, sharedErrors.CodeNotFound, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"media": meta})
}

func (h *MediaHandler) source(c *gin.Context) {
	trackID := c.Param("trackId")
	url, err := h.svc.SourceURL(c.Request.Context(), trackID)
	if err != nil {
		httpx.JSONError(c, http.StatusNotFound, sharedErrors.CodeNotFound, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"url": url})
}
