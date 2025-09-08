package application

import (
	"context"

	"jcourse_go/internal/application/vo"
	"jcourse_go/internal/domain/shared"
)

type UserQueryService interface {
	GetUserInfo(ctx context.Context, userID shared.IDType) (*vo.UserInfoVO, error)
}
