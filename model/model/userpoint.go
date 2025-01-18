package model

type PointEventType = string

const (
	PointEventReview      PointEventType = "review"
	PointEventLike        PointEventType = "like"
	PointEventBeLiked     PointEventType = "be_liked"
	PointEventAdminChange PointEventType = "admin_change"
	PointEventInit        PointEventType = "init"
	PointEventTransfer    PointEventType = "transfer"
	PointEventReward      PointEventType = "reward"
	PointEventPunish      PointEventType = "punish"
	PointEventWithdraw    PointEventType = "withdraw"
	PointEventConsume     PointEventType = "consume"
	PointEventRedeem      PointEventType = "redeem" // 兑换积分
)

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
