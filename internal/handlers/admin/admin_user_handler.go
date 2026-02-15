package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/unify-z/go-surl/internal/dto/admin_dto"
	"github.com/unify-z/go-surl/internal/handlers"
	"github.com/unify-z/go-surl/internal/repository/user"
	UserService "github.com/unify-z/go-surl/internal/service/user"
	"github.com/unify-z/go-surl/internal/utils"
)

type AdminUserHandler struct {
	userSvc *UserService.UserManager
}

func NewAdminUserHandler(userSvc *UserService.UserManager) *AdminUserHandler {
	return &AdminUserHandler{userSvc: userSvc}
}

func (auh *AdminUserHandler) DeleteUser(c *gin.Context) {
	var req admin_dto.AdminDeleteUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		if utils.IsInternalServerError(err) {
			handlers.Fail(c, http.StatusInternalServerError, "Internal server error")
			return
		}
		handlers.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	err := auh.userSvc.DeleteUser(req.ID)
	if err != nil {
		if utils.IsInternalServerError(err) {
			handlers.Fail(c, http.StatusInternalServerError, "Internal server error")
			return
		}
		handlers.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	handlers.Success[any](c, nil)
}

func (auh *AdminUserHandler) GetUser(c *gin.Context) {
	users, err := auh.userSvc.ListUsersWithoutSURLs()
	if err != nil {
		if utils.IsInternalServerError(err) {
			handlers.Fail(c, http.StatusInternalServerError, "Internal server error")
			return
		}
		handlers.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	handlers.Success[[]user.UserModel](c, users)
}

func (auh *AdminUserHandler) BanUser(c *gin.Context) {
	var req admin_dto.AdminBanUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		if utils.IsInternalServerError(err) {
			handlers.Fail(c, http.StatusInternalServerError, "Internal server error")
			return
		}
		handlers.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	_user, err := auh.userSvc.GetUserByID(req.UserID)
	if err != nil {
		if utils.IsInternalServerError(err) {
			handlers.Fail(c, http.StatusInternalServerError, "Internal server error")
			return
		}
		handlers.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	_user.IsBanned = true
	updateerr := auh.userSvc.UpdateUser(_user)
	if updateerr != nil {
		if utils.IsInternalServerError(err) {
			handlers.Fail(c, http.StatusInternalServerError, "Internal server error")
			return
		}
		handlers.Fail(c, http.StatusBadRequest, updateerr.Error())
		return
	}
	handlers.Success[any](c, nil)
}

func (auh *AdminUserHandler) UnbanUser(c *gin.Context) {
	var req admin_dto.AdminUnbanUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		if utils.IsInternalServerError(err) {
			handlers.Fail(c, http.StatusInternalServerError, "Internal server error")
			return
		}
		handlers.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	_user, err := auh.userSvc.GetUserByID(req.UserID)
	if err != nil {
		if utils.IsInternalServerError(err) {
			handlers.Fail(c, http.StatusInternalServerError, "Internal server error")
			return
		}
		handlers.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	_user.IsBanned = false
	updateerr := auh.userSvc.UpdateUser(_user)
	if updateerr != nil {
		if utils.IsInternalServerError(err) {
			handlers.Fail(c, http.StatusInternalServerError, "Internal server error")
			return
		}
		handlers.Fail(c, http.StatusBadRequest, updateerr.Error())
		return
	}
	handlers.Success[any](c, nil)
}
