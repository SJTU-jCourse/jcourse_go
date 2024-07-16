package main

import (
	"github.com/joho/godotenv"
	"jcourse_go/dal"
	"jcourse_go/model/po"
)

func main() {
	_ = godotenv.Load()
	dal.InitDBClient()
	db := dal.GetDBClient()
	err := db.AutoMigrate(&po.UserPO{},
		&po.BaseCoursePO{}, &po.CoursePO{}, &po.TeacherPO{}, &po.CourseCategoryPO{},
		&po.OfferedCoursePO{}, &po.OfferedCourseTeacherPO{},
		&po.ReviewPO{}, &po.TrainingPlanCoursePO{}, &po.TrainingPlanPO{})
	if err != nil {
		panic(err)
	}
}
