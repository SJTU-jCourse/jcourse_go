package dto

import (
	model2 "jcourse_go/internal/model/model"
)

type UserPointDetailListRequest struct {
	StartTime int64 `json:"start_time" form:"start_time"` // unix timestamp, 单位秒
	EndTime   int64 `json:"end_time" form:"end_time"`
	model2.PaginationFilterForQuery
}
type UserPointDetailListAdminRequest struct {
	UserID int64 `json:"user_id" form:"user_id"`
	UserPointDetailListRequest
}

type UserPointDetailListResponse struct {
	CurrentPoint int64 `json:"current_point"`
	BasePaginateResponse[model2.UserPointDetailItem]
}

type ChangeUserPointRequest struct {
	UserID int64 `json:"user_id" form:"user_id" binding:"required"`
	Value  int64 `json:"value" form:"value" binding:"required"`
}
type ChangeUserPointResponse = BaseResponse

type TransferUserPointRequest struct {
	// sender is the current user
	Receiver int64 `json:"receiver" form:"receiver" binding:"required"`
	Value    int64 `json:"value" form:"value" binding:"required"`
}

type TransferUserPointAdminRequest struct {
	Sender int64 `json:"sender" form:"sender" binding:"required"`
	TransferUserPointRequest
}
