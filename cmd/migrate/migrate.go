package main

import (
	"jcourse_go/dal"
	"jcourse_go/model/po"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	dal.InitDBClient()
	db := dal.GetDBClient()
	err := db.AutoMigrate(&po.UserPO{},
		&po.BaseCoursePO{}, &po.CoursePO{}, &po.TeacherPO{}, &po.CourseCategoryPO{},
		&po.OfferedCoursePO{}, &po.OfferedCourseTeacherPO{},
		&po.ReviewPO{}, &po.TrainingPlanCoursePO{}, &po.TrainingPlanPO{},
		&po.TrainingPlanRatePO{})
	if err != nil {
		panic(err)
	}
}
