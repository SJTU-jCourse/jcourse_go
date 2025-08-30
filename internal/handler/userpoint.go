package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"jcourse_go/internal/middleware"
	dto2 "jcourse_go/internal/model/dto"
	"jcourse_go/internal/model/model"
	"jcourse_go/internal/service"
)

func TransferUserPointHandler(c *gin.Context) {
	var request dto2.TransferUserPointRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto2.BaseResponse{Message: "参数错误"})
		return
	}
	user := middleware.GetCurrentUser(c)
	if user.ID == request.Receiver {
		c.JSON(http.StatusBadRequest, dto2.BaseResponse{Message: "不能给自己转账。"})
		return
	}
	err := service.TransferUserPoints(c, user.ID, request.Receiver, request.Value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto2.BaseResponse{Message: "用户积分转账失败。"})
		return
	}
	c.JSON(http.StatusOK, dto2.BaseResponse{Message: "用户积分转账成功。"})
}
func AdminTransferUserPoint(c *gin.Context) {
	var request dto2.TransferUserPointAdminRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto2.BaseResponse{Message: "参数错误"})
		return
	}
	if request.Sender == request.Receiver {
		c.JSON(http.StatusBadRequest, dto2.BaseResponse{Message: "不能给自己转账。"})
		return
	}
	err := service.TransferUserPoints(c, request.Sender, request.Receiver, request.Value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto2.BaseResponse{Message: "用户积分转账失败。"})
		return
	}
	c.JSON(http.StatusOK, dto2.BaseResponse{Message: "用户积分转账成功。"})
}

func GetUserPointDetailListHandler(c *gin.Context) {
	var request dto2.UserPointDetailListRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto2.BaseResponse{Message: "参数错误"})
		return
	}
	user := middleware.GetCurrentUser(c)
	filter := model.UserPointDetailFilter{
		UserID:    user.ID,
		StartTime: request.StartTime,
		EndTime:   request.EndTime,
	}
	totalValue, userPointDetails, err := service.GetUserPointDetailList(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto2.BaseResponse{Message: "内部错误。"})
	}

	total, _ := service.GetUserPointDetailCount(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto2.BaseResponse{Message: "内部错误。"})
	}
	response := dto2.UserPointDetailListResponse{
		BasePaginateResponse: dto2.BasePaginateResponse[model.UserPointDetailItem]{
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
	var request dto2.UserPointDetailListAdminRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto2.BaseResponse{Message: "参数错误"})
		return
	}
	filter := model.UserPointDetailFilter{
		UserID:    request.UserID,
		StartTime: request.StartTime,
		EndTime:   request.EndTime,
	}
	_, userPointDetails, err := service.GetUserPointDetailList(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto2.BaseResponse{Message: "内部错误。"})
	}

	total, _ := service.GetUserPointDetailCount(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto2.BaseResponse{Message: "内部错误。"})
	}
	response := dto2.UserPointDetailListResponse{
		BasePaginateResponse: dto2.BasePaginateResponse[model.UserPointDetailItem]{
			Page:     request.Page,
			PageSize: request.PageSize,
			Total:    total,
			Data:     userPointDetails,
		},
	}
	c.JSON(http.StatusOK, response)
}
