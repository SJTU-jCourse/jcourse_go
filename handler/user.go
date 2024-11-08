package handler

import (
	"errors"
	"jcourse_go/constant"
	"jcourse_go/middleware"
	"jcourse_go/model/converter"
	"jcourse_go/model/dto"
	"jcourse_go/model/model"
	"jcourse_go/service"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetSuggestedUserHandler(c *gin.Context) {}

func GetUserListHandler(c *gin.Context) {
	var request dto.UserListRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{Message: "参数错误"})
		return
	}

	filter := model.UserFilterForQuery{
		PaginationFilterForQuery: request.PaginationFilterForQuery,
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

func getUserIDFromRequest(c *gin.Context) (int64, error) {
	userIDStr := c.Param("userID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return -1, errors.New("非法用户ID")
	}
	return int64(userID), nil
}

// 非公开信息？
func GetUserActivityHandler(c *gin.Context) {
	userID, err := getUserIDFromRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{Message: "非法用户ID"})
		return
	}

	user := middleware.GetCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, dto.BaseResponse{Message: "用户未登录！"})
		return
	}

	if user.ID != userID {
		c.JSON(http.StatusForbidden, dto.BaseResponse{Message: "无权查看他人信息！"})
		return
	}

	userSummary, err := service.GetUserActivityByID(c, user.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.BaseResponse{Message: "此用户不存在！"})
		return
	}
	c.JSON(http.StatusOK, userSummary)
}

// 公开信息
func GetUserDetailHandler(c *gin.Context) {
	userID, err := getUserIDFromRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{Message: "非法用户ID"})
		return
	}

	currentUserID := int64(0)
	currentUser := middleware.GetCurrentUser(c)
	if currentUser != nil {
		currentUserID = currentUser.ID
	}

	user, err := service.GetUserDetailByID(c, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.BaseResponse{Message: "此用户不存在！"})
		return
	}
	converter.RemoveUserEmail(user, currentUserID)
	log.Printf("user: %v", user)
	c.JSON(http.StatusOK, user)
}

func WatchUserHandler(c *gin.Context) {}

func UnWatchUserHandler(c *gin.Context) {}

func UpdateUserProfileHandler(c *gin.Context) {
	userInterface, exists := c.Get(constant.CtxKeyUser)
	if !exists {
		c.JSON(http.StatusNotFound, dto.BaseResponse{Message: "用户未登录！"})
		return
	}
	user, _ := userInterface.(*model.UserDetail)

	var request dto.UserProfileDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{Message: "参数错误"})
		return
	}

	err := service.UpdateUserProfileByID(c, request, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{Message: "用户信息更新失败。"})
		return
	}
	c.JSON(http.StatusOK, dto.BaseResponse{Message: "用户信息更新成功。"})
}

func TransferUserPointHandler(c *gin.Context) {
	var request dto.TransferUserPointRequest
	if err := c.ShouldBind(&request); err != nil {
		log.Printf("TransferUserPointHandler: %v, %v", request, c.Request.URL.Query())
		c.JSON(http.StatusBadRequest, dto.BaseResponse{Message: "参数错误"})
		return
	}
	if request.Sender == request.Receiver {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{Message: "不能给自己转账。"})
		return
	}
	err := service.TransferUserPoints(c, request.Sender, request.Receiver, request.Value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{Message: "用户积分转账失败。"})
		log.Printf("TransferUserPointHandler: %v", err)
		return
	}
	c.JSON(http.StatusOK, dto.BaseResponse{Message: "用户积分转账成功。"})
}

func GetUserPointDetailHandler(c *gin.Context) {
	var requestUri dto.UserPointDetailRequestURI
	if err := c.ShouldBindUri(&requestUri); err != nil {
		c.JSON(http.StatusNotFound, dto.BaseResponse{Message: "参数错误"})
		return
	}
	var requestJson dto.UserPointDetailRequestJSON
	if err := c.ShouldBind(&requestJson); err != nil {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{Message: "参数错误"})
		return
	}
	var request = dto.UserPointDetailRequest{
		UserID:   requestJson.UserID,
		DetailID: requestUri.DetailID,
	}
	log.Printf("userID: %v", request.UserID)
	if !middleware.IsMineOrAdmin(c, request.UserID) {
		c.JSON(http.StatusForbidden, dto.BaseResponse{Message: "无权查看他人积分。"})
		return
	}
	userPointDetails, err := service.GetUserPointDetailList(c, model.UserPointDetailFilter{UserPointDetailID: request.DetailID})
	if err != nil || len(userPointDetails) == 0 {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{Message: "内部错误。"})
		return
	}
	c.JSON(http.StatusOK, userPointDetails[0])
}

func GetUserPointDetailListHandler(c *gin.Context) {
	var request = dto.UserPointDetailListRequest{
		Page:     constant.DefaultPage,
		PageSize: constant.DefaultPageSize,
	}
	if err := c.ShouldBind(&request); err != nil {
		log.Printf("GetUserPointDetailListHandler: %v, %v", request, c.Request.URL.Query())
		c.JSON(http.StatusBadRequest, dto.BaseResponse{Message: "参数错误"})
		return
	}
	if !middleware.IsMineOrAdmin(c, request.UserID) {
		c.JSON(http.StatusForbidden, dto.BaseResponse{Message: "无权查看他人积分。"})
		return
	}
	filter := model.UserPointDetailFilter{
		UserID:    request.UserID,
		StartTime: time.Unix(request.StartTime, 0),
		EndTime:   time.Unix(request.EndTime, 0),
	}
	userPointDetails, err := service.GetUserPointDetailList(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{Message: "内部错误。"})
	}

	total, _ := service.GetUserPointDetailCount(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{Message: "内部错误。"})
	}
	response := dto.UserPointDetailListResponse{
		Page:     request.Page,
		PageSize: request.PageSize,
		Total:    total,
		Data:     userPointDetails,
	}
	c.JSON(http.StatusOK, response)
}
