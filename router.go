package main

import (
	"jcourse_go/handler"
	"jcourse_go/middleware"
	"jcourse_go/util"

	"github.com/gin-gonic/gin"
)

func registerRouter(r *gin.Engine) {
	middleware.InitSession(r)

	api := r.Group("/api")

	authGroup := api.Group("/auth")
	authGroup.POST("/login", handler.LoginHandler)
	authGroup.POST("/logout", handler.LogoutHandler)
	authGroup.POST("/register", handler.RegisterHandler)
	authGroup.POST("/send-verify-code", handler.SendVerifyCodeHandler)
	authGroup.POST("/reset-password", handler.ResetPasswordHandler)

	needAuthGroup := api.Group("")
	if !util.IsNoLoginMode() {
		needAuthGroup.Use(middleware.RequireAuth())
	}

	needAuthGroup.GET("/common", handler.GetCommonInfo)

	teacherGroup := needAuthGroup.Group("/teacher")
	teacherGroup.GET("", handler.GetTeacherListHandler)
	teacherGroup.GET("/:teacherID", handler.GetTeacherDetailHandler)

	baseCourseGroup := needAuthGroup.Group("/base_course")
	baseCourseGroup.GET("/:code", handler.GetBaseCourse)

	courseGroup := needAuthGroup.Group("/course")
	courseGroup.GET("", handler.GetCourseListHandler)
	// courseGroup.GET("/suggest", handler.GetSuggestedCourseHandler)
	courseGroup.GET("/:courseID", handler.GetCourseDetailHandler)
	// courseGroup.POST("/:courseID/watch", handler.WatchCourseHandler)
	// courseGroup.POST("/:courseID/unwatch", handler.UnWatchCourseHandler)

	trainingPlanGroup := needAuthGroup.Group("/training_plan")
	trainingPlanGroup.GET("", handler.SearchTrainingPlanHandler)
	trainingPlanGroup.GET("/:trainingPlanID", handler.GetTrainingPlanHandler)

	ratingGroup := needAuthGroup.Group("/rating")
	ratingGroup.POST("", handler.CreateRatingHandler)

	reviewGroup := needAuthGroup.Group("/review")
	reviewGroup.GET("", handler.GetReviewListHandler)
	reviewGroup.POST("", handler.CreateReviewHandler)
	// reviewGroup.GET("/suggest", handler.GetSuggestedReviewHandler)
	reviewGroup.GET("/:reviewID", handler.GetReviewDetailHandler)
	reviewGroup.PUT("/:reviewID", handler.UpdateReviewHandler)
	reviewGroup.DELETE("/:reviewID", handler.DeleteReviewHandler)

	reviewReactionGroup := needAuthGroup.Group("/review-reaction")
	reviewReactionGroup.POST("", handler.CreateReviewReactionHandler)
	reviewReactionGroup.DELETE("/:reactionID", handler.DeleteReviewReactionHandler)

	userGroup := needAuthGroup.Group("/user")
	userGroup.GET("", handler.GetUserListHandler)
	// userGroup.GET("/suggest", handler.GetSuggestedUserHandler)
	userGroup.GET("/:userID/activity", handler.GetUserActivityHandler)
	userGroup.GET("/:userID", handler.GetUserDetailHandler)
	// userGroup.POST("/:userID/watch", handler.WatchUserHandler)
	// userGroup.POST("/:userID/unwatch", handler.UnWatchUserHandler)
	userGroup.PUT("/:userID/profile", handler.UpdateUserProfileHandler)

	adminGroup := needAuthGroup.Group("/admin")
	adminGroup.Use(middleware.RequireAdmin())
	adminGroup.GET("/user", handler.AdminGetUserList)

	adminGroup.GET("")

	llmGroup := needAuthGroup.Group(("/llm"))
	llmGroup.GET("/review/opt", handler.OptCourseReviewHandler)
	llmGroup.GET("/course/summary/:courseID", handler.GetCourseSummaryHandler)
	llmGroup.GET("/course/match", handler.GetMatchCoursesHandler)

	if util.IsDebug() {
		llmGroup.GET("/vectorize/:courseID", handler.VectorizeCourseHandler)
	}
}
