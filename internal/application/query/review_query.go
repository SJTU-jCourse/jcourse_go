package query

import (
	"context"

	"gorm.io/gorm"

	"jcourse_go/internal/application/vo"
	"jcourse_go/internal/domain/shared"
)

type ReviewQueryService interface {
	GetLatestReviews(ctx context.Context, query shared.PaginationQuery) ([]vo.ReviewVO, error)
	GetCourseReviews(ctx context.Context, courseID shared.IDType, query shared.PaginationQuery) ([]vo.ReviewVO, error)
	GetUserReviews(ctx context.Context, userID shared.IDType, query shared.PaginationQuery) ([]vo.ReviewVO, error)
}

type reviewQueryService struct {
	db *gorm.DB
}

func (r *reviewQueryService) GetLatestReviews(ctx context.Context, query shared.PaginationQuery) ([]vo.ReviewVO, error) {
	// TODO implement me
	panic("implement me")
}

func (r *reviewQueryService) GetCourseReviews(ctx context.Context, courseID shared.IDType, query shared.PaginationQuery) ([]vo.ReviewVO, error) {
	// TODO implement me
	panic("implement me")
}

func (r *reviewQueryService) GetUserReviews(ctx context.Context, userID shared.IDType, query shared.PaginationQuery) ([]vo.ReviewVO, error) {
	// TODO implement me
	panic("implement me")
}

func NewReviewQueryService(db *gorm.DB) ReviewQueryService {
	return &reviewQueryService{db: db}
}
