package repository

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"jcourse_go/model/converter"
	"jcourse_go/model/model"
	"jcourse_go/model/po"
)

type IReviewQuery interface {
	GetReviewCount(ctx context.Context, opts ...DBOption) (int64, error)
	GetReview(ctx context.Context, opts ...DBOption) ([]po.ReviewPO, error)
	CreateReview(ctx context.Context, review po.ReviewPO) (int64, error)
	UpdateReview(ctx context.Context, review po.ReviewPO) error
	DeleteReview(ctx context.Context, opts ...DBOption) error
	// SyncCourseRating(ctx context.Context, courseID int64) error
}

type ReviewQuery struct {
	db *gorm.DB
}

func (c *ReviewQuery) SyncCourseRating(ctx context.Context, db *gorm.DB, courseID int64) error {
	// 1. Aggregate review data
	var info po.CourseReviewInfo
	if err := db.WithContext(ctx).
		Model(&po.RatingPO{}).
		Select("COUNT(*) AS count, AVG(rating) AS average").
		Where("related_type = ? and related_id = ?", model.RelatedTypeCourse, courseID).
		Scan(&info).
		Error; err != nil {
		return err
	}

	// 2. Update the matching course row
	if err := db.WithContext(ctx).
		Model(&po.CoursePO{}).
		Where("id = ?", courseID).
		Updates(map[string]interface{}{
			"rating_count": info.Count,
			"rating_avg":   info.Average,
		}).Error; err != nil {
		return err
	}
	return nil
}

func (c *ReviewQuery) GetReviewCount(ctx context.Context, opts ...DBOption) (int64, error) {
	db := c.optionDB(ctx, opts...)
	count := int64(0)
	result := db.Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

func (c *ReviewQuery) optionDB(ctx context.Context, opts ...DBOption) *gorm.DB {
	db := c.db.WithContext(ctx).Model(&po.ReviewPO{})
	for _, opt := range opts {
		db = opt(db)
	}
	return db
}

func (c *ReviewQuery) GetReview(ctx context.Context, opts ...DBOption) ([]po.ReviewPO, error) {
	db := c.optionDB(ctx, opts...)
	reviews := make([]po.ReviewPO, 0)
	result := db.WithContext(ctx).Find(&reviews)
	if result.Error != nil {
		return nil, result.Error
	}
	return reviews, nil
}

func (c *ReviewQuery) CreateReview(ctx context.Context, review po.ReviewPO) (int64, error) {
	db := c.db.WithContext(ctx)
	ratingPO := converter.BuildRatingFromReview(review)
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&po.ReviewPO{}).Create(&review).Error; err != nil {
			return err
		}
		if err := tx.Model(&po.RatingPO{}).Create(&ratingPO).Error; err != nil {
			return err
		}
		if err := c.SyncCourseRating(ctx, tx, review.CourseID); err != nil {
			return err
		}
		return nil
	})
	return int64(review.ID), err
}

func (c *ReviewQuery) UpdateReview(ctx context.Context, review po.ReviewPO) error {
	db := c.db.WithContext(ctx)
	ratingPO := converter.BuildRatingFromReview(review)
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&po.ReviewPO{}).Where("id = ?", review.ID).Updates(&review).Error; err != nil {
			return err
		}
		if err := tx.Model(&po.RatingPO{}).
			Where("related_type = ? and related_id = ?", model.RelatedTypeCourse, review.ID).
			Updates(&ratingPO).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func (c *ReviewQuery) DeleteReview(ctx context.Context, opts ...DBOption) error {
	err := c.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		tx1 := tx.Session(&gorm.Session{})
		for _, opt := range opts {
			tx1 = opt(tx1)
		}
		reviews := make([]po.ReviewPO, 0)
		if err := tx1.Model(&po.ReviewPO{}).Clauses(clause.Returning{}).Delete(&reviews).Error; err != nil {
			return err
		}

		ids := make([]int64, 0)
		for _, review := range reviews {
			ids = append(ids, int64(review.ID))
		}
		tx2 := tx.Session(&gorm.Session{})
		if err := tx2.Model(&po.RatingPO{}).Delete(&po.RatingPO{}, "related_id in ? and related_type = ?", ids, model.RelatedTypeCourse).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func NewReviewQuery(db *gorm.DB) IReviewQuery {
	return &ReviewQuery{db: db}
}
