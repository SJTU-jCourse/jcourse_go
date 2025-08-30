package types

//go:generate go run ../../cmd/codegen/codegen.go -type=RatingRelatedType,PointEventType,UserRole,UserType,SettingType -file types.go

type RatingRelatedType string

const (
	RelatedTypeCourse       RatingRelatedType = "course"
	RelatedTypeTeacher      RatingRelatedType = "teacher"
	RelatedTypeTrainingPlan RatingRelatedType = "training_plan"
)

type PointEventType string

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

type UserRole string

const (
	UserRoleNormal UserRole = "normal"
	UserRoleAdmin  UserRole = "admin"
)

type UserType string

const (
	UserTypeStudent UserType = "student"
	UserTypeFaculty UserType = "faculty"
)

type SettingType string

const (
	SettingTypeString SettingType = "string"
	SettingTypeInt    SettingType = "int"
	SettingTypeBool   SettingType = "bool"
)
