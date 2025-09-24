package user

type PointEventType string

const (
	PointEventReview      PointEventType = "review"
	PointEventLike        PointEventType = "like"
	PointEventBeLiked     PointEventType = "be_liked"
	PointEventAdminChange PointEventType = "admin_change"
	PointEventInit        PointEventType = "init"
	PointEventTransferIn  PointEventType = "transfer_in"
	PointEventTransferOut PointEventType = "transfer_out"
	PointEventReward      PointEventType = "reward"
	PointEventPunish      PointEventType = "punish"
	PointEventWithdraw    PointEventType = "withdraw"
	PointEventConsume     PointEventType = "consume"
	PointEventRedeem      PointEventType = "redeem" // 兑换积分
)
