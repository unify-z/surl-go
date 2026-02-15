package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/unify-z/go-surl/internal/dto/admin_dto"
	"github.com/unify-z/go-surl/internal/handlers"
	ShortURLService "github.com/unify-z/go-surl/internal/service/surl"
	"github.com/unify-z/go-surl/internal/utils"
)

type AdminSURLHandler struct {
	surlSvc *ShortURLService.ShortURLManager
}

func NewAdminSURLHandler(surlSvc *ShortURLService.ShortURLManager) *AdminSURLHandler {
	return &AdminSURLHandler{surlSvc: surlSvc}
}

func (ash *AdminSURLHandler) DeleteSURL(c *gin.Context) {
	var req admin_dto.AdminDeleteShortURLReq
	if err := c.ShouldBindJSON(&req); err != nil {
		if utils.IsInternalServerError(err) {
			handlers.Fail(c, http.StatusInternalServerError, "Internal server error")
			return
		}
		handlers.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	err := ash.surlSvc.DeleteShortURL(req.ID)
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

func (ash *AdminSURLHandler) ListShortURL(c *gin.Context) {
	var req admin_dto.AdminListShortURLsReq
	if err := c.ShouldBindQuery(&req); err != nil {
		if utils.IsInternalServerError(err) {
			handlers.Fail(c, http.StatusInternalServerError, "Internal server error")
			return
		}
		handlers.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	shortURLs, err := ash.surlSvc.ListAllShortURLsWithPagination(req.Page, req.PageSize)
	if err != nil {
		if utils.IsInternalServerError(err) {
			handlers.Fail(c, http.StatusInternalServerError, "Internal server error")
			return
		}
		handlers.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	count, err := ash.surlSvc.CountShortURLs()
	if err != nil {
		if utils.IsInternalServerError(err) {
			handlers.Fail(c, http.StatusInternalServerError, "Internal server error")
			return
		}
		handlers.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	handlers.Success[map[string]any](c, map[string]any{
		"surls":       shortURLs,
		"total_count": count,
	})
}

func (ash *AdminSURLHandler) EditShortURL(c *gin.Context) {
	var req admin_dto.AdminEditShortURLReq
	if err := c.ShouldBindJSON(&req); err != nil {
		if utils.IsInternalServerError(err) {
			handlers.Fail(c, http.StatusInternalServerError, "Internal server error")
			return
		}
		handlers.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	surlData, err := ash.surlSvc.GetShortURLByShortCode(req.ShortCode)
	if err != nil {
		if utils.IsInternalServerError(err) {
			handlers.Fail(c, http.StatusInternalServerError, "Internal server error")
			return
		}
		handlers.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	surlData.OriginalURL = req.OriginalURL
	surlData.UserID = utils.ToUint(req.UserID)
	err = ash.surlSvc.UpdateShortURL(surlData)
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
