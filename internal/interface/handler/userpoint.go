package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"jcourse_go/internal/domain/user"
	"jcourse_go/internal/interface/dto"
	"jcourse_go/internal/interface/middleware"
	"jcourse_go/internal/service"
)

func TransferUserPointHandler(c *gin.Context) {
	var request olddto.TransferUserPointRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, olddto.BaseResponse{Message: "参数错误"})
		return
	}
	user := middleware.GetCurrentUser(c)
	if user.ID == request.Receiver {
		c.JSON(http.StatusBadRequest, olddto.BaseResponse{Message: "不能给自己转账。"})
		return
	}
	err := service.TransferUserPoints(c, user.ID, request.Receiver, request.Value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, olddto.BaseResponse{Message: "用户积分转账失败。"})
		return
	}
	c.JSON(http.StatusOK, olddto.BaseResponse{Message: "用户积分转账成功。"})
}
func AdminTransferUserPoint(c *gin.Context) {
	var request olddto.TransferUserPointAdminRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, olddto.BaseResponse{Message: "参数错误"})
		return
	}
	if request.Sender == request.Receiver {
		c.JSON(http.StatusBadRequest, olddto.BaseResponse{Message: "不能给自己转账。"})
		return
	}
	err := service.TransferUserPoints(c, request.Sender, request.Receiver, request.Value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, olddto.BaseResponse{Message: "用户积分转账失败。"})
		return
	}
	c.JSON(http.StatusOK, olddto.BaseResponse{Message: "用户积分转账成功。"})
}

func GetUserPointDetailListHandler(c *gin.Context) {
	var request olddto.UserPointDetailListRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, olddto.BaseResponse{Message: "参数错误"})
		return
	}
	user := middleware.GetCurrentUser(c)
	filter := user.UserPointDetailFilter{
		UserID:    user.ID,
		StartTime: request.StartTime,
		EndTime:   request.EndTime,
	}
	totalValue, userPointDetails, err := service.GetUserPointDetailList(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, olddto.BaseResponse{Message: "内部错误。"})
	}

	total, _ := service.GetUserPointDetailCount(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, olddto.BaseResponse{Message: "内部错误。"})
	}
	response := olddto.UserPointDetailListResponse{
		BasePaginateResponse: olddto.BasePaginateResponse[user.UserPointDetailItem]{
			Page:     request.Page,
			PageSize: request.PageSize,
			Total:    total,
			Data:     userPointDetails,
		},
		CurrentPoint: totalValue,
	}
	c.JSON(http.StatusOK, response)
}

func AdminGetUserPointDetailList(c *gin.Context) {
	var request olddto.UserPointDetailListAdminRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, olddto.BaseResponse{Message: "参数错误"})
		return
	}
	filter := user.UserPointQuery{
		UserID:    request.UserID,
		StartTime: request.StartTime,
		EndTime:   request.EndTime,
	}
	_, userPointDetails, err := service.GetUserPointDetailList(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, olddto.BaseResponse{Message: "内部错误。"})
	}

	total, _ := service.GetUserPointDetailCount(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, olddto.BaseResponse{Message: "内部错误。"})
	}
	response := olddto.UserPointDetailListResponse{
		BasePaginateResponse: olddto.BasePaginateResponse[user.UserPointDetailItem]{
			Page:     request.Page,
			PageSize: request.PageSize,
			Total:    total,
			Data:     userPointDetails,
		},
	}
	c.JSON(http.StatusOK, response)
}
