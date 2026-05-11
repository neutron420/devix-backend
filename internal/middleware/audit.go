package middleware

import (
	"devix-backend/internal/modules/audit"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Audit(auditSvc *audit.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process request
		c.Next()

		// Only log mutating requests or specific statuses
		if c.Request.Method == "GET" {
			return
		}

		actorIDVal, exists := c.Get(ContextKeyUserID)
		var actorID uuid.UUID
		if exists {
			actorID = actorIDVal.(uuid.UUID)
		}

		// Log the action
		auditSvc.Log(
			c.Request.Context(),
			actorID,
			c.Request.Method,
			c.Request.URL.Path,
			"", // TargetID could be extracted from path if needed
			fmt.Sprintf("Status: %d", c.Writer.Status()),
			c.ClientIP(),
			c.GetHeader("User-Agent"),
		)
	}
}
