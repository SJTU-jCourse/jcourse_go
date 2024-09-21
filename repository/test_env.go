package repository

import (
	"context"

	"gorm.io/gorm"

	"jcourse_go/model/po"
)

func createTestBaseCourse(db *gorm.DB) error {
	baseCourses := []po.BaseCoursePO{
		{
			Model:  gorm.Model{ID: 1},
			Code:   "MARX1001",
			Name:   "思想道德修养与法律基础",
			Credit: 3,
		},
		{
			Model:  gorm.Model{ID: 2},
			Code:   "CS1500",
			Name:   "计算机科学导论",
			Credit: 4,
		},
		{
			Model:  gorm.Model{ID: 3},
			Code:   "CS2500",
			Name:   "算法与复杂性",
			Credit: 2,
		}}
	err := db.Create(&baseCourses).Error
	return err
}

func CreateTestEnv(ctx context.Context, db *gorm.DB) error {
	db = db.WithContext(ctx)
	err := createTestBaseCourse(db)
	if err != nil {
		return err
	}
	return nil
}
