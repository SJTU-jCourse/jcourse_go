package announcement

import (
	"context"
	"time"

	"jcourse_go/internal/domain/shared"
)

type Announcement struct {
	ID shared.IDType

	Title   string
	Content string

	CreatedAt time.Time
}

type AnnouncementRepository interface {
	GetAvailableAnnouncement(ctx context.Context, now time.Time) ([]Announcement, error)
}
