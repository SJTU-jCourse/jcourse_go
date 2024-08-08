package main

import (
	"github.com/gin-gonic/gin"
	"jcourse_go/handler"
	"jcourse_go/middleware"
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
	needAuthGroup.Use(middleware.RequireAuth())

	teacherGroup := needAuthGroup.Group("/teacher")
	teacherGroup.GET("", handler.GetTeacherListHandler)
	teacherGroup.GET("/:teacherID", handler.GetTeacherDetailHandler)

	courseGroup := needAuthGroup.Group("/course")
	courseGroup.GET("", handler.GetCourseListHandler)
	courseGroup.GET("/suggest", handler.GetSuggestedCourseHandler)
	courseGroup.GET("/:courseID", handler.GetCourseDetailHandler)
	courseGroup.POST("/:courseID/watch", handler.WatchCourseHandler)
	courseGroup.POST("/:courseID/unwatch", handler.UnWatchCourseHandler)

	ratingGroup := needAuthGroup.Group("/rating")
	ratingGroup.POST("", handler.CreateRatingHandler)

	reviewGroup := needAuthGroup.Group("/review")
	reviewGroup.GET("", handler.GetReviewListHandler)
	reviewGroup.POST("", handler.CreateReviewHandler)
	reviewGroup.GET("/suggest", handler.GetSuggestedReviewHandler)
	reviewGroup.GET("/:reviewID", handler.GetReviewDetailHandler)
	reviewGroup.PUT("/:reviewID", handler.UpdateReviewHandler)
	reviewGroup.DELETE("/:reviewID", handler.DeleteReviewHandler)

	reviewReactionGroup := needAuthGroup.Group("/review-reaction")
	reviewReactionGroup.POST("", handler.CreateReviewReactionHandler)
	reviewReactionGroup.DELETE("/:reactionID", handler.DeleteReviewReactionHandler)

	userGroup := needAuthGroup.Group("/user")
	userGroup.GET("", handler.GetUserListHandler)
	userGroup.GET("/suggest", handler.GetSuggestedUserHandler)
	userGroup.GET("/me", handler.GetCurrentUserHandler)
	userGroup.GET("/:userID", handler.GetUserDetailHandler)
	userGroup.POST("/:userID/unwatch", handler.UnWatchUserHandler)

	adminGroup := needAuthGroup.Group("/admin")
	adminGroup.Use(middleware.RequireAdmin())

	adminGroup.GET("")
}
