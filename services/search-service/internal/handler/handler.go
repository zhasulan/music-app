package handler

import (
	"net/http"

	"github.com/freedom-music/search-service/internal/service"
	"github.com/freedom-music/shared/httpx"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *service.Service
}

func New(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(r *gin.Engine) {
	r.GET("/health", httpx.Health)
	api := r.Group("/api/v1/search")
	{
		api.GET("", h.search)
		api.GET("/suggest", h.suggest)
	}
}

func (h *Handler) search(c *gin.Context) {
	q := c.Query("q")
	results := h.svc.Search(q)
	c.JSON(http.StatusOK, gin.H{"results": results})
}

func (h *Handler) suggest(c *gin.Context) {
	q := c.Query("q")
	results := h.svc.Suggest(q)
	c.JSON(http.StatusOK, gin.H{"suggestions": results})
}
