package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/unify-z/go-surl/internal/repository/user"
)

func SafeGetUserId(c *gin.Context) uint {
	userIdInterface, exists := c.Get("user_id")
	if !exists {
		return 0
	}
	userId, ok := userIdInterface.(uint)
	if !ok {
		return 0
	}
	return userId
}

func SafeGetUserData(c *gin.Context) (user.UserModel, bool) {
	v, exists := c.Get("User")
	if !exists {
		return user.UserModel{}, false
	}

	_user, ok := v.(user.UserModel)
	if !ok {
		return user.UserModel{}, false
	}

	return _user, true
}
