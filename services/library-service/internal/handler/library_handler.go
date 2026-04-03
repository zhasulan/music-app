package handler

import (
	"net/http"
	"strconv"

	"github.com/freedom-music/library-service/internal/service"
	sharedErrors "github.com/freedom-music/shared/errors"
	"github.com/freedom-music/shared/httpx"
	"github.com/freedom-music/shared/middleware"
	"github.com/gin-gonic/gin"
)

type LibraryHandler struct {
	svc *service.LibraryService
}

type likeRequest struct{}

func NewLibraryHandler(svc *service.LibraryService) *LibraryHandler {
	return &LibraryHandler{svc: svc}
}

func (h *LibraryHandler) RegisterRoutes(r *gin.Engine, jwtSecret string) {
	r.GET("/health", httpx.Health)

	api := r.Group("/api/v1")
	api.Use(middleware.AuthRequired(jwtSecret))
	{
		api.PUT("/library/liked-tracks/:trackId", h.like)
		api.DELETE("/library/liked-tracks/:trackId", h.unlike)
		api.GET("/library/liked-tracks", h.list)
	}
}

func (h *LibraryHandler) like(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		httpx.JSONError(c, http.StatusUnauthorized, sharedErrors.CodeUnauthorized, "missing user context")
		return
	}
	trackID := c.Param("trackId")
	lt, err := h.svc.Like(c.Request.Context(), userID, trackID)
	if err != nil {
		if err.Error() == "track not found" {
			httpx.JSONError(c, http.StatusNotFound, sharedErrors.CodeNotFound, "track not found")
			return
		}
		httpx.JSONError(c, http.StatusBadRequest, sharedErrors.CodeInvalidRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"liked": lt})
}

func (h *LibraryHandler) unlike(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		httpx.JSONError(c, http.StatusUnauthorized, sharedErrors.CodeUnauthorized, "missing user context")
		return
	}
	trackID := c.Param("trackId")
	if err := h.svc.Unlike(c.Request.Context(), userID, trackID); err != nil {
		httpx.JSONError(c, http.StatusInternalServerError, sharedErrors.CodeInternal, "could not unlike")
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "removed"})
}

func (h *LibraryHandler) list(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		httpx.JSONError(c, http.StatusUnauthorized, sharedErrors.CodeUnauthorized, "missing user context")
		return
	}
	items, err := h.svc.List(c.Request.Context(), userID)
	if err != nil {
		httpx.JSONError(c, http.StatusInternalServerError, sharedErrors.CodeInternal, "could not list")
		return
	}
	c.JSON(http.StatusOK, gin.H{"liked_tracks": items})
}

func userIDFromContext(c *gin.Context) (int64, bool) {
	raw, ok := c.Get(middleware.UserIDContextKey)
	if !ok {
		return 0, false
	}
	id, err := strconv.ParseInt(raw.(string), 10, 64)
	if err != nil {
		return 0, false
	}
	return id, true
}
