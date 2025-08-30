package model

// 用户积分明细
type UserPointDetailItem struct {
	Time        string `json:"time"`
	Value       int64  `json:"value"` // 积分变化值: +1, -3
	Description string `json:"description"`
}

type UserPointDetailFilter struct {
	PaginationFilterForQuery
	UserPointDetailID int64
	UserID            int64
	StartTime         int64
	EndTime           int64
}
