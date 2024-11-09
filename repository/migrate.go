package repository

import (
	"gorm.io/gorm"

	"jcourse_go/model/po"
)

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&po.UserPO{},
		&po.BaseCoursePO{}, &po.CoursePO{}, &po.TeacherPO{}, &po.CourseCategoryPO{},
		&po.OfferedCoursePO{}, &po.OfferedCourseTeacherPO{},
		&po.TrainingPlanPO{}, &po.TrainingPlanCoursePO{},
		&po.ReviewPO{}, &po.RatingPO{},
		&po.SettingPO{}, &po.UserPointDetailPO{})
	if err != nil {
		return err
	}
	return nil
}
