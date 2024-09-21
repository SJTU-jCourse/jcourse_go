package repository

import (
	"gorm.io/gorm"

	"jcourse_go/model/po"
)

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&po.UserPO{},
		&po.BaseCoursePO{}, &po.CoursePO{}, &po.TeacherPO{}, &po.CourseCategoryPO{},
		&po.OfferedCoursePO{}, &po.OfferedCourseTeacherPO{},
		&po.ReviewPO{}, &po.RatingPO{}, &po.TrainingPlanPO{}, &po.TrainingPlanCoursePO{})
	if err != nil {
		return err
	}
	return nil
}
