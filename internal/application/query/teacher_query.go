package query

import (
	"context"

	"gorm.io/gorm"

	"jcourse_go/internal/application/vo"
	"jcourse_go/internal/domain/course"
	"jcourse_go/internal/domain/shared"
	"jcourse_go/internal/infrastructure/entity"
)

type TeacherQueryService interface {
	GetTeacherList(ctx context.Context, query course.TeacherListQuery) ([]vo.TeacherListItemVO, error)
	GetTeacherFilter(ctx context.Context) (*course.TeacherFilter, error)
}

type teacherQueryService struct {
	db *gorm.DB
}

func (t *teacherQueryService) GetTeacherList(ctx context.Context, query course.TeacherListQuery) ([]vo.TeacherListItemVO, error) {
	te := make([]entity.Teacher, 0)
	db := t.db.WithContext(ctx).Model(&entity.Teacher{})
	if len(query.Titles) > 0 {
		db.Where("title IN (?)", query.Titles)
	}
	if len(query.Departments) > 0 {
		db.Where("department IN (?)", query.Departments)
	}
	if err := db.Find(&te).Error; err != nil {
		return nil, err
	}
	teVO := make([]vo.TeacherListItemVO, 0)
	for _, tt := range te {
		tVO := vo.NewTeacherListItemVOFromEntity(&tt)
		teVO = append(teVO, tVO)
	}
	return teVO, nil
}

func (t *teacherQueryService) GetTeacherDetail(ctx context.Context, teacherID shared.IDType) (*vo.TeacherDetailVO, error) {
	te := entity.Teacher{}
	if err := t.db.WithContext(ctx).
		Model(&entity.Teacher{}).
		Where("id = ?", teacherID).
		Take(&te).Error; err != nil {
		return nil, err
	}
	tVO := vo.NewTeacherDetailVOFromEntity(&te)
	return &tVO, nil
}

func (t *teacherQueryService) GetTeacherFilter(ctx context.Context) (*course.TeacherFilter, error) {
	filter := &course.TeacherFilter{
		Titles:      make([]course.FilterAggregateItem, 0),
		Departments: make([]course.FilterAggregateItem, 0),
	}

	if err := t.db.WithContext(ctx).Model(&entity.Teacher{}).
		Select("title as value, count(*) as count").
		Group("title").
		Order("title ASC").
		Scan(&filter.Titles).Error; err != nil {
		return nil, err
	}
	if err := t.db.WithContext(ctx).Model(&entity.Teacher{}).
		Select("department as value, count(*) as count").
		Group("department").
		Order("department ASC").
		Scan(&filter.Titles).Error; err != nil {
		return nil, err
	}
	return filter, nil
}

func NewTeacherQueryService(db *gorm.DB) TeacherQueryService {
	return &teacherQueryService{db: db}
}
