package main

import (
	handler2 "jcourse_go/internal/handler"
	middleware2 "jcourse_go/internal/middleware"
	"jcourse_go/pkg/util"

	"github.com/gin-gonic/gin"
)

func registerRouter(r *gin.Engine) {
	middleware2.InitSession(r)
	r.Use(middleware2.UV.UVStatistic())
	r.Use(middleware2.PV.PVStatistic())
	api := r.Group("/api")
	authGroup := api.Group("/auth")
	authGroup.POST("/login", handler2.LoginHandler)
	authGroup.POST("/logout", handler2.LogoutHandler)
	authGroup.POST("/register", handler2.RegisterHandler)
	authGroup.POST("/send-verify-code", handler2.SendVerifyCodeHandler)
	authGroup.POST("/reset-password", handler2.ResetPasswordHandler)

	needAuthGroup := api.Group("")
	if !util.IsNoLoginMode() {
		needAuthGroup.Use(middleware2.RequireAuth())
	}
	needAuthGroup.GET("/common", handler2.GetCommonInfo)

	teacherGroup := needAuthGroup.Group("/teacher")
	teacherGroup.GET("", handler2.GetTeacherListHandler)
	teacherGroup.GET("/filter", handler2.GetTeacherFilter)
	teacherGroup.GET("/:teacherID", handler2.GetTeacherDetailHandler)

	baseCourseGroup := needAuthGroup.Group("/base_course")
	baseCourseGroup.GET("/:code", handler2.GetBaseCourse)

	courseGroup := needAuthGroup.Group("/course")
	courseGroup.GET("", handler2.GetCourseListHandler)
	courseGroup.GET("/filter", handler2.GetCourseFilterHandler)
	// courseGroup.GET("/suggest", handler.GetSuggestedCourseHandler)
	courseGroup.GET("/:courseID", handler2.GetCourseDetailHandler)
	// courseGroup.POST("/:courseID/watch", handler.WatchCourseHandler)
	// courseGroup.POST("/:courseID/unwatch", handler.UnWatchCourseHandler)

	trainingPlanGroup := needAuthGroup.Group("/training_plan")
	trainingPlanGroup.GET("", handler2.SearchTrainingPlanHandler)
	trainingPlanGroup.GET("/filter", handler2.GetTrainingPlanFilter)
	trainingPlanGroup.GET("/:trainingPlanID", handler2.GetTrainingPlanHandler)

	ratingGroup := needAuthGroup.Group("/rating")
	ratingGroup.POST("", handler2.CreateRatingHandler)

	reviewGroup := needAuthGroup.Group("/review")
	reviewGroup.GET("", handler2.GetReviewListHandler)
	reviewGroup.POST("", handler2.CreateReviewHandler)
	// reviewGroup.GET("/suggest", handler.GetSuggestedReviewHandler)
	reviewGroup.GET("/:reviewID", handler2.GetReviewDetailHandler)
	reviewGroup.PUT("/:reviewID", handler2.UpdateReviewHandler)
	reviewGroup.DELETE("/:reviewID", handler2.DeleteReviewHandler)

	reviewReactionGroup := needAuthGroup.Group("/review-reaction")
	reviewReactionGroup.POST("", handler2.CreateReviewReactionHandler)
	reviewReactionGroup.DELETE("/:reactionID", handler2.DeleteReviewReactionHandler)

	userGroup := needAuthGroup.Group("/user")
	userGroup.GET("", handler2.GetUserListHandler)
	// userGroup.GET("/suggest", handler.GetSuggestedUserHandler)
	userGroup.GET("/:userID/activity", handler2.GetUserActivityHandler)
	userGroup.GET("/:userID", handler2.GetUserDetailHandler)
	// userGroup.POST("/:userID/watch", handler.WatchUserHandler)
	// userGroup.POST("/:userID/unwatch", handler.UnWatchUserHandler)
	userGroup.PUT("/:userID/profile", handler2.UpdateUserProfileHandler)

	userPointGroup := userGroup.Group("/point")
	userPointGroup.GET("", handler2.GetUserPointDetailListHandler)
	userPointGroup.POST("/transfer", handler2.TransferUserPointHandler)

	statisticGroup := needAuthGroup.Group("/statistic")
	statisticGroup.GET("", handler2.GetStatisticHandler)

	adminGroup := needAuthGroup.Group("/admin")
	adminGroup.Use(middleware2.RequireAdmin())
	adminGroup.GET("/user", handler2.AdminGetUserList)
	adminGroup.POST("/user/point/change", handler2.AdminChangeUserPoint)
	adminGroup.GET("/user/point/detail", handler2.AdminGetUserPointDetailList)
	adminGroup.GET("/user/point/transfer", handler2.AdminTransferUserPoint)

	llmGroup := needAuthGroup.Group("/llm")
	llmGroup.POST("/review/opt", handler2.OptCourseReviewHandler)
	llmGroup.GET("/course/summary/:courseID", handler2.GetCourseSummaryHandler)
	llmGroup.POST("/course/match", handler2.GetMatchCoursesHandler)

	if util.IsDebug() {
		llmGroup.GET("/vectorize/:courseID", handler2.VectorizeCourseHandler)
	}

}
