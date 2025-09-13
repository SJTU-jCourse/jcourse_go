package vo

import "jcourse_go/internal/infrastructure/entity"

type UserPointVO struct {
	ID          int64
	Type        string
	Description string
	Value       int64
	CreatedAt   int64
}

func NewUserPointFromEntity(e *entity.UserPoint) UserPointVO {
	return UserPointVO{
		ID:          e.ID,
		Type:        e.Type,
		Description: e.Description,
		Value:       e.Value,
		CreatedAt:   e.CreatedAt.Unix(),
	}
}
