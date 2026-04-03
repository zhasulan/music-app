package handler

import (
	"net/http"

	"github.com/freedom-music/recommendation-service/internal/service"
	sharedErrors "github.com/freedom-music/shared/errors"
	"github.com/freedom-music/shared/httpx"
	"github.com/freedom-music/shared/middleware"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc       *service.Service
	jwtSecret string
}

func New(svc *service.Service, jwtSecret string) *Handler {
	return &Handler{svc: svc, jwtSecret: jwtSecret}
}

func (h *Handler) Register(r *gin.Engine) {
	r.GET("/health", httpx.Health)

	api := r.Group("/api/v1/recommendations")
	{
		api.GET("/trending", h.trending)
		protected := api.Group("")
		protected.Use(middleware.AuthRequired(h.jwtSecret))
		{
			protected.GET("/recently-played", h.recent)
			protected.GET("/for-you", h.forYou)
		}
	}
}

func (h *Handler) trending(c *gin.Context) {
	tracks, _ := h.svc.Trending(c.Request.Context(), 20)
	c.JSON(http.StatusOK, gin.H{"tracks": tracks})
}

func (h *Handler) recent(c *gin.Context) {
	user, ok := userID(c)
	if !ok {
		httpx.JSONError(c, http.StatusUnauthorized, sharedErrors.CodeUnauthorized, "unauthorized")
		return
	}
	tracks, _ := h.svc.RecentlyPlayed(c.Request.Context(), user, 20)
	c.JSON(http.StatusOK, gin.H{"tracks": tracks})
}

func (h *Handler) forYou(c *gin.Context) {
	user, ok := userID(c)
	if !ok {
		httpx.JSONError(c, http.StatusUnauthorized, sharedErrors.CodeUnauthorized, "unauthorized")
		return
	}
	tracks, _ := h.svc.ForYou(c.Request.Context(), user, 20)
	c.JSON(http.StatusOK, gin.H{"tracks": tracks})
}

func userID(c *gin.Context) (int64, bool) {
	if v, ok := c.Get(middleware.UserIDContextKey); ok {
		if s, ok2 := v.(string); ok2 {
			return service.ParseUserID(s)
		}
	}
	return 0, false
}
