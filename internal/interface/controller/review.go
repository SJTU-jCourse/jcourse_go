package controller

import (
	"github.com/gin-gonic/gin"

	"jcourse_go/internal/application/command"
	"jcourse_go/internal/application/query"
)

type ReviewController struct {
	reviewQuery   query.ReviewQueryService
	reviewCommand command.ReviewCommandService
}

func NewReviewController(
	reviewQuery query.ReviewQueryService,
	reviewCommand command.ReviewCommandService,
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
