package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireRole retorna um middleware que restringe acesso por role.
// Exemplo: RequireRole("admin") permite apenas usuarios com role "admin".
func RequireRole(allowed ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden: role not found"})
			return
		}

		roleStr, ok := role.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden: invalid role"})
			return
		}

		for _, a := range allowed {
			if roleStr == a {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
	}
}
