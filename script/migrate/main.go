package main

import (
	"log"

	flag "github.com/spf13/pflag"

	"jcourse_go/internal/config"
	"jcourse_go/internal/infrastructure/dal"
	"jcourse_go/internal/infrastructure/entity"
)

func main() {
	configPath := flag.StringP("config", "c", "config/config.yaml", "config file path")
	flag.Parse()

	c := config.InitConfig(*configPath)

	db, err := dal.NewPostgresSQL(c.DB)
	if err != nil {
		log.Fatal(err)
	}
	err = db.AutoMigrate(&entity.User{},
		&entity.Semester{}, &entity.Department{},
		&entity.Course{}, &entity.Teacher{},
		&entity.CourseOffering{}, &entity.CourseOfferingCategory{}, &entity.CourseOfferingTeacher{},
		&entity.TrainingPlan{}, &entity.Curriculum{}, &entity.TrainingPlanCurriculum{},
		&entity.Review{}, &entity.ReviewRevision{}, &entity.ReviewReaction{},
		&entity.UserPoint{}, &entity.Statistic{})
	if err != nil {
		log.Fatal(err)
	}
}
