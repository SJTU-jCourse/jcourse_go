package domain

type RewardType = string

const (
	RewardTypePoint RewardType = "point"
)

type RewardRecord struct {
	ID           uint
	UserID       int64
	RewardType   RewardType
	RewardAmount int64
}
