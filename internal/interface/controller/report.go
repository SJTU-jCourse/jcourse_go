package controller

import (
	"github.com/gin-gonic/gin"

	"jcourse_go/internal/application/query"
	"jcourse_go/internal/domain/shared"
)

type ReportController struct {
	reportQuery query.ReportQueryService
}

func NewReportController(reportQuery query.ReportQueryService) *ReportController {
	return &ReportController{
		reportQuery: reportQuery,
	}
}

func (c *ReportController) GetUserReports(ctx *gin.Context) {
	var req shared.PaginationQuery
	if ctx.ShouldBind(&req) != nil {
		WriteBadArgumentResponse(ctx)
		return
	}

	userID := shared.IDType(0)

	reports, err := c.reportQuery.GetUserReports(ctx, userID, req)
	if err != nil {
		WriteErrorResponse(ctx, err)
		return
	}
	WriteDataResponse(ctx, reports)
}
