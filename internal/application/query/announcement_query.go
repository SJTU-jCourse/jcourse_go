package query

import (
	"context"
	"errors"

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
	if s.db == nil {
		return nil, errors.New("db not initialized")
	}

	var announcements []entity.Announcement
	err := s.db.WithContext(ctx).
		Where("available = ?", true).
		Order("created_at DESC").
		Find(&announcements).Error
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
