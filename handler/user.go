package handler

import (
	"github.com/gin-gonic/gin"
	"jcourse_go/constant"
	"jcourse_go/model/converter"
	"jcourse_go/model/domain"
	"jcourse_go/model/dto"
	"jcourse_go/service"
	"net/http"
	"strconv"
)

func GetSuggestedUserHandler(c *gin.Context) {}

func GetUserListHandler(c *gin.Context) {
	/*
		管理员权限验证
	*/
	var request dto.UserListRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{Message: "参数错误"})
		return
	}

	filter := domain.UserFilter{
		Page:     request.Page,
		PageSize: request.PageSize,
	}
	users, err := service.GetUserList(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{Message: "内部错误。"})
	}
	total, _ := service.GetUserCount(c, filter)
	response := dto.UserListResponse{
		Page:     request.Page,
		PageSize: request.PageSize,
		Total:    total,
		Data:     users,
	}
	c.JSON(http.StatusOK, response)
}

func GetCurrentUserHandler(c *gin.Context) {
	userInterface, exists := c.Get(constant.CtxKeyUser)
	if !exists {
		c.JSON(http.StatusNotFound, dto.BaseResponse{Message: "用户未登录！"})
		return
	}
	user, _ := userInterface.(*domain.User)

	me, err := service.GetUserDetailsByID(c, user.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.BaseResponse{Message: "此用户不存在！"})
		return
	}
	c.JSON(http.StatusOK, me)
	return
}

func GetUserDetailHandler(c *gin.Context) {
	userIDStr := c.Param("userID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "非法用户ID"})
		return
	}

	userDetail, errDetail := service.GetUserDetailsByID(c, int64(userID))
	if errDetail != nil {
		c.JSON(http.StatusNotFound, dto.BaseResponse{Message: "此用户不存在！"})
	}
	c.JSON(http.StatusOK, userDetail)
}

func WatchUserHandler(c *gin.Context) {}

func UnWatchUserHandler(c *gin.Context) {}

func UpdateUserInfoHandler(c *gin.Context) {}

func GetUserReviewsHandler(c *gin.Context) {
	userIDStr := c.Param("userID")

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "非法用户ID"})
		return
	}

	var request dto.ReviewListRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{Message: "参数错误"})
		return
	}

	filter := domain.ReviewFilter{
		Page:     request.Page,
		PageSize: request.PageSize,
		UserID:   int64(userID),
	}

	reviews, err := service.GetReviewList(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{Message: "内部错误。"})
		return
	}

	total, err := service.GetReviewCount(c, filter)

	response := dto.ReviewListResponse{
		Page:     request.Page,
		PageSize: request.PageSize,
		Total:    total,
		Data:     converter.ConvertReviewDomainToListDTO(reviews, true),
	}
	c.JSON(http.StatusOK, response)
}
