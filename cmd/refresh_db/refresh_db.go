package main

import (
	"github.com/joho/godotenv"

	"jcourse_go/internal/infra"
	"jcourse_go/model/po"
	"jcourse_go/util"
)

func main() {
	_ = godotenv.Load()
	infra.InitDBClient()
	_ = util.InitSegWord()
	db := infra.GetDBClient()

	println("refreshing course...")
	courses := make([]po.CoursePO, 0)
	if err := db.Model(&po.CoursePO{}).Find(&courses).Error; err != nil {
		println(err)
	}
	for _, course := range courses {
		db.Save(&course)
	}

	println("refreshing teacher...")
	teachers := make([]po.TeacherPO, 0)
	if err := db.Model(&po.TeacherPO{}).Find(&teachers).Error; err != nil {
		println(err)
	}
	for _, teacher := range teachers {
		db.Save(&teacher)
	}

	println("refreshing training plan...")
	trainingPlans := make([]po.TrainingPlanPO, 0)
	if err := db.Model(&po.TrainingPlanPO{}).Find(&trainingPlans).Error; err != nil {
		println(err)
	}
	for _, tp := range trainingPlans {
		db.Save(&tp)
	}

	println("refreshing review...")
	reviews := make([]po.ReviewPO, 0)
	if err := db.Model(&po.ReviewPO{}).Find(&reviews).Error; err != nil {
		println(err)
	}
	for _, r := range reviews {
		db.Save(&r)
	}
}
