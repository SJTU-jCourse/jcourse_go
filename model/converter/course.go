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
