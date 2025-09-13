package main

import (
	"gorm.io/gen"

	"jcourse_go/internal/infrastructure/entity"
)

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "./internal/infrastructure/query",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	// Generate basic type-safe DAO API for struct `model.User` following conventions
	g.ApplyBasic(&entity.User{},
		&entity.Semester{}, &entity.Department{},
		&entity.Course{}, &entity.Teacher{},
		&entity.CourseOffering{}, &entity.CourseOfferingTeacher{}, &entity.CourseOfferingCategory{},
		&entity.TrainingPlan{}, &entity.TrainingPlanCurriculum{}, &entity.Curriculum{},
		&entity.Review{}, &entity.ReviewRevision{}, &entity.ReviewReaction{},
		&entity.UserPoint{}, &entity.Statistic{})

	// Generate the code
	g.Execute()
}
