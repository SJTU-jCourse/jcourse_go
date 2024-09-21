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

func createTestTeacher(db *gorm.DB) error {
	teachers := []po.TeacherPO{
		{
			Model:      gorm.Model{ID: 1},
			Code:       "10001",
			Name:       "高女士",
			Email:      "gaoxiaofeng@example.com",
			Department: "SEIEE",
			Pinyin:     "gaoxiaofeng",
			PinyinAbbr: "gxf",
			Title:      "教授",
		},
		{
			Model:      gorm.Model{ID: 2},
			Code:       "10002",
			Name:       "潘老师",
			Email:      "panli@example.com",
			Department: "SEIEE",
			Pinyin:     "panli",
			PinyinAbbr: "pl",
			Title:      "教授",
		},
		{
			Model:      gorm.Model{ID: 3},
			Code:       "10003",
			Name:       "梁女士",
			Email:      "liangqin@example.com",
			Department: "PHYSICS",
			Pinyin:     "liangqin",
			PinyinAbbr: "lq",
			Title:      "教授",
		},
		{
			Model:      gorm.Model{ID: 4},
			Code:       "10004",
			Name:       "赵先生",
			Email:      "zhaohao@example.com",
			Pinyin:     "zhaohao",
			PinyinAbbr: "zh",
			Title:      "讲师",
		},
	}
	err := db.Create(&teachers).Error
	return err
}

func createTestCourse(db *gorm.DB) error {
	courses := []po.CoursePO{
		{
			Model:           gorm.Model{ID: 1},
			Code:            "MARX1001",
			Name:            "思想道德修养与法律基础",
			Credit:          3,
			MainTeacherID:   3,
			MainTeacherName: "梁女士",
			Department:      "MARX",
		},
		{
			Model:           gorm.Model{ID: 2},
			Code:            "MARX1001",
			Name:            "思想道德修养与法律基础",
			Credit:          3,
			MainTeacherID:   4,
			MainTeacherName: "赵先生",
			Department:      "MARX",
		},
		{
			Model:           gorm.Model{ID: 3},
			Code:            "CS1500",
			Name:            "计算机科学导论",
			Credit:          3,
			MainTeacherID:   1,
			MainTeacherName: "高女士",
			Department:      "SEIEE",
		},
		{
			Model:           gorm.Model{ID: 4},
			Code:            "CS2500",
			Name:            "算法与复杂性",
			Credit:          3,
			MainTeacherID:   1,
			MainTeacherName: "高女士",
			Department:      "SEIEE",
		},
	}
	err := db.Create(&courses).Error
	return err
}

func createCourseCategories(db *gorm.DB) error {
	categories := []po.CourseCategoryPO{
		{
			CourseID: 1,
			Category: "通识",
		},
		{
			CourseID: 2,
			Category: "通识",
		},
		{
			CourseID: 2,
			Category: "必修",
		},
	}
	err := db.Create(&categories).Error
	return err
}

func CreateTestEnv(ctx context.Context, db *gorm.DB) error {
	db = db.WithContext(ctx)
	createFunc := []func(db *gorm.DB) error{
		createTestBaseCourse,
		createTestTeacher,
		createTestCourse,
		createCourseCategories,
	}
	for _, fn := range createFunc {
		err := fn(db)
		if err != nil {
			return err
		}
	}
	return nil
}
