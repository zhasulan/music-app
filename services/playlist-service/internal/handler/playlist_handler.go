package handler

import (
	"net/http"
	"strconv"

	"github.com/freedom-music/playlist-service/internal/service"
	sharedErrors "github.com/freedom-music/shared/errors"
	"github.com/freedom-music/shared/httpx"
	"github.com/freedom-music/shared/middleware"
	"github.com/gin-gonic/gin"
)

type PlaylistHandler struct {
	svc *service.PlaylistService
}

type createPlaylistRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"omitempty"`
}

type updatePlaylistRequest struct {
	Name        string `json:"name" binding:"omitempty"`
	Description string `json:"description" binding:"omitempty"`
}

type addTrackRequest struct {
	TrackID  string `json:"track_id" binding:"required"`
	Position *int   `json:"position"`
}

func NewPlaylistHandler(svc *service.PlaylistService) *PlaylistHandler {
	return &PlaylistHandler{svc: svc}
}

func (h *PlaylistHandler) RegisterRoutes(r *gin.Engine, jwtSecret string) {
	r.GET("/health", httpx.Health)

	api := r.Group("/api/v1")
	api.Use(middleware.AuthRequired(jwtSecret))
	{
		api.POST("/playlists", h.create)
		api.GET("/playlists", h.list)
		api.GET("/playlists/:id", h.get)
		api.PATCH("/playlists/:id", h.update)
		api.DELETE("/playlists/:id", h.delete)
		api.POST("/playlists/:id/tracks", h.addTrack)
		api.DELETE("/playlists/:id/tracks/:trackId", h.removeTrack)
	}
}

func (h *PlaylistHandler) create(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		httpx.JSONError(c, http.StatusUnauthorized, sharedErrors.CodeUnauthorized, "missing user context")
		return
	}

	var req createPlaylistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.JSONError(c, http.StatusBadRequest, sharedErrors.CodeInvalidRequest, "invalid payload")
		return
	}

	pl, err := h.svc.Create(c.Request.Context(), userID, req.Name, req.Description)
	if err != nil {
		httpx.JSONError(c, http.StatusInternalServerError, sharedErrors.CodeInternal, "could not create playlist")
		return
	}
	c.JSON(http.StatusOK, gin.H{"playlist": pl})
}

func (h *PlaylistHandler) list(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		httpx.JSONError(c, http.StatusUnauthorized, sharedErrors.CodeUnauthorized, "missing user context")
		return
	}
	pls, err := h.svc.List(c.Request.Context(), userID)
	if err != nil {
		httpx.JSONError(c, http.StatusInternalServerError, sharedErrors.CodeInternal, "could not list playlists")
		return
	}
	c.JSON(http.StatusOK, gin.H{"playlists": pls})
}

func (h *PlaylistHandler) get(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		httpx.JSONError(c, http.StatusUnauthorized, sharedErrors.CodeUnauthorized, "missing user context")
		return
	}
	playlistID, err := parseID(c.Param("id"))
	if err != nil {
		httpx.JSONError(c, http.StatusBadRequest, sharedErrors.CodeInvalidRequest, "invalid playlist id")
		return
	}
	pl, tracks, err := h.svc.Get(c.Request.Context(), userID, playlistID)
	if err != nil {
		if err == service.ErrPlaylistNotFound {
			httpx.JSONError(c, http.StatusNotFound, sharedErrors.CodeNotFound, "playlist not found")
			return
		}
		httpx.JSONError(c, http.StatusInternalServerError, sharedErrors.CodeInternal, "could not load playlist")
		return
	}
	c.JSON(http.StatusOK, gin.H{"playlist": pl, "tracks": tracks})
}

func (h *PlaylistHandler) update(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		httpx.JSONError(c, http.StatusUnauthorized, sharedErrors.CodeUnauthorized, "missing user context")
		return
	}
	playlistID, err := parseID(c.Param("id"))
	if err != nil {
		httpx.JSONError(c, http.StatusBadRequest, sharedErrors.CodeInvalidRequest, "invalid playlist id")
		return
	}
	var req updatePlaylistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.JSONError(c, http.StatusBadRequest, sharedErrors.CodeInvalidRequest, "invalid payload")
		return
	}
	pl, err := h.svc.Update(c.Request.Context(), userID, playlistID, req.Name, req.Description)
	if err != nil {
		if err == service.ErrPlaylistNotFound {
			httpx.JSONError(c, http.StatusNotFound, sharedErrors.CodeNotFound, "playlist not found")
			return
		}
		httpx.JSONError(c, http.StatusInternalServerError, sharedErrors.CodeInternal, "could not update playlist")
		return
	}
	c.JSON(http.StatusOK, gin.H{"playlist": pl})
}

func (h *PlaylistHandler) delete(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		httpx.JSONError(c, http.StatusUnauthorized, sharedErrors.CodeUnauthorized, "missing user context")
		return
	}
	playlistID, err := parseID(c.Param("id"))
	if err != nil {
		httpx.JSONError(c, http.StatusBadRequest, sharedErrors.CodeInvalidRequest, "invalid playlist id")
		return
	}
	if err := h.svc.Delete(c.Request.Context(), userID, playlistID); err != nil {
		httpx.JSONError(c, http.StatusInternalServerError, sharedErrors.CodeInternal, "could not delete playlist")
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func (h *PlaylistHandler) addTrack(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		httpx.JSONError(c, http.StatusUnauthorized, sharedErrors.CodeUnauthorized, "missing user context")
		return
	}
	playlistID, err := parseID(c.Param("id"))
	if err != nil {
		httpx.JSONError(c, http.StatusBadRequest, sharedErrors.CodeInvalidRequest, "invalid playlist id")
		return
	}
	var req addTrackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.JSONError(c, http.StatusBadRequest, sharedErrors.CodeInvalidRequest, "invalid payload")
		return
	}
	position := 0
	if req.Position != nil {
		position = *req.Position
	}
	tracks, err := h.svc.AddTrack(c.Request.Context(), userID, playlistID, req.TrackID, position)
	if err != nil {
		if err == service.ErrPlaylistNotFound {
			httpx.JSONError(c, http.StatusNotFound, sharedErrors.CodeNotFound, "playlist not found")
			return
		}
		if err.Error() == "track not found" {
			httpx.JSONError(c, http.StatusNotFound, sharedErrors.CodeNotFound, "track not found")
			return
		}
		httpx.JSONError(c, http.StatusInternalServerError, sharedErrors.CodeInternal, "could not add track")
		return
	}
	c.JSON(http.StatusOK, gin.H{"tracks": tracks})
}

func (h *PlaylistHandler) removeTrack(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		httpx.JSONError(c, http.StatusUnauthorized, sharedErrors.CodeUnauthorized, "missing user context")
		return
	}
	playlistID, err := parseID(c.Param("id"))
	if err != nil {
		httpx.JSONError(c, http.StatusBadRequest, sharedErrors.CodeInvalidRequest, "invalid playlist id")
		return
	}
	trackID := c.Param("trackId")
	tracks, err := h.svc.RemoveTrack(c.Request.Context(), userID, playlistID, trackID)
	if err != nil {
		if err == service.ErrPlaylistNotFound {
			httpx.JSONError(c, http.StatusNotFound, sharedErrors.CodeNotFound, "playlist not found")
			return
		}
		httpx.JSONError(c, http.StatusBadRequest, sharedErrors.CodeInvalidRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"tracks": tracks})
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

func parseID(raw string) (int64, error) {
	return strconv.ParseInt(raw, 10, 64)
}
