package httpx

import (
	"net/http"

	sharedErrors "github.com/freedom-music/shared/errors"
	"github.com/gin-gonic/gin"
)

func JSONError(c *gin.Context, status int, code sharedErrors.ErrorCode, message string) {
	c.AbortWithStatusJSON(status, sharedErrors.New(code, message))
}

func JSONOK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
