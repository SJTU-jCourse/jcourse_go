package repository

import (
	"context"
	"testing"

	"gorm.io/gorm"

	"jcourse_go/dal"
	"jcourse_go/model/po"
	"jcourse_go/model/types"
)

func createTestBaseCourse(db *gorm.DB) error {
	baseCourses := []po.BaseCoursePO{
		{
			ID:     1,
			Code:   "MARX1001",
			Name:   "思想道德修养与法律基础",
			Credit: 3,
		},
		{
			ID:     2,
			Code:   "CS1500",
			Name:   "计算机科学导论",
			Credit: 4,
		},
		{
			ID:     3,
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
			ID:         1,
			Code:       "10001",
			Name:       "高女士",
			Email:      "gaoxiaofeng@example.com",
			Department: "SEIEE",
			Pinyin:     "gaoxiaofeng",
			PinyinAbbr: "gxf",
			Title:      "教授",
		},
		{
			ID:         2,
			Code:       "10002",
			Name:       "潘老师",
			Email:      "panli@example.com",
			Department: "SEIEE",
			Pinyin:     "panli",
			PinyinAbbr: "pl",
			Title:      "教授",
		},
		{
			ID:         3,
			Code:       "10003",
			Name:       "梁女士",
			Email:      "liangqin@example.com",
			Department: "PHYSICS",
			Pinyin:     "liangqin",
			PinyinAbbr: "lq",
			Title:      "教授",
		},
		{
			ID:         4,
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
			ID:              1,
			Code:            "MARX1001",
			Name:            "思想道德修养与法律基础",
			Credit:          3,
			MainTeacherID:   3,
			MainTeacherName: "梁女士",
			Department:      "MARX",
		},
		{
			ID:              2,
			Code:            "MARX1001",
			Name:            "思想道德修养与法律基础",
			Credit:          3,
			MainTeacherID:   4,
			MainTeacherName: "赵先生",
			Department:      "MARX",
		},
		{
			ID:              3,
			Code:            "CS1500",
			Name:            "计算机科学导论",
			Credit:          3,
			MainTeacherID:   1,
			MainTeacherName: "高女士",
			Department:      "SEIEE",
		},
		{
			ID:              4,
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

func createTrainingPlan(db *gorm.DB) error {
	trainingPlans := []po.TrainingPlanPO{
		{
			ID:         1,
			Degree:     "本科",
			Department: "SEIEE",
			EntryYear:  "2020",
			Major:      "计算机科学与技术",
			MajorClass: "工学",
		},
	}
	err := db.Create(&trainingPlans).Error
	if err != nil {
		return err
	}

	trainingPlanCourses := []po.TrainingPlanCoursePO{
		{
			TrainingPlanID: 1,
			BaseCourseID:   1,
		},
		{
			TrainingPlanID: 1,
			BaseCourseID:   2,
		},
		{
			TrainingPlanID: 1,
			BaseCourseID:   3,
		},
	}
	err = db.Create(&trainingPlanCourses).Error
	if err != nil {
		return err
	}
	return nil
}

func createTestUser(db *gorm.DB) error {
	users := []po.UserPO{
		{
			ID:       1,
			Username: "test1",
			Email:    "test1@example.com",
		},
		{
			ID:       2,
			Username: "test2",
			Email:    "test2@example.com",
		},
		{
			ID:       3,
			Username: "test3",
			Email:    "test3@example.com",
		},
	}
	err := db.Create(&users).Error
	return err
}

func createTestReview(db *gorm.DB) error {
	reviews := []po.ReviewPO{
		{
			ID:       1,
			CourseID: 1,
			UserID:   1,
			Comment:  "test review",
			Rating:   5,
		},
		{
			ID:       2,
			CourseID: 2,
			UserID:   1,
			Comment:  "test review",
			Rating:   4,
		},
	}
	err := db.Create(&reviews).Error
	if err != nil {
		return err
	}
	ratings := []po.RatingPO{
		{
			ID:          1,
			Rating:      5,
			RelatedID:   1,
			UserID:      1,
			RelatedType: string(types.RelatedTypeCourse),
		},
		{
			ID:          2,
			Rating:      5,
			RelatedID:   2,
			UserID:      1,
			RelatedType: string(types.RelatedTypeCourse),
		},
		{
			ID:          3,
			Rating:      5,
			RelatedID:   2,
			UserID:      1,
			RelatedType: string(types.RelatedTypeTeacher),
		},
	}

	err = db.Create(&ratings).Error
	return err
}

func CreateTestEnv(ctx context.Context, db *gorm.DB) error {
	db = db.WithContext(ctx)
	createFunc := []func(db *gorm.DB) error{
		createTestBaseCourse,
		createTestTeacher,
		createTestCourse,
		createCourseCategories,
		createTrainingPlan,
		createTestUser,
		createTestReview,
	}
	for _, fn := range createFunc {
		err := fn(db)
		if err != nil {
			panic(err)
		}
	}
	return nil
}

func setup() {
	ctx := context.Background()
	dal.InitTestMemDBClient()
	db := dal.GetDBClient()
	_ = Migrate(db)
	_ = CreateTestEnv(ctx, db)
}

func tearDown() {
	db, _ := dal.GetDBClient().DB()
	_ = db.Close()
}

func TestMain(m *testing.M) {
	setup()
	m.Run()
	tearDown()
}
