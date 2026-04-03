package middleware

import (
	"net/http"
	"strings"

	sharedAuth "github.com/freedom-music/shared/auth"
	sharedErrors "github.com/freedom-music/shared/errors"
	"github.com/gin-gonic/gin"
)

const UserIDContextKey = "user_id"

// AuthRequired validates bearer tokens and injects subject (user id) into context.
func AuthRequired(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			unauthorized(c)
			return
		}
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			unauthorized(c)
			return
		}
		claims, err := sharedAuth.ParseToken(parts[1], secret)
		if err != nil {
			unauthorized(c)
			return
		}
		c.Set(UserIDContextKey, claims.Subject)
		c.Next()
	}
}

func unauthorized(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, sharedErrors.New(sharedErrors.CodeUnauthorized, "unauthorized"))
}
