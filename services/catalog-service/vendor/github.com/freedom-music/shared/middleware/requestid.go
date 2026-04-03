package middleware

import (
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

const RequestIDKey = "request_id"

// RequestID sets a request id header if missing and stores it on the context.
func RequestID() gin.HandlerFunc {
    return func(c *gin.Context) {
        rid := c.GetHeader("X-Request-ID")
        if rid == "" {
            rid = uuid.NewString()
        }
        c.Set(RequestIDKey, rid)
        c.Writer.Header().Set("X-Request-ID", rid)
        c.Next()
    }
}
