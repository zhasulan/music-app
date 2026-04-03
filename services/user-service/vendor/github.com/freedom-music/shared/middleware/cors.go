package middleware

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

// CORS allows browser-based local development.
// For Phase 1 we allow all origins and common headers/methods.
func CORS() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Request-ID")
        c.Writer.Header().Set("Access-Control-Expose-Headers", "X-Request-ID")
        if c.Request.Method == http.MethodOptions {
            c.AbortWithStatus(http.StatusNoContent)
            return
        }
        c.Next()
    }
}
