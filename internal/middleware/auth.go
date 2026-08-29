package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/steven230500/cartagena-api/internal/service"
)

const UserIDContextKey = "user_id"

// RequireUser exige "Authorization: Bearer <jwt>" válido — sistema de auth
// de usuarios, separado por completo de RequireAdminKey (Fase 3).
func RequireUser(authSvc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		token, ok := strings.CutPrefix(header, "Bearer ")
		if !ok || token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		userID, err := authSvc.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		c.Set(UserIDContextKey, userID)
		c.Next()
	}
}
