package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuditTrail logs API accesses.
func AuditTrail() gin.HandlerFunc {
	return func(c *gin.Context) {
		auditID := uuid.NewString()
		c.Set("auditID", auditID)
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		user := UserFrom(c.Request.Context())
		username := "anonymous"
		if user != nil {
			username = user.Username
		}
		log.Printf("AUDIT id=%s user=%s method=%s path=%s status=%d duration=%s", auditID, username, c.Request.Method, c.FullPath(), c.Writer.Status(), duration)
	}
}
