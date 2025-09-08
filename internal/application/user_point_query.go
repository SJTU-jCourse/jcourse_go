package application

import (
	"context"

	"jcourse_go/internal/application/vo"
	"jcourse_go/internal/domain/shared"
)

type UserPointQueryService interface {
	GetUserPoint(ctx context.Context, userID shared.IDType) (int, []vo.UserPointVO, error)
}
