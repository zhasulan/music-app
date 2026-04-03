package handler

import (
	"net/http"

	"github.com/freedom-music/auth-service/internal/service"
	sharedErrors "github.com/freedom-music/shared/errors"
	"github.com/freedom-music/shared/httpx"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	svc *service.AuthService
}

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/health", httpx.Health)
	api := r.Group("/api/v1")
	{
		api.POST("/auth/register", h.register)
		api.POST("/auth/login", h.login)
		api.POST("/auth/refresh", h.refresh)
		api.POST("/auth/logout", h.logout)
	}
}

func (h *AuthHandler) register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.JSONError(c, http.StatusBadRequest, sharedErrors.CodeInvalidRequest, "invalid payload")
		return
	}

	user, access, refresh, err := h.svc.Register(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		switch err {
		case service.ErrEmailExists:
			httpx.JSONError(c, http.StatusConflict, sharedErrors.CodeInvalidRequest, err.Error())
		default:
			httpx.JSONError(c, http.StatusInternalServerError, sharedErrors.CodeInternal, "could not register")
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user":          user,
		"access_token":  access,
		"refresh_token": refresh,
	})
}

func (h *AuthHandler) login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.JSONError(c, http.StatusBadRequest, sharedErrors.CodeInvalidRequest, "invalid payload")
		return
	}

	user, access, refresh, err := h.svc.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		switch err {
		case service.ErrInvalidLogin:
			httpx.JSONError(c, http.StatusUnauthorized, sharedErrors.CodeUnauthorized, "invalid credentials")
		default:
			httpx.JSONError(c, http.StatusInternalServerError, sharedErrors.CodeInternal, "could not login")
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":          user,
		"access_token":  access,
		"refresh_token": refresh,
	})
}

func (h *AuthHandler) refresh(c *gin.Context) {
	var req refreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.JSONError(c, http.StatusBadRequest, sharedErrors.CodeInvalidRequest, "invalid payload")
		return
	}

	access, refresh, err := h.svc.Refresh(req.RefreshToken)
	if err != nil {
		switch err {
		case service.ErrInvalidToken:
			httpx.JSONError(c, http.StatusUnauthorized, sharedErrors.CodeUnauthorized, "invalid refresh token")
		default:
			httpx.JSONError(c, http.StatusInternalServerError, sharedErrors.CodeInternal, "could not refresh")
		}
		return
	}

	c.JSON(http.StatusOK, tokenResponse{AccessToken: access, RefreshToken: refresh})
}

func (h *AuthHandler) logout(c *gin.Context) {
	// Placeholder for future token revocation. For now just acknowledge.
	c.JSON(http.StatusOK, gin.H{"message": "logout acknowledged"})
}
