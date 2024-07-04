package converter

import (
	"jcourse_go/model/domain"
	"jcourse_go/model/po"
)

func ConvertBaseCourseDomainToPO(course *domain.BaseCourse) *po.BaseCoursePO {
	if course == nil {
		return nil
	}
	return &po.BaseCoursePO{
		Code:   course.Code,
		Name:   course.Name,
		Credit: course.Credit,
	}
}

func ConvertBaseCoursePOToDomain(course *po.BaseCoursePO) *domain.BaseCourse {
	if course == nil {
		return nil
	}
	return &domain.BaseCourse{
		ID:     int64(course.ID),
		Code:   course.Code,
		Name:   course.Name,
		Credit: course.Credit,
	}
}
