package handler

import (
	"log"
	"net/http"

	"jcourse_go/model/types"
	"jcourse_go/service"

	"github.com/gin-gonic/gin"

	"jcourse_go/model/dto"
)

func AdminGetUserList(c *gin.Context) {
	var request dto.UserListRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{Message: "参数错误"})
		return
	}
	/*
		filter := dto.UserFilter{
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
	*/
}
func AdminChangeUserPoint(c *gin.Context) {
	var request dto.ChangeUserPointRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{Message: "参数错误"})
		return
	}
	err := service.ChangeUserPoints(c, request.UserID, types.PointEventAdminChange, request.Value, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{Message: "用户积分更新失败。"})
		log.Printf("ChangeUserPointHandler: %v", err)
		return
	}
	c.JSON(http.StatusOK, dto.BaseResponse{Message: "用户积分更新成功。"})
}
