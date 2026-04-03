package handler

import (
	"net/http"
	"strconv"

	"github.com/freedom-music/playback-service/internal/service"
	sharedErrors "github.com/freedom-music/shared/errors"
	"github.com/freedom-music/shared/httpx"
	"github.com/freedom-music/shared/middleware"
	"github.com/gin-gonic/gin"
)

type PlaybackHandler struct {
	svc *service.PlaybackService
}

type startRequest struct {
	TrackID    string `json:"track_id" binding:"required"`
	PositionMs int    `json:"position_ms" binding:"omitempty"`
}

type positionRequest struct {
	PositionMs int `json:"position_ms" binding:"required"`
}

func NewPlaybackHandler(svc *service.PlaybackService) *PlaybackHandler {
	return &PlaybackHandler{svc: svc}
}

func (h *PlaybackHandler) RegisterRoutes(r *gin.Engine, jwtSecret string) {
	r.GET("/health", httpx.Health)

	api := r.Group("/api/v1")
	api.Use(middleware.AuthRequired(jwtSecret))
	{
		api.POST("/playback/start", h.start)
		api.POST("/playback/pause", h.pause)
		api.POST("/playback/resume", h.resume)
		api.POST("/playback/seek", h.seek)
		api.GET("/playback/current", h.current)
	}
}

func (h *PlaybackHandler) start(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		httpx.JSONError(c, http.StatusUnauthorized, sharedErrors.CodeUnauthorized, "missing user context")
		return
	}
	var req startRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.JSONError(c, http.StatusBadRequest, sharedErrors.CodeInvalidRequest, "invalid payload")
		return
	}
	sess, err := h.svc.Start(c.Request.Context(), userID, req.TrackID, req.PositionMs)
	if err != nil {
		if err.Error() == "track not found" {
			httpx.JSONError(c, http.StatusNotFound, sharedErrors.CodeNotFound, "track not found")
			return
		}
		httpx.JSONError(c, http.StatusInternalServerError, sharedErrors.CodeInternal, "could not start playback")
		return
	}
	c.JSON(http.StatusOK, gin.H{"session": sess})
}

func (h *PlaybackHandler) pause(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		httpx.JSONError(c, http.StatusUnauthorized, sharedErrors.CodeUnauthorized, "missing user context")
		return
	}
	var req positionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.JSONError(c, http.StatusBadRequest, sharedErrors.CodeInvalidRequest, "invalid payload")
		return
	}
	sess, err := h.svc.Pause(c.Request.Context(), userID, req.PositionMs)
	if err != nil {
		if err == service.ErrSessionNotFound {
			httpx.JSONError(c, http.StatusNotFound, sharedErrors.CodeNotFound, "session not found")
			return
		}
		httpx.JSONError(c, http.StatusInternalServerError, sharedErrors.CodeInternal, "could not pause")
		return
	}
	c.JSON(http.StatusOK, gin.H{"session": sess})
}

func (h *PlaybackHandler) resume(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		httpx.JSONError(c, http.StatusUnauthorized, sharedErrors.CodeUnauthorized, "missing user context")
		return
	}
	sess, err := h.svc.Resume(c.Request.Context(), userID)
	if err != nil {
		if err == service.ErrSessionNotFound {
			httpx.JSONError(c, http.StatusNotFound, sharedErrors.CodeNotFound, "session not found")
			return
		}
		httpx.JSONError(c, http.StatusInternalServerError, sharedErrors.CodeInternal, "could not resume")
		return
	}
	c.JSON(http.StatusOK, gin.H{"session": sess})
}

func (h *PlaybackHandler) seek(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		httpx.JSONError(c, http.StatusUnauthorized, sharedErrors.CodeUnauthorized, "missing user context")
		return
	}
	var req positionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.JSONError(c, http.StatusBadRequest, sharedErrors.CodeInvalidRequest, "invalid payload")
		return
	}
	sess, err := h.svc.Seek(c.Request.Context(), userID, req.PositionMs)
	if err != nil {
		if err == service.ErrSessionNotFound {
			httpx.JSONError(c, http.StatusNotFound, sharedErrors.CodeNotFound, "session not found")
			return
		}
		httpx.JSONError(c, http.StatusInternalServerError, sharedErrors.CodeInternal, "could not seek")
		return
	}
	c.JSON(http.StatusOK, gin.H{"session": sess})
}

func (h *PlaybackHandler) current(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		httpx.JSONError(c, http.StatusUnauthorized, sharedErrors.CodeUnauthorized, "missing user context")
		return
	}
	sess, err := h.svc.Current(c.Request.Context(), userID)
	if err != nil {
		if err == service.ErrSessionNotFound {
			httpx.JSONError(c, http.StatusNotFound, sharedErrors.CodeNotFound, "session not found")
			return
		}
		httpx.JSONError(c, http.StatusInternalServerError, sharedErrors.CodeInternal, "could not load session")
		return
	}
	c.JSON(http.StatusOK, gin.H{"session": sess})
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
