package middlewares

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/unify-z/go-surl/internal/config"
	intlerror "github.com/unify-z/go-surl/internal/errors"
	"github.com/unify-z/go-surl/internal/handlers"
	UserService "github.com/unify-z/go-surl/internal/service/user"
	"github.com/unify-z/go-surl/internal/utils"
)

func AuthMiddleware(userSvc *UserService.UserManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if path == "/api/surl/create" && config.ConfigManagerInstance.Config.Site.AllowGuestToCreateURL {
			var authHeader string = c.GetHeader("Authorization")
			if authHeader == "" {
				c.Set("user_id", uint(0))
				c.Next()
				return
			}
			var jwtToken string = authHeader[len("Bearer "):]
			verifyResult, err := userSvc.ValidateJWTToken(jwtToken)
			if err != nil {
				if utils.IsInternalServerError(err) {
					handlers.Fail(c, http.StatusInternalServerError, "Internal server error")
					c.Abort()
					return
				}
				if errors.As(err, &intlerror.ErrUserBanned) {
					handlers.Fail(c, http.StatusForbidden, "Forbidden: "+err.Error())
					c.Abort()
					return
				}
				handlers.Fail(c, http.StatusUnauthorized, "Unauthorized: "+err.Error())
				c.Abort()
				return
			}
			c.Set("user_id", verifyResult.ID)
			c.Set("User", verifyResult)
			c.Next()
			return
		}
		// 其他路由要求认证
		var authHeader string = c.GetHeader("Authorization")
		if authHeader == "" {
			handlers.Fail(c, http.StatusUnauthorized, "Unauthorized: No token provided")
			c.Abort()
			return
		}
		var jwtToken string = authHeader[len("Bearer "):]
		verifyResult, err := userSvc.ValidateJWTToken(jwtToken)
		if err != nil {
			if utils.IsInternalServerError(err) {
				handlers.Fail(c, http.StatusInternalServerError, "Internal server error")
				c.Abort()
				return
			}
			if errors.As(err, &intlerror.ErrUserBanned) {
				handlers.Fail(c, http.StatusForbidden, "Forbidden: "+err.Error())
				c.Abort()
				return
			}
			handlers.Fail(c, http.StatusUnauthorized, "Unauthorized: "+err.Error())
			c.Abort()
			return
		}
		c.Set("user_id", verifyResult.ID)
		c.Set("User", verifyResult)
		c.Next()
	}
}
