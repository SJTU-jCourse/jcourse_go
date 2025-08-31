package main

import (
	"github.com/joho/godotenv"

	"jcourse_go/internal/dal"
	entity2 "jcourse_go/internal/infrastructure/entity"
	"jcourse_go/pkg/util"
)

func main() {
	_ = godotenv.Load()
	dal.InitDBClient()
	_ = util.InitSegWord()
	db := dal.GetDBClient()

	println("refreshing course...")
	courses := make([]entity2.Course, 0)
	if err := db.Model(&entity2.Course{}).Find(&courses).Error; err != nil {
		println(err)
	}
	for _, course := range courses {
		db.Save(&course)
	}

	println("refreshing teacher...")
	teachers := make([]entity2.TeacherPO, 0)
	if err := db.Model(&entity2.TeacherPO{}).Find(&teachers).Error; err != nil {
		println(err)
	}
	for _, teacher := range teachers {
		db.Save(&teacher)
	}

	println("refreshing training plan...")
	trainingPlans := make([]entity2.TrainingPlanPO, 0)
	if err := db.Model(&entity2.TrainingPlanPO{}).Find(&trainingPlans).Error; err != nil {
		println(err)
	}
	for _, tp := range trainingPlans {
		db.Save(&tp)
	}

	println("refreshing review...")
	reviews := make([]entity2.ReviewPO, 0)
	if err := db.Model(&entity2.ReviewPO{}).Find(&reviews).Error; err != nil {
		println(err)
	}
	for _, r := range reviews {
		db.Save(&r)
	}
}
