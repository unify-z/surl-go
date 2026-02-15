package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/unify-z/go-surl/internal/repository/user"
)

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		userInterface, exists := c.Get("User")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: User not login"})
			c.Abort()
			return
		}
		_user, ok := userInterface.(*user.UserModel)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			c.Abort()
			return
		}
		if !_user.IsAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Admin access permission required"})
			c.Abort()
			return
		}
		c.Next()
	}
}
