package handler

import (
	"net/http"
	"strconv"

	sharedErrors "github.com/freedom-music/shared/errors"
	"github.com/freedom-music/shared/httpx"
	"github.com/freedom-music/shared/middleware"
	"github.com/freedom-music/user-service/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc *service.UserService
}

type updateRequest struct {
	Name    string `json:"name" binding:"omitempty"`
	Country string `json:"country" binding:"omitempty"`
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) RegisterRoutes(r *gin.Engine, jwtSecret string) {
	r.GET("/health", httpx.Health)

	api := r.Group("/api/v1")
	api.Use(middleware.AuthRequired(jwtSecret))
	{
		api.GET("/users/me", h.me)
		api.PUT("/users/me", h.update)
		api.GET("/users/me/preferences", h.preferences)
	}
}

func (h *UserHandler) me(c *gin.Context) {
	userIDStr, ok := c.Get(middleware.UserIDContextKey)
	if !ok {
		httpx.JSONError(c, http.StatusUnauthorized, sharedErrors.CodeUnauthorized, "missing user context")
		return
	}
	userID, err := strconv.ParseInt(userIDStr.(string), 10, 64)
	if err != nil {
		httpx.JSONError(c, http.StatusUnauthorized, sharedErrors.CodeUnauthorized, "invalid user id")
		return
	}

	profile, err := h.svc.Get(c.Request.Context(), userID)
	if err != nil {
		httpx.JSONError(c, http.StatusInternalServerError, sharedErrors.CodeInternal, "could not load profile")
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": profile})
}

func (h *UserHandler) update(c *gin.Context) {
	userIDStr, ok := c.Get(middleware.UserIDContextKey)
	if !ok {
		httpx.JSONError(c, http.StatusUnauthorized, sharedErrors.CodeUnauthorized, "missing user context")
		return
	}
	userID, err := strconv.ParseInt(userIDStr.(string), 10, 64)
	if err != nil {
		httpx.JSONError(c, http.StatusUnauthorized, sharedErrors.CodeUnauthorized, "invalid user id")
		return
	}

	var req updateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.JSONError(c, http.StatusBadRequest, sharedErrors.CodeInvalidRequest, "invalid payload")
		return
	}

	profile, err := h.svc.Update(c.Request.Context(), userID, req.Name, req.Country)
	if err != nil {
		httpx.JSONError(c, http.StatusInternalServerError, sharedErrors.CodeInternal, "could not update profile")
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": profile})
}

func (h *UserHandler) preferences(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"preferences": gin.H{"theme": "light", "explicit_filter": false}})
}
