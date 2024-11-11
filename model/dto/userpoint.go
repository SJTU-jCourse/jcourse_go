package dto

import "jcourse_go/model/model"

type UserListResponse = BasePaginateResponse[model.UserMinimal]

type UserPointDetailRequest struct {
	DetailID int64 `uri:"detailID" binding:"required"`
}
type UserPointDetailResponse = BaseResponse
type UserPointDetailListRequest struct {
	StartTime int64 `json:"start_time" form:"start_time"` // unix timestamp, 单位秒
	EndTime   int64 `json:"end_time" form:"end_time"`
	model.PaginationFilterForQuery
}
type UserPointDetailListAdminRequest struct {
	UserID int64 `json:"user_id" form:"user_id"`
	UserPointDetailListRequest
}

type UserPointDetailListResponse = BasePaginateResponse[model.UserPointDetailItem]

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
