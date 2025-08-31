package event

type Type string

const (
	TypeUserCreated   Type = "user.created"
	TypeReviewCreated Type = "review.created"
)

type Event struct {
	ID      string
	Type    Type
	Payload any
}
