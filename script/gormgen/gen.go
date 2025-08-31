package main

import (
	"github.com/joho/godotenv"
	"gorm.io/gen"

	"jcourse_go/internal/dal"
	entity2 "jcourse_go/internal/infrastructure/entity"

	"jcourse_go/internal/entity"
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
	g.ApplyBasic(entity2.UserPO{}, entity2.UserPointDetailPO{},
		entity2.TeacherPO{}, entity2.Course{}, entity2.BaseCourse{}, entity2.OfferedCoursePO{}, entity2.TrainingPlanPO{},
		entity2.CourseCategoryPO{}, entity2.TrainingPlanCoursePO{},
		entity2.ReviewPO{}, entity2.RatingPO{}, entity2.ReviewReactionPO{}, entity2.ReviewRevisionPO{},
		entity2.SettingPO{}, entity2.StatisticPO{}, entity.StatisticDataPO{})

	// Generate the code
	g.Execute()
}
