package controller

import (
	"github.com/gin-gonic/gin"

	"jcourse_go/internal/application"
)

type ReviewController struct {
	reviewQuery   application.ReviewQueryService
	reviewCommand application.ReviewCommandService
}

func NewReviewController(
	reviewQuery application.ReviewQueryService,
	reviewCommand application.ReviewCommandService,
) *ReviewController {
	return &ReviewController{
		reviewQuery:   reviewQuery,
		reviewCommand: reviewCommand,
	}
}

func (c *ReviewController) CreateReview(ctx *gin.Context) {

}

func (c *ReviewController) UpdateReview(ctx *gin.Context) {

}

func (c *ReviewController) DeleteReview(ctx *gin.Context) {

}

func (c *ReviewController) GetLatestReviews(ctx *gin.Context) {

}

func (c *ReviewController) GetCourseReviews(ctx *gin.Context) {

}
