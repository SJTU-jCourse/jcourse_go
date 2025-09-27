package controller

import (
	"github.com/gin-gonic/gin"

	"jcourse_go/internal/application/query"
)

type AnnouncementController struct {
	announcementQuery query.AnnouncementQueryService
}

func NewAnnouncementController(
	announcementQuery query.AnnouncementQueryService,
) *AnnouncementController {
	return &AnnouncementController{
		announcementQuery: announcementQuery,
	}
}

func (c *AnnouncementController) GetAnnouncements(ctx *gin.Context) {
	announcements, err := c.announcementQuery.GetAnnouncements(ctx)
	if err != nil {
		WriteErrorResponse(ctx, err)
		return
	}
	WriteDataResponse(ctx, announcements)
}
