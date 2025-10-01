package user

//go:generate go run github.com/dmarkham/enumer -type=PointEvent -transform=snake -trimprefix=PointEvent
type PointEvent int

const (
	PointEventReview PointEvent = iota
	PointEventLike
	PointEventBeLiked
	PointEventAdminChange
	PointEventInit
	PointEventTransferIn
	PointEventTransferOut
	PointEventReward
	PointEventPunish
	PointEventWithdraw
	PointEventConsume
	PointEventRedeem
)
