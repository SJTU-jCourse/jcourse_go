package controller

import (
	"github.com/gin-gonic/gin"

	"jcourse_go/internal/domain/course"
)

type ReviewController struct {
	reviewRepo       course.ReviewRepository
	reviewCmdService course.ReviewCommandService
}

func NewReviewController(
	reviewRepo course.ReviewRepository,
	reviewCmdService course.ReviewCommandService,
) *ReviewController {
	return &ReviewController{
		reviewRepo:       reviewRepo,
		reviewCmdService: reviewCmdService,
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
