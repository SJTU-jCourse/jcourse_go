package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"jcourse_go/internal/application/command"
	"jcourse_go/internal/application/query"
	"jcourse_go/internal/domain/course"
	"jcourse_go/internal/domain/shared"
	"jcourse_go/internal/interface/dto"
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

func (c *ReviewController) GetReview(ctx *gin.Context) {
	reviewIDStr := ctx.Param("reviewID")
	if reviewIDStr == "" {
		dto.WriteBadArgumentResponse(ctx)
		return
	}
	reviewID, err := strconv.Atoi(reviewIDStr)
	if err != nil || reviewID == 0 {
		dto.WriteBadArgumentResponse(ctx)
		return
	}

	review, err := c.reviewQuery.GetReview(ctx, shared.IDType(reviewID))
	if err != nil {
		dto.WriteErrorResponse(ctx, err)
		return
	}
	dto.WriteDataResponse(ctx, review)
}

func (c *ReviewController) CreateReview(ctx *gin.Context) {
	var req course.WriteReviewCommand
	if ctx.ShouldBind(&req) != nil {
		dto.WriteBadArgumentResponse(ctx)
		return
	}

	reqCtx := shared.NewRequestCtx(0, shared.UserRoleNormal)
	err := c.reviewCommand.WriteReview(ctx, reqCtx, req)
	if err != nil {
		dto.WriteErrorResponse(ctx, err)
		return
	}
	dto.WriteDataResponse(ctx, nil)
}

func (c *ReviewController) UpdateReview(ctx *gin.Context) {
	var req course.UpdateReviewCommand
	if ctx.ShouldBind(&req) != nil {
		dto.WriteBadArgumentResponse(ctx)
		return
	}

	reqCtx := shared.NewRequestCtx(0, shared.UserRoleNormal)
	err := c.reviewCommand.UpdateReview(ctx, reqCtx, req)
	if err != nil {
		dto.WriteErrorResponse(ctx, err)
		return
	}
	dto.WriteDataResponse(ctx, nil)

}

func (c *ReviewController) DeleteReview(ctx *gin.Context) {
	var req course.DeleteReviewCommand
	if ctx.ShouldBind(&req) != nil {
		dto.WriteBadArgumentResponse(ctx)
		return
	}

	reqCtx := shared.NewRequestCtx(0, shared.UserRoleNormal)
	err := c.reviewCommand.DeleteReview(ctx, reqCtx, req)
	if err != nil {
		dto.WriteErrorResponse(ctx, err)
		return
	}
	dto.WriteDataResponse(ctx, nil)

}

func (c *ReviewController) GetLatestReviews(ctx *gin.Context) {
	var req shared.PaginationQuery
	if ctx.ShouldBind(&req) != nil {
		dto.WriteBadArgumentResponse(ctx)
		return
	}
	reviews, err := c.reviewQuery.GetLatestReviews(ctx, req)
	if err != nil {
		dto.WriteErrorResponse(ctx, err)
		return
	}
	dto.WriteDataResponse(ctx, reviews)
}

func (c *ReviewController) GetCourseReviews(ctx *gin.Context) {
	courseIDStr := ctx.Param("courseID")
	if courseIDStr == "" {
		dto.WriteBadArgumentResponse(ctx)
		return
	}
	courseID, err := strconv.Atoi(courseIDStr)
	if err != nil || courseID == 0 {
		dto.WriteBadArgumentResponse(ctx)
		return
	}
	var req shared.PaginationQuery
	if ctx.ShouldBind(&req) != nil {
		dto.WriteBadArgumentResponse(ctx)
		return
	}

	reviews, err := c.reviewQuery.GetCourseReviews(ctx, shared.IDType(courseID), req)
	if err != nil {
		dto.WriteErrorResponse(ctx, err)
		return
	}
	dto.WriteDataResponse(ctx, reviews)
}

func (c *ReviewController) GetUserReviews(ctx *gin.Context) {
	var req shared.PaginationQuery
	if ctx.ShouldBind(&req) != nil {
		dto.WriteBadArgumentResponse(ctx)
		return
	}

	userID := 0

	reviews, err := c.reviewQuery.GetUserReviews(ctx, shared.IDType(userID), req)
	if err != nil {
		dto.WriteErrorResponse(ctx, err)
		return
	}
	dto.WriteDataResponse(ctx, reviews)
}
