package repository

import (
	"context"

	"gorm.io/gorm"

	"jcourse_go/model/converter"
	"jcourse_go/model/model"
	"jcourse_go/model/po"
)

type IReviewQuery interface {
	GetReviewCount(ctx context.Context, opts ...DBOption) (int64, error)
	GetReviewDetail(ctx context.Context, opts ...DBOption) (*po.ReviewPO, error)
	GetReviewList(ctx context.Context, opts ...DBOption) ([]po.ReviewPO, error)
	CreateReview(ctx context.Context, review po.ReviewPO) (int64, error)
	UpdateReview(ctx context.Context, review po.ReviewPO) error
	DeleteReview(ctx context.Context, opts ...DBOption) error
	GetCourseReviewInfo(ctx context.Context, courseIDs []int64) (map[int64]po.CourseReviewInfo, error)
}

type ReviewQuery struct {
	db *gorm.DB
}

func (c *ReviewQuery) GetCourseReviewInfo(ctx context.Context, courseIDs []int64) (map[int64]po.CourseReviewInfo, error) {
	infoMap := make(map[int64]po.CourseReviewInfo)
	infos := make([]po.CourseReviewInfo, 0)
	result := c.db.WithContext(ctx).Model(&po.ReviewPO{}).
		Select("count(*) as count, avg(rating) as average, course_id").
		Group("course_id").
		Where("course_id in (?)", courseIDs).
		Find(&infos)
	if result.Error != nil {
		return infoMap, result.Error
	}
	for _, info := range infos {
		infoMap[info.CourseID] = info
	}
	return infoMap, nil
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

func (c *ReviewQuery) GetReviewDetail(ctx context.Context, opts ...DBOption) (*po.ReviewPO, error) {
	db := c.optionDB(ctx, opts...)
	review := po.ReviewPO{}
	result := db.WithContext(ctx).First(&review)
	if result.Error != nil {
		return nil, result.Error
	}
	return &review, nil
}

func (c *ReviewQuery) GetReviewList(ctx context.Context, opts ...DBOption) ([]po.ReviewPO, error) {
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
		return nil
	})
	return int64(review.ID), err
}

func (c *ReviewQuery) UpdateReview(ctx context.Context, review po.ReviewPO) error {
	db := c.db.WithContext(ctx)
	ratingPO := converter.BuildRatingFromReview(review)
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&po.ReviewPO{}).Updates(&review).Error; err != nil {
			return err
		}
		if err := tx.Model(&po.RatingPO{}).Updates(&ratingPO).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func (c *ReviewQuery) DeleteReview(ctx context.Context, opts ...DBOption) error {
	db := c.optionDB(ctx, opts...)
	reviews := make([]po.ReviewPO, 0)
	err := db.Find(&reviews).Error
	if err != nil {
		return err
	}
	ids := make([]int64, 0)
	for _, review := range reviews {
		ids = append(ids, int64(review.ID))
	}
	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&po.ReviewPO{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(&po.RatingPO{}, "related_id in ? and related_type = ?", ids, model.RelatedTypeCourse).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func NewReviewQuery(db *gorm.DB) IReviewQuery {
	return &ReviewQuery{db: db}
}
