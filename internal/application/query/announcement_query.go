package query

import (
	"context"

	"gorm.io/gorm"

	"jcourse_go/internal/application/vo"
	"jcourse_go/internal/infrastructure/entity"
)

type AnnouncementQueryService interface {
	GetAnnouncements(ctx context.Context) ([]vo.AnnouncementVO, error)
}

type announcementQueryService struct {
	db *gorm.DB
}

func (s *announcementQueryService) GetAnnouncements(ctx context.Context) ([]vo.AnnouncementVO, error) {
	announcements, err := gorm.G[entity.Announcement](s.db).Where("available = ?", true).Order("created_at desc").Find(ctx)

	if err != nil {
		return nil, err
	}

	result := make([]vo.AnnouncementVO, 0, len(announcements))
	for _, a := range announcements {
		result = append(result, vo.AnnouncementVO{
			ID:    a.ID,
			Title: a.Title,
			Body:  a.Body,
			URL:   a.URL,
		})
	}

	return result, nil
}

func NewAnnouncementQueryService(db *gorm.DB) AnnouncementQueryService {
	return &announcementQueryService{db: db}
}
