package dal

import (
	"gorm.io/gorm"

	"jcourse_go/internal/infrastructure/entity"
)

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&entity.UserPO{},
		&entity.BaseCourse{}, &entity.Course{}, &entity.TeacherPO{},
		&entity.OfferedCoursePO{}, &entity.OfferedCourseTeacherPO{},
		&entity.TrainingPlanPO{}, &entity.TrainingPlanCoursePO{},
		&entity.ReviewPO{}, &entity.RatingPO{}, &entity.ReviewRevisionPO{}, &entity.ReviewReactionPO{},
		&entity.UserPointDetailPO{}, &entity.StatisticPO{})
	if err != nil {
		return err
	}
	return nil
}
