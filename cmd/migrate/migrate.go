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
	// note: po with search index require manual migration
	err := db.AutoMigrate(&po.UserPO{},
		&po.BaseCoursePO{}, &po.CoursePO{}, &po.TeacherPO{}, &po.CourseCategoryPO{},
		&po.OfferedCoursePO{}, &po.OfferedCourseTeacherPO{},
		&po.ReviewPO{}, &po.RatingPO{}, &po.TrainingPlanPO{}, &po.TrainingPlanCoursePO{})
	if err != nil {
		panic(err)
	}
}
