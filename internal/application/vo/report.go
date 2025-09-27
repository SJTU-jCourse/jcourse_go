package vo

import "jcourse_go/internal/infrastructure/entity"

type ReportVO struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	Content   string `json:"content"`
	Reply     string `json:"reply"`
	Solved    bool   `json:"solved"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

func NewReportVOFromEntity(e *entity.Report) ReportVO {
	return ReportVO{
		ID:        e.ID,
		UserID:    e.UserID,
		Content:   e.Content,
		Reply:     e.Reply,
		Solved:    e.Solved,
		CreatedAt: e.CreatedAt.Unix(),
		UpdatedAt: e.UpdatedAt.Unix(),
	}
}
