package query

import (
	"context"

	"gorm.io/gorm"

	"jcourse_go/internal/application/vo"
	"jcourse_go/internal/domain/shared"
	"jcourse_go/internal/infrastructure/entity"
)

type ReportQueryService interface {
	GetUserReports(ctx context.Context, userID shared.IDType, query shared.PaginationQuery) ([]vo.ReportVO, error)
}

type reportQueryService struct {
	db *gorm.DB
}

func (r *reportQueryService) GetUserReports(ctx context.Context, userID shared.IDType, query shared.PaginationQuery) ([]vo.ReportVO, error) {
	reports, err := gorm.G[entity.Report](r.db).Where("user_id = ?", userID).Order("created_at desc").Find(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]vo.ReportVO, 0)
	for _, report := range reports {
		result = append(result, vo.NewReportVOFromEntity(&report))
	}
	return result, nil
}

func NewReportQueryService(db *gorm.DB) ReportQueryService {
	return &reportQueryService{db: db}
}
