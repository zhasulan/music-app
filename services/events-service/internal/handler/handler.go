package handler

import (
	"net/http"
	"strconv"

	"github.com/freedom-music/events-service/internal/domain"
	"github.com/freedom-music/events-service/internal/repository"
	sharedErrors "github.com/freedom-music/shared/errors"
	"github.com/freedom-music/shared/httpx"
	"github.com/freedom-music/shared/middleware"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo      *repository.Repository
	jwtSecret string
}

func New(repo *repository.Repository, jwtSecret string) *Handler {
	return &Handler{repo: repo, jwtSecret: jwtSecret}
}

func (h *Handler) Register(r *gin.Engine) {
	r.GET("/health", httpx.Health)

	api := r.Group("/api/v1/events")
	api.Use(middleware.AuthRequired(h.jwtSecret))
	{
		api.POST("", h.generic)
		api.POST("/track-played", h.trackPlayed)
		api.POST("/track-paused", h.trackPaused)
		api.POST("/track-liked", h.trackLiked)
		api.POST("/search", h.searchPerformed)
		api.POST("/playlist-created", h.playlistCreated)
		api.POST("/playlist-track-added", h.playlistTrackAdded)
	}
}

type eventRequest struct {
	TrackID    string `json:"track_id" binding:"omitempty"`
	PlaylistID int64  `json:"playlist_id" binding:"omitempty"`
	Query      string `json:"query" binding:"omitempty"`
	EventType  string `json:"event_type" binding:"required"`
}

func (h *Handler) generic(c *gin.Context) {
	h.handle(c, eventRequest{})
}

func (h *Handler) trackPlayed(c *gin.Context)        { h.handleTyped(c, "track_played") }
func (h *Handler) trackPaused(c *gin.Context)        { h.handleTyped(c, "track_paused") }
func (h *Handler) trackLiked(c *gin.Context)         { h.handleTyped(c, "track_liked") }
func (h *Handler) searchPerformed(c *gin.Context)    { h.handleTyped(c, "search_performed") }
func (h *Handler) playlistCreated(c *gin.Context)    { h.handleTyped(c, "playlist_created") }
func (h *Handler) playlistTrackAdded(c *gin.Context) { h.handleTyped(c, "playlist_track_added") }

func (h *Handler) handleTyped(c *gin.Context, t string) {
	h.handle(c, eventRequest{EventType: t})
}

func (h *Handler) handle(c *gin.Context, base eventRequest) {
	var req eventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.JSONError(c, http.StatusBadRequest, sharedErrors.CodeInvalidRequest, "invalid payload")
		return
	}
	if base.EventType != "" {
		req.EventType = base.EventType
	}
	if req.EventType == "" {
		httpx.JSONError(c, http.StatusBadRequest, sharedErrors.CodeInvalidRequest, "event_type required")
		return
	}

	var userID *int64
	if val, ok := c.Get(middleware.UserIDContextKey); ok {
		if s, ok2 := val.(string); ok2 {
			if parsed, err := strconv.ParseInt(s, 10, 64); err == nil {
				userID = &parsed
			}
		}
	}

	var trackID *string
	if req.TrackID != "" {
		trackID = &req.TrackID
	}
	var playlistID *int64
	if req.PlaylistID != 0 {
		playlistID = &req.PlaylistID
	}
	var query *string
	if req.Query != "" {
		query = &req.Query
	}

	ev := &domain.Event{
		UserID:     userID,
		EventType:  req.EventType,
		TrackID:    trackID,
		PlaylistID: playlistID,
		Query:      query,
	}
	if err := h.repo.Insert(c.Request.Context(), ev); err != nil {
		httpx.JSONError(c, http.StatusInternalServerError, sharedErrors.CodeInternal, "failed to persist event")
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
