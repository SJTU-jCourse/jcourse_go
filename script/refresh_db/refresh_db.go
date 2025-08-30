package main

import (
	"github.com/joho/godotenv"

	"jcourse_go/internal/dal"
	po2 "jcourse_go/internal/model/po"
	"jcourse_go/pkg/util"
)

func main() {
	_ = godotenv.Load()
	dal.InitDBClient()
	_ = util.InitSegWord()
	db := dal.GetDBClient()

	println("refreshing course...")
	courses := make([]po2.CoursePO, 0)
	if err := db.Model(&po2.CoursePO{}).Find(&courses).Error; err != nil {
		println(err)
	}
	for _, course := range courses {
		db.Save(&course)
	}

	println("refreshing teacher...")
	teachers := make([]po2.TeacherPO, 0)
	if err := db.Model(&po2.TeacherPO{}).Find(&teachers).Error; err != nil {
		println(err)
	}
	for _, teacher := range teachers {
		db.Save(&teacher)
	}

	println("refreshing training plan...")
	trainingPlans := make([]po2.TrainingPlanPO, 0)
	if err := db.Model(&po2.TrainingPlanPO{}).Find(&trainingPlans).Error; err != nil {
		println(err)
	}
	for _, tp := range trainingPlans {
		db.Save(&tp)
	}

	println("refreshing review...")
	reviews := make([]po2.ReviewPO, 0)
	if err := db.Model(&po2.ReviewPO{}).Find(&reviews).Error; err != nil {
		println(err)
	}
	for _, r := range reviews {
		db.Save(&r)
	}
}
