package app

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"jcourse_go/internal/application/command"
	"jcourse_go/internal/application/query"
	"jcourse_go/internal/config"
	"jcourse_go/internal/infrastructure/dal"
	"jcourse_go/internal/infrastructure/repository"
	"jcourse_go/internal/infrastructure/rpc"
)

type ServiceContainer struct {
	DB                 *gorm.DB
	Redis              *redis.Client
	Auth               command.AuthService
	Reaction           command.ReactionService
	ReviewCommand      command.ReviewCommandService
	Notification       command.CourseNotificationService
	StatisticCommand   command.StatisticService
	UserProfileCommand command.UserProfileService

	CourseQuery       query.CourseQueryService
	ReviewQuery       query.ReviewQueryService
	StatisticQuery    query.StatisticQueryService
	TeacherQuery      query.TeacherQueryService
	TrainingPlanQuery query.TrainingPlanQueryService
	UserQuery         query.UserQueryService
	UserPointQuery    query.UserPointQueryService
	AnnouncementQuery query.AnnouncementQueryService
	ReportQuery       query.ReportQueryService
}

func NewServiceContainer(c config.AppConfig) (*ServiceContainer, error) {
	db, err := dal.NewPostgresSQL(c.Postgres)
	if err != nil {
		return nil, err
	}
	rdb, err := dal.NewRedisClient(c.Redis)
	if err != nil {
		return nil, err
	}

	codeRepo := repository.NewVerificationCodeRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	userProfileRepo := repository.NewUserProfileRepository(db)
	authUserRepo := repository.NewAuthUserRepository(db)
	reactionRepo := repository.NewReactionRepository(db)
	reviewRepo := repository.NewReviewRepository(db)
	courseRepo := repository.NewCourseRepository(db)
	statisticRepo := repository.NewStatisticRepository(db)
	notificationRepo := repository.NewCourseNotificationRepository(db)

	emailSender := rpc.NewSMTPEmailSender(c.SMTP)

	authCommandService := command.NewAuthService(codeRepo, sessionRepo, authUserRepo, emailSender)
	reviewCommandService := command.NewReviewCommandService(courseRepo, reviewRepo)
	notificationService := command.NewCourseNotificationService(courseRepo, notificationRepo)
	reactionService := command.NewReactionService(reactionRepo, reviewRepo)
	statisticService := command.NewStatisticService(db, statisticRepo)
	userProfileService := command.NewUserProfileService(userProfileRepo)

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
		AnnouncementQuery: query.NewAnnouncementQueryService(db),
		ReportQuery:       query.NewReportQueryService(db),

		Auth:               authCommandService,
		Reaction:           reactionService,
		Notification:       notificationService,
		ReviewCommand:      reviewCommandService,
		StatisticCommand:   statisticService,
		UserProfileCommand: userProfileService,
	}, nil
}
