package entity

import "gorm.io/gorm"

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, &UserPoint{}, &Report{},
		&Semester{}, &Department{}, &Course{}, &Teacher{},
		&CourseNotification{}, &UserCourseEnrollment{},
		&CourseOffering{}, &CourseOfferingCategory{}, &CourseOfferingTeacher{},
		&TrainingPlan{}, &Curriculum{}, &TrainingPlanCurriculum{},
		&Review{}, &ReviewRevision{}, &ReviewReaction{},
		&Announcement{}, &ApiKey{}, &Statistic{},
	)
}
