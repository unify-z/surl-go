package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/unify-z/go-surl/internal/dto"
	"github.com/unify-z/go-surl/internal/repository/surl"
	ShortURLService "github.com/unify-z/go-surl/internal/service/surl"
	"github.com/unify-z/go-surl/internal/utils"
)

type SURLHandler struct {
	surlSvc *ShortURLService.ShortURLManager
}

func NewSURLHandler(surlSvc *ShortURLService.ShortURLManager) *SURLHandler {
	return &SURLHandler{surlSvc: surlSvc}
}

func (sh *SURLHandler) RedirectToOriginalURL(c *gin.Context) {
	shortCode := c.Param("short_code")
	originalURL, err := sh.surlSvc.GetOriginalURLFromShortCode(shortCode)
	if err != nil {
		if utils.IsInternalServerError(err) {
			Fail(c, http.StatusInternalServerError, "Internal server error")
			return
		}
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	c.Redirect(http.StatusFound, originalURL)
}

func (sh *SURLHandler) CreateShortURL(c *gin.Context) {
	var req dto.CreateShortURLReq
	if err := c.ShouldBindJSON(&req); err != nil {
		if utils.IsInternalServerError(err) {
			Fail(c, http.StatusInternalServerError, "Internal server error")
			return
		}
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		Fail(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userID, ok := userIDInterface.(uint)
	if !ok {
		Fail(c, http.StatusInternalServerError, "Internal server error")
		return
	}
	shortURL, err := sh.surlSvc.CreateShortURL(req.OriginalURL, userID)
	if err != nil {
		if utils.IsInternalServerError(err) {
			Fail(c, http.StatusInternalServerError, "Internal server error")
			return
		}
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	Success[map[string]any](c, map[string]any{
		"short_code":   shortURL.ShortCode,
		"original_url": shortURL.OriginalURL,
	})
}

func (sh *SURLHandler) ListUserShortURLs(c *gin.Context) {
	userId := utils.SafeGetUserId(c)
	if userId == 0 {
		Fail(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	urls, err := sh.surlSvc.ListShortURLsByUserID(userId)
	if err != nil {
		if utils.IsInternalServerError(err) {
			Fail(c, http.StatusInternalServerError, "Internal server error")
			return
		}
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	Success[[]surl.ShortURL](c, urls)
}

func (sh *SURLHandler) DeleteShortURL(c *gin.Context) {
	var req dto.DeleteShortURLReq
	if err := c.ShouldBindJSON(&req); err != nil {
		if utils.IsInternalServerError(err) {
			Fail(c, http.StatusInternalServerError, "Internal server error")
			return
		}
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	urlData, err := sh.surlSvc.GetShortURLByShortCode(req.ShortCode)
	if err != nil {
		if utils.IsInternalServerError(err) {
			Fail(c, http.StatusInternalServerError, "Internal server error")
			return
		}
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	userId := utils.SafeGetUserId(c)
	if userId == 0 {
		Fail(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	if urlData.UserID != userId {
		Fail(c, http.StatusForbidden, "Forbidden: You dont have enough permission to do this action")
	}
	err = sh.surlSvc.DeleteShortURL(urlData.ID)
	if err != nil {
		if utils.IsInternalServerError(err) {
			Fail(c, http.StatusInternalServerError, "Internal server error")
		}
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	Success[any](c, nil)
}

func (sh *SURLHandler) EditShortURL(c *gin.Context) {
	var req dto.EditShortURLReq
	if err := c.ShouldBindJSON(&req); err != nil {
		if utils.IsInternalServerError(err) {
			Fail(c, http.StatusInternalServerError, "Internal server error")
			return
		}
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	urlData, err := sh.surlSvc.GetShortURLByShortCode(req.ShortCode)
	if err != nil {
		if utils.IsInternalServerError(err) {
			Fail(c, http.StatusInternalServerError, "Internal server error")
			return
		}
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	userId := utils.SafeGetUserId(c)
	if userId == 0 {
		Fail(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	if urlData.UserID != userId {
		Fail(c, http.StatusForbidden, "Forbidden: You dont have enough permission to do this action")
	}
	urlData.OriginalURL = req.OriginalURL
	err = sh.surlSvc.UpdateShortURL(urlData)
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
