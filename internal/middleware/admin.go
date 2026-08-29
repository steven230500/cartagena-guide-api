package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// RequireAdminKey exige el header X-Admin-Key con el valor de ADMIN_API_KEY.
// Sin sistema de usuarios todavía — ver ARCHITECTURE.md, Fase 1.
func RequireAdminKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		expected := os.Getenv("ADMIN_API_KEY")
		got := c.GetHeader("X-Admin-Key")

		if expected == "" || got == "" || got != expected {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		c.Next()
	}
}
