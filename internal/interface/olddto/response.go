package olddto

import (
	"errors"

	"github.com/gin-gonic/gin"

	"jcourse_go/pkg/apperror"
)

type PaginationResponse struct {
	Page     int64 `json:"page"`
	PageSize int64 `json:"page_size"`
	Total    int64 `json:"total"`
}
type APIResponse struct {
	Code       int                 `json:"code"`
	Msg        string              `json:"msg"`
	Data       any                 `json:"data,omitempty"`
	Pagination *PaginationResponse `json:"pagination,omitempty"`
}

func ResponseWithData(c *gin.Context, data any) {
	c.JSON(200, APIResponse{
		Code: 0,
		Msg:  "",
		Data: data,
	})
}

func ResponseWithPaginationData(c *gin.Context, data any, pagination *PaginationResponse) {
	c.JSON(200, APIResponse{
		Code:       0,
		Msg:        "",
		Data:       data,
		Pagination: pagination,
	})
}

func ResponseWithError(c *gin.Context, err error) {
	var appError *apperror.AppError
	if !errors.As(err, &appError) {
		appError = apperror.ErrSomethingWrong
	}
	c.JSON(appError.StatusCode(), APIResponse{Code: appError.Code, Msg: appError.Msg})
	return
}
