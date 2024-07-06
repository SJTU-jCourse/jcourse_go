package converter

import (
	"jcourse_go/model/domain"
	"jcourse_go/model/dto"
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

func ConvertCoursePOToDomain(course *po.CoursePO, categories []string) *domain.Course {
	if course == nil {
		return nil
	}
	return &domain.Course{
		ID:          int64(course.ID),
		Code:        course.Code,
		Name:        course.Name,
		Credit:      course.Credit,
		Department:  course.Department,
		MainTeacher: domain.Teacher{Name: course.MainTeacherName, ID: course.MainTeacherID},
		Categories:  categories,
	}
}

func ConvertCourseDomainToListDTO(course *domain.Course) *dto.CourseListDTO {
	if course == nil {
		return nil
	}
	return &dto.CourseListDTO{
		ID:              course.ID,
		Code:            course.Code,
		Name:            course.Name,
		Credit:          course.Credit,
		MainTeacherName: course.MainTeacher.Name,
		Categories:      course.Categories,
		Department:      course.Department,
	}
}

func ConvertCourseListDomainToDTO(courses []domain.Course) []dto.CourseListDTO {
	result := make([]dto.CourseListDTO, len(courses))
	if len(courses) == 0 {
		return result
	}
	for _, course := range courses {
		result = append(result, *ConvertCourseDomainToListDTO(&course))
	}
	return result
}

func ConvertOfferedCoursePOToDomain(offeredCourse *po.OfferedCoursePO) *domain.OfferedCourse {
	if offeredCourse == nil {
		return nil
	}
	return &domain.OfferedCourse{}
}
