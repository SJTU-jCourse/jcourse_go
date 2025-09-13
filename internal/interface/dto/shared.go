package dto

import (
	"errors"

	"github.com/gin-gonic/gin"

	"jcourse_go/pkg/apperror"
)

type BaseResponse struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func WriteDataResponse(ctx *gin.Context, data any) {
	ctx.JSON(200, BaseResponse{
		Code:    0,
		Message: "",
		Data:    data,
	})
}

func WriteErrorResponse(ctx *gin.Context, err error) {
	var appErr *apperror.AppError
	if !errors.As(err, &appErr) {
		appErr = apperror.ErrSomethingWrong
	}
	ctx.JSON(appErr.StatusCode(), BaseResponse{
		Code:    appErr.Code,
		Message: appErr.Msg,
	})
}

func WriteBadArgumentResponse(ctx *gin.Context) {
	WriteErrorResponse(ctx, apperror.ErrBadRequest)
}
