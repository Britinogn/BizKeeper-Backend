package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminOnly(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get role from context (you must set it earlier)
		role, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Role not found",
			})
			return
		}

		roleStr, ok := role.(string)
		if !ok || roleStr != requiredRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "access denied",
			})
			return
		}

		c.Next()
	}
}