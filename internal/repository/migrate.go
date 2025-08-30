package repository

import (
	"gorm.io/gorm"

	po2 "jcourse_go/internal/model/po"
)

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&po2.UserPO{},
		&po2.BaseCoursePO{}, &po2.CoursePO{}, &po2.TeacherPO{}, &po2.CourseCategoryPO{},
		&po2.OfferedCoursePO{}, &po2.OfferedCourseTeacherPO{},
		&po2.TrainingPlanPO{}, &po2.TrainingPlanCoursePO{},
		&po2.ReviewPO{}, &po2.RatingPO{}, &po2.ReviewRevisionPO{}, &po2.ReviewReactionPO{},
		&po2.SettingPO{}, &po2.UserPointDetailPO{}, &po2.StatisticPO{}, &po2.StatisticDataPO{})
	if err != nil {
		return err
	}
	return nil
}
