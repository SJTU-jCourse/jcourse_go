package dal

import (
	"gorm.io/gorm"

	entity2 "jcourse_go/internal/infrastructure/entity"
)

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&entity2.UserPO{},
		&entity2.BaseCourse{}, &entity2.Course{}, &entity2.TeacherPO{}, &entity2.CourseCategoryPO{},
		&entity2.OfferedCoursePO{}, &entity2.OfferedCourseTeacherPO{},
		&entity2.TrainingPlanPO{}, &entity2.TrainingPlanCoursePO{},
		&entity2.ReviewPO{}, &entity2.RatingPO{}, &entity2.ReviewRevisionPO{}, &entity2.ReviewReactionPO{},
		&entity2.SettingPO{}, &entity2.UserPointDetailPO{}, &entity2.StatisticPO{})
	if err != nil {
		return err
	}
	return nil
}
