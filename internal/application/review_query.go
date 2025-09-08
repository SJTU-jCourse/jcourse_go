package application

import (
	"context"

	"jcourse_go/internal/application/vo"
	"jcourse_go/internal/domain/shared"
)

type ReviewQueryService interface {
	GetLatestReviews(ctx context.Context, query shared.PaginationQuery) ([]vo.ReviewVO, error)
	GetCourseReviews(ctx context.Context, courseID shared.IDType, query shared.PaginationQuery) ([]vo.ReviewVO, error)
	GetUserReviews(ctx context.Context, userID shared.IDType, query shared.PaginationQuery) ([]vo.ReviewVO, error)
}
