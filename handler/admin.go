package handler

import (
	"jcourse_go/model/domain"
	"jcourse_go/model/dto"
	"jcourse_go/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminGetUserList(c *gin.Context) {
	var request dto.UserListRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{Message: "参数错误"})
		return
	}

	filter := domain.UserFilter{
		Page:     request.Page,
		PageSize: request.PageSize,
	}
	users, err := service.AdminGetUserList(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{Message: "内部错误。"})
	}
	total, _ := service.GetUserCount(c, filter)
	response := dto.UserListResponseForAdmin{
		Page:     request.Page,
		PageSize: request.PageSize,
		Total:    total,
		Data:     users,
	}
	c.JSON(http.StatusOK, response)
}