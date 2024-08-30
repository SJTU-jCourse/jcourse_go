package repository

import (
	"context"
	"jcourse_go/dal"
	"jcourse_go/model/po"

	"gorm.io/gorm"
)

type ICourseVectorQuery interface {
	InsertCourseVector(ctx context.Context, courseID int64, vector []float32) error
	UpdateCourseVector(ctx context.Context, courseID int64, vector []float32) error
}
type CourseVectorQuery struct {
	db *gorm.DB
}

func (b *CourseVectorQuery) optionDB(ctx context.Context, opts ...DBOption) *gorm.DB {
	db := b.db.WithContext(ctx)
	for _, opt := range opts {
		db = opt(db)
	}
	return db
}

func (c *CourseVectorQuery) InsertCourseVector(ctx context.Context, courseID int64, vector []float32) error {
	db := c.optionDB(ctx)
	courseVectorPO := po.CourseVectorPO{
		ID:     courseID,
		Vector: vector,
	}
	result := db.Create(&courseVectorPO)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (c *CourseVectorQuery) UpdateCourseVector(ctx context.Context, courseID int64, vector []float32) error {
	db := c.optionDB(ctx)
	courseVectorPO := po.CourseVectorPO{
		ID:     courseID,
		Vector: vector,
	}
	result := db.Model(&po.CourseVectorPO{}).Where("id = ?", courseID).Updates(&courseVectorPO)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func NewCourseVectorQuery() ICourseVectorQuery {
	return &CourseVectorQuery{db: dal.GetDBClient()}
}
