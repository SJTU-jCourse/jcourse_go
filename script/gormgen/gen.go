package main

import (
	"github.com/joho/godotenv"
	"gorm.io/gen"

	"jcourse_go/internal/dal"
	po2 "jcourse_go/internal/model/po"
)

func main() {
	_ = godotenv.Load()
	dal.InitDBClient()
	db := dal.GetDBClient()

	g := gen.NewGenerator(gen.Config{
		OutPath: "./repository",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})
	// gormdb, _ := gorm.Open(mysql.Open("root:@(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local"))
	g.UseDB(db) // reuse your gorm db

	// Generate basic type-safe DAO API for struct `model.User` following conventions
	g.ApplyBasic(po2.UserPO{}, po2.UserPointDetailPO{},
		po2.TeacherPO{}, po2.CoursePO{}, po2.BaseCoursePO{}, po2.OfferedCoursePO{}, po2.TrainingPlanPO{},
		po2.CourseCategoryPO{}, po2.TrainingPlanCoursePO{},
		po2.ReviewPO{}, po2.RatingPO{}, po2.ReviewReactionPO{}, po2.ReviewRevisionPO{},
		po2.SettingPO{}, po2.StatisticPO{}, po2.StatisticDataPO{})

	// Generate the code
	g.Execute()
}
