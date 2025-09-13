package app

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"jcourse_go/internal/application/command"
	"jcourse_go/internal/application/query"
	"jcourse_go/internal/config"
	"jcourse_go/internal/infrastructure/dal"
)

type ServiceContainer struct {
	DB            *gorm.DB
	Redis         *redis.Client
	Auth          command.AuthService
	Reaction      command.ReactionService
	ReviewCommand command.ReviewCommandService
	Subscription  command.SubscriptionService

	CourseQuery       query.CourseQueryService
	ReviewQuery       query.ReviewQueryService
	StatisticQuery    query.StatisticQueryService
	TeacherQuery      query.TeacherQueryService
	TrainingPlanQuery query.TrainingPlanQueryService
	UserQuery         query.UserQueryService
	UserPointQuery    query.UserPointQueryService
}

func NewServiceContainer(c config.AppConfig) (*ServiceContainer, error) {
	db, err := dal.NewPostgresSQL(c.DB)
	if err != nil {
		return nil, err
	}
	rdb, err := dal.NewRedisClient(c.Redis)
	if err != nil {
		return nil, err
	}
	return &ServiceContainer{
		DB:                db,
		Redis:             rdb,
		CourseQuery:       query.NewCourseQueryService(db),
		ReviewQuery:       query.NewReviewQueryService(db),
		StatisticQuery:    query.NewStatisticQueryService(db),
		TeacherQuery:      query.NewTeacherQueryService(db),
		TrainingPlanQuery: query.NewTrainingPlanQueryService(db),
		UserQuery:         query.NewUserQueryService(db),
		UserPointQuery:    query.NewUserPointQueryService(db),
	}, nil
}
