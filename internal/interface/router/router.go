package router

import (
	"github.com/gin-gonic/gin"

	"jcourse_go/internal/app"
	"jcourse_go/internal/interface/controller"
)

func RegisterRouter(s *app.ServiceContainer) *gin.Engine {
	r := gin.Default()

	authController := controller.NewAuthController(s.Auth)
	courseController := controller.NewCourseController(s.CourseQuery, s.Notification)
	teacherController := controller.NewTeacherController(s.TeacherQuery)
	reviewController := controller.NewReviewController(s.ReviewQuery, s.ReviewCommand)
	reactionController := controller.NewReviewReactionController(s.Reaction)
	trainingPlanController := controller.NewTrainingPlanController(s.TrainingPlanQuery)
	userController := controller.NewUserController(s.UserQuery, s.UserProfileCommand)
	userPointController := controller.NewUserPointController(s.UserPointQuery)
	statisticController := controller.NewStatisticController(s.StatisticQuery)

	api := r.Group("/api")
	authGroup := api.Group("/auth")
	authGroup.POST("/login", authController.LoginHandler)
	authGroup.POST("/logout", authController.LogoutHandler)
	authGroup.POST("/register", authController.RegisterHandler)
	authGroup.POST("/send-verify-code", authController.SendVerifyCodeHandler)
	authGroup.POST("/reset-password", authController.ResetPasswordHandler)

	needAuthGroup := api.Group("")

	courseGroup := needAuthGroup.Group("/course")
	courseGroup.GET("", courseController.GetCourseList)
	courseGroup.GET("/filter", courseController.GetCourseFilter)
	courseGroup.GET("/:courseID", courseController.GetCourseDetail)
	courseGroup.POST("/:courseID/notification", courseController.ChangeNotification)

	teacherGroup := needAuthGroup.Group("/teacher")
	teacherGroup.GET("", teacherController.GetTeacherList)
	teacherGroup.GET("/filter", teacherController.GetTeacherFilter)
	teacherGroup.GET("/:teacherID", teacherController.GetTeacherDetail)

	trainingPlanGroup := needAuthGroup.Group("/training_plan")
	trainingPlanGroup.GET("", trainingPlanController.GetTrainingPlanList)
	trainingPlanGroup.GET("/filter", trainingPlanController.GetTrainingPlanFilter)
	trainingPlanGroup.GET("/:trainingPlanID", trainingPlanController.GetTrainingPlanDetail)

	reviewGroup := needAuthGroup.Group("/review")
	reviewGroup.GET("", reviewController.GetLatestReviews)
	reviewGroup.POST("", reviewController.CreateReview)
	reviewGroup.GET("/:reviewID", reviewController.GetReview)
	reviewGroup.PUT("/:reviewID", reviewController.UpdateReview)
	reviewGroup.DELETE("/:reviewID", reviewController.DeleteReview)

	reviewReactionGroup := needAuthGroup.Group("/review-reaction")
	reviewReactionGroup.POST("", reactionController.CreateReaction)
	reviewReactionGroup.DELETE("/:reactionID", reactionController.DeleteReaction)

	userGroup := needAuthGroup.Group("/user")
	userGroup.GET("/profile", userController.GetUserInfo)
	userGroup.PUT("/profile", userController.UpdateUserInfo)
	userGroup.GET("/review", reviewController.GetUserReviews)

	userPointGroup := userGroup.Group("/point")
	userPointGroup.GET("", userPointController.GetUserPoints)

	statisticGroup := needAuthGroup.Group("/statistic")
	statisticGroup.GET("", statisticController.GetStatistics)

	return r
}
