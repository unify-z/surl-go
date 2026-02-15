package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/unify-z/go-surl/internal/dto"
)

func Success[T any](c *gin.Context, data T) {
	c.JSON(http.StatusOK, dto.CommonDTO{
		Code:    200,
		Message: "ok",
		Data:    data,
	})
}

func Fail(c *gin.Context, httpStatus int, msg string) {
	c.JSON(httpStatus, dto.CommonDTO{
		Code:    httpStatus,
		Message: msg,
	})
}
