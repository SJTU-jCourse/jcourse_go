package repository

import (
	"context"

	"gorm.io/gorm"

	"jcourse_go/internal/domain/course"
	"jcourse_go/internal/domain/shared"
	"jcourse_go/internal/infrastructure/entity"
)

type ReviewRepository struct {
	db *gorm.DB
}

func (r *ReviewRepository) Get(ctx context.Context, id shared.IDType) (*course.Review, error) {
	e := &entity.Review{}
	if err := r.db.Model(&entity.Review{}).Where("id = ?", id).First(&e).Error; err != nil {
		return nil, err
	}
	d := newReviewDomainFromEntity(e)
	return &d, nil
}

func (r *ReviewRepository) Create(ctx context.Context, review *course.Review) error {
	e := newReviewEntityFromDomain(review)
	if err := r.db.Model(&entity.Review{}).Create(&e).Error; err != nil {
		return err
	}
	review.ID = shared.IDType(e.ID)
	return nil
}

func (r *ReviewRepository) Update(ctx context.Context, review *course.Review, revision *course.ReviewRevision) error {
	re := newReviewEntityFromDomain(review)
	rve := newRevisionEntityFromDomain(revision)
	if err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&entity.ReviewRevision{}).Create(rve).Error; err != nil {
			return err
		}
		if err := tx.Model(&entity.Review{}).Updates(re).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	revision.ID = shared.IDType(rve.ID)
	return nil
}

func (r *ReviewRepository) Delete(ctx context.Context, reviewID shared.IDType) error {
	if err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&entity.Review{}, "id = ?", reviewID).Error; err != nil {
			return err
		}
		if err := tx.Delete(&entity.ReviewRevision{}, "review_id = ?", reviewID).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func NewReviewRepository(db *gorm.DB) course.ReviewRepository {
	return &ReviewRepository{db: db}
}

func newReviewDomainFromEntity(r *entity.Review) course.Review {
	return course.Review{
		ID:        shared.IDType(r.ID),
		CourseID:  shared.IDType(r.CourseID),
		UserID:    shared.IDType(r.UserID),
		Score:     r.Score,
		Comment:   r.Comment,
		Rating:    r.Rating,
		IsPublic:  r.IsPublic,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

func newReviewEntityFromDomain(r *course.Review) entity.Review {
	return entity.Review{
		ID:        int64(r.ID),
		CourseID:  int64(r.CourseID),
		UserID:    int64(r.UserID),
		Score:     r.Score,
		Comment:   r.Comment,
		Rating:    r.Rating,
		IsPublic:  r.IsPublic,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

func newRevisionEntityFromDomain(rv *course.ReviewRevision) entity.ReviewRevision {
	return entity.ReviewRevision{
		ID:        int64(rv.ID),
		ReviewID:  int64(rv.ReviewID),
		UserID:    int64(rv.UserID),
		Score:     rv.Score,
		Comment:   rv.Comment,
		Rating:    rv.Rating,
		IsPublic:  rv.IsPublic,
		CreatedAt: rv.CreatedAt,
	}
}
