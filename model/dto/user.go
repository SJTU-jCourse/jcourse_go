package dto

import "jcourse_go/model/model"

type UserProfileDTO struct {
	Username   string `json:"username"`
	Avatar     string `json:"avatar"`
	Bio        string `json:"bio"`
	Type       string `json:"type"`
	Department string `json:"department"`
	Major      string `json:"major"`
	Grade      string `json:"grade"`
	Degree     string `json:"degree"`
}

type UserListRequest struct {
	model.PaginationFilterForQuery
}

type UserListResponse = BasePaginateResponse[model.UserMinimal]

type UserPointDetailRequestURI struct {
	DetailID int64 `uri:"detailID" binding:"required"`
}
type UserPointDetailRequestJSON struct {
	UserID int64 `json:"user_id" form:"user_id" binding:"required"` // use for auth
}
type UserPointDetailRequest struct {
	UserID   int64
	DetailID int64
}

type UserPointDetailResponse = BaseResponse
type UserPointDetailListRequest struct {
	Page      int64 `json:"page" form:"page" binding:"required"`
	PageSize  int64 `json:"page_size" form:"page_size" binding:"required"`
	UserID    int64 `json:"user_id" form:"user_id" binding:"required"`
	StartTime int64 `json:"start_time" form:"start_time"` // unix timestamp, 单位秒
	EndTime   int64 `json:"end_time" form:"end_time"`
}

type UserPointDetailListResponse = BasePaginateResponse[model.UserPointDetailItem]

type ChangeUserPointRequest struct {
	UserID int64 `json:"user_id" form:"user_id" binding:"required"`
	Value  int64 `json:"value" form:"value" binding:"required"`
}
type ChangeUserPointResponse = BaseResponse

type TransferUserPointRequest struct {
	Sender   int64 `json:"sender" form:"sender" binding:"required"`
	Receiver int64 `json:"receiver" form:"receiver" binding:"required"`
	Value    int64 `json:"value" form:"value" binding:"required"`
}
type TransferUserPointResponse = BaseResponse

type RedeemUserPointRequest struct {
	UserID int64 `json:"user_id" form:"user_id" binding:"required"`
	Value  int64 `json:"value" form:"value" binding:"required"`
}
type RedeemUserPointResponse = BaseResponse
