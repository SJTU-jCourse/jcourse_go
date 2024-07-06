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

func ConvertCoursePOToDomain(course *po.CoursePO) *domain.Course {
	if course == nil {
		return nil
	}
	return &domain.Course{
		ID:     int64(course.ID),
		Code:   course.Code,
		Name:   course.Name,
		Credit: course.Credit,
	}
}

func ConvertOfferedCoursePOToDomain(offeredCourse *po.OfferedCoursePO) *domain.OfferedCourse {
	if offeredCourse == nil {
		return nil
	}
	return &domain.OfferedCourse{}
}
