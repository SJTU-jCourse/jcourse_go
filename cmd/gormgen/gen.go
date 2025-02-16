package main

import (
	"github.com/joho/godotenv"
	"gorm.io/gen"

	"jcourse_go/dal"
	"jcourse_go/model/po"
)

func main() {
	_ = godotenv.Load()
	dal.InitDBClient()
	db := dal.GetDBClient()

	g := gen.NewGenerator(gen.Config{
		OutPath: "./query",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})
	// gormdb, _ := gorm.Open(mysql.Open("root:@(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local"))
	g.UseDB(db) // reuse your gorm db

	// Generate basic type-safe DAO API for struct `model.User` following conventions
	g.ApplyBasic(po.UserPO{}, po.UserPointDetailPO{},
		po.TeacherPO{}, po.CoursePO{}, po.BaseCoursePO{}, po.OfferedCoursePO{}, po.TrainingPlanPO{},
		po.ReviewPO{}, po.RatingPO{}, po.ReviewReactionPO{}, po.ReviewRevisionPO{},
		po.SettingPO{}, po.StatisticPO{})

	// Generate the code
	g.Execute()
}
