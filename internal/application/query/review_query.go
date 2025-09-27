package query

import (
	"context"

	"gorm.io/gorm"

	"jcourse_go/internal/application/vo"
	"jcourse_go/internal/domain/shared"
	"jcourse_go/internal/infrastructure/entity"
)

type ReviewQueryService interface {
	GetReview(ctx context.Context, reviewID shared.IDType) (*vo.ReviewVO, error)
	GetLatestReviews(ctx context.Context, query shared.PaginationQuery) ([]vo.ReviewVO, error)
	GetCourseReviews(ctx context.Context, courseID shared.IDType, query shared.PaginationQuery) ([]vo.ReviewVO, error)
	GetUserReviews(ctx context.Context, userID shared.IDType, query shared.PaginationQuery) ([]vo.ReviewVO, error)
}

type reviewQueryService struct {
	db *gorm.DB
}

func (r *reviewQueryService) GetReview(ctx context.Context, reviewID shared.IDType) (*vo.ReviewVO, error) {
	review := entity.Review{}
	if err := r.db.WithContext(ctx).
		Model(&entity.Review{}).
		Joins("Course").
		Preload("Course.MainTeacher").
		Preload("Reactions").
		Where("id = ?", reviewID).
		Take(&review).Error; err != nil {
		return nil, err
	}
	reviewVO := vo.NewReviewVOFromEntity(&review)
	return &reviewVO, nil
}

func (r *reviewQueryService) GetLatestReviews(ctx context.Context, query shared.PaginationQuery) ([]vo.ReviewVO, error) {
	rs := make([]entity.Review, 0)
	if err := r.db.WithContext(ctx).
		Model(&entity.Review{}).
		Joins("Course.MainTeacher").
		Preload("Reactions").
		Order("updated_at desc").
		Find(&rs).Error; err != nil {
		return nil, err
	}
	res := make([]vo.ReviewVO, 0)
	for _, e := range rs {
		res = append(res, vo.NewReviewVOFromEntity(&e))
	}
	return res, nil
}

func (r *reviewQueryService) GetCourseReviews(ctx context.Context, courseID shared.IDType, query shared.PaginationQuery) ([]vo.ReviewVO, error) {
	rs := make([]entity.Review, 0)
	if err := r.db.WithContext(ctx).
		Model(&entity.Review{}).
		Preload("Reactions").
		Where("course_id = ?", courseID).
		Order("updated_at desc").
		Find(&rs).Error; err != nil {
		return nil, err
	}
	res := make([]vo.ReviewVO, 0)
	for _, e := range rs {
		res = append(res, vo.NewReviewVOFromEntity(&e))
	}
	return res, nil
}

func (r *reviewQueryService) GetUserReviews(ctx context.Context, userID shared.IDType, query shared.PaginationQuery) ([]vo.ReviewVO, error) {
	rs := make([]entity.Review, 0)
	if err := r.db.WithContext(ctx).
		Model(&entity.Review{}).
		Joins("Course.MainTeacher").
		Preload("Reactions").
		Where("user_id = ?", userID).
		Order("updated_at desc").
		Find(&rs).Error; err != nil {
		return nil, err
	}
	res := make([]vo.ReviewVO, 0)
	for _, e := range rs {
		res = append(res, vo.NewReviewVOFromEntity(&e))
	}
	return res, nil
}

func NewReviewQueryService(db *gorm.DB) ReviewQueryService {
	return &reviewQueryService{db: db}
}
