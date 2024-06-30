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
	err := db.AutoMigrate(&po.UserPO{}, &po.BaseCoursePO{}, &po.CoursePO{}, &po.TeacherPO{}, &po.OfferedCourseCategoryPO{}, &po.OfferedCoursePO{}, &po.OfferedCourseTeacherPO{})
	if err != nil {
		panic(err)
	}
}
