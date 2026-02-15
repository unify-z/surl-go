package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/unify-z/go-surl/internal/config"
	"github.com/unify-z/go-surl/internal/dto"
	UserService "github.com/unify-z/go-surl/internal/service/user"
	"github.com/unify-z/go-surl/internal/service/verify_code"
	"github.com/unify-z/go-surl/internal/utils"
)

type UserHandler struct {
	userSvc       *UserService.UserManager
	verifyCodeSvc *verify_code.VerifyCodeManager
}

func NewUserHandler(userSvc *UserService.UserManager, verifyCodeSvc *verify_code.VerifyCodeManager) *UserHandler {
	return &UserHandler{userSvc: userSvc, verifyCodeSvc: verifyCodeSvc}
}

func (uh *UserHandler) GetUserInfo(c *gin.Context) {
	userId := utils.SafeGetUserId(c)
	if userId == 0 {
		Fail(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	user, err := uh.userSvc.GetUserByID(userId)
	if err != nil {
		if utils.IsInternalServerError(err) {
			Fail(c, http.StatusInternalServerError, "Internal server error")
			return
		}
	}
	Success[map[string]any](c, map[string]any{
		"user_id":  strconv.Itoa(int(user.ID)),
		"username": user.Username,
		"is_admin": user.IsAdmin,
	})
}

func (uh *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		if utils.IsInternalServerError(err) {
			Fail(c, http.StatusInternalServerError, "Internal server error")
			return
		}
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	if !config.ConfigManagerInstance.Config.Site.AllowRegistration {
		Fail(c, http.StatusForbidden, "Registration is disabled")
		return
	}
	user, err := uh.userSvc.RegisterUser(req.Username, req.PasswordMD5, req.Email, req.EmailVerifyCode, false)
	if err != nil {
		if utils.IsInternalServerError(err) {
			Fail(c, http.StatusInternalServerError, "Internal server error")
			return
		}
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	_, _, token := uh.userSvc.LoginUser(req.Username, req.PasswordMD5)
	Success[map[string]any](c, map[string]any{
		"user_id":  strconv.Itoa(int(user.ID)),
		"username": user.Username,
		"is_admin": user.IsAdmin,
		"token":    token,
	})
}

func (uh *UserHandler) Login(c *gin.Context) {
	var req dto.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		if utils.IsInternalServerError(err) {
			Fail(c, http.StatusInternalServerError, "Internal server error")
			return
		}
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	user, err, token := uh.userSvc.LoginUser(req.Username, req.PasswordMD5)
	if err != nil {
		if utils.IsInternalServerError(err) {
			Fail(c, http.StatusInternalServerError, "Internal server error")
			return
		}
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	Success[map[string]any](c, map[string]any{
		"user_id":  strconv.Itoa(int(user.ID)),
		"username": user.Username,
		"is_admin": user.IsAdmin,
		"token":    token,
	})
}

func (uh *UserHandler) SendEmailVerifyCode(c *gin.Context) {
	var req dto.SendEmailVerifyCodeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		if utils.IsInternalServerError(err) {
			Fail(c, http.StatusInternalServerError, "Internal server error")
			return
		}
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	err := uh.verifyCodeSvc.SendCode(req.Email)
	if err != nil {
		if utils.IsInternalServerError(err) {
			Fail(c, http.StatusInternalServerError, "Internal server error")
			return
		}
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	Success[any](c, nil)
}
