package entity

type Announcement struct {
	ID        int64 `gorm:"primaryKey"`
	Title     string
	Body      string
	URL       string
	Available bool
	CreatedAt int64
}

func (Announcement) TableName() string {
	return "announcement"
}
