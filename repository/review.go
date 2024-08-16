package repository

import (
	"context"

	"gorm.io/gorm"

	"jcourse_go/dal"
	"jcourse_go/model/po"
)

type IReviewQuery interface {
	GetReviewCount(ctx context.Context, opts ...DBOption) (int64, error)
	GetReviewDetail(ctx context.Context, opts ...DBOption) (*po.ReviewPO, error)
	GetReviewList(ctx context.Context, opts ...DBOption) ([]po.ReviewPO, error)
	CreateReview(ctx context.Context, review po.ReviewPO) (int64, error)
	UpdateReview(ctx context.Context, review po.ReviewPO) (int64, error)
	DeleteReview(ctx context.Context, opts ...DBOption) (int64, error)
	GetCourseReviewInfo(ctx context.Context, courseIDs []int64) (map[int64]po.CourseReviewInfo, error)
	WithID(id int64) DBOption
	WithCourseID(courseID int64) DBOption
	WithUserID(userID int64) DBOption
	WithSemester(semester string) DBOption
	WithOrderBy(orderBy string, ascending bool) DBOption
	WithLimit(limit int64) DBOption
	WithOffset(offset int64) DBOption
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

func (c *ReviewQuery) WithLimit(limit int64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(int(limit))
	}
}

func (c *ReviewQuery) WithOffset(offset int64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(int(offset))
	}
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
	db := c.optionDB(ctx)
	result := db.Create(&review)
	return int64(review.ID), result.Error
}

func (c *ReviewQuery) UpdateReview(ctx context.Context, review po.ReviewPO) (int64, error) {
	result := c.db.WithContext(ctx).Save(&review)
	return result.RowsAffected, result.Error
}

func (c *ReviewQuery) DeleteReview(ctx context.Context, opts ...DBOption) (int64, error) {
	db := c.optionDB(ctx, opts...)
	result := db.Delete(&po.ReviewPO{})
	return result.RowsAffected, result.Error
}

func (c *ReviewQuery) WithID(id int64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func (c *ReviewQuery) WithCourseID(courseID int64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("course_id = ?", courseID)
	}
}

func (c *ReviewQuery) WithUserID(userID int64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", userID)
	}
}

func (c *ReviewQuery) WithSemester(semester string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("semester = ?", semester)
	}
}

func (c *ReviewQuery) WithOrderBy(orderBy string, ascending bool) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		if ascending {
			orderBy = orderBy + " asc"
		} else {
			orderBy = orderBy + " desc"
		}
		return db.Order(orderBy)
	}
}

func NewReviewQuery() IReviewQuery {
	return &ReviewQuery{db: dal.GetDBClient()}
}
