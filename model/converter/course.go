package converter

import (
	"jcourse_go/model/domain"
	"jcourse_go/model/dto"
	"jcourse_go/model/po"
)

func ConvertBaseCourseDomainToPO(course domain.BaseCourse) po.BaseCoursePO {
	return po.BaseCoursePO{
		Code:   course.Code,
		Name:   course.Name,
		Credit: course.Credit,
	}
}

func ConvertBaseCoursePOToDomain(course po.BaseCoursePO) domain.BaseCourse {
	return domain.BaseCourse{
		ID:     int64(course.ID),
		Code:   course.Code,
		Name:   course.Name,
		Credit: course.Credit,
	}
}

func ConvertCoursePOToDomain(course po.CoursePO, categories []string) domain.Course {
	if categories == nil {
		categories = make([]string, 0)
	}
	return domain.Course{
		ID:          int64(course.ID),
		Code:        course.Code,
		Name:        course.Name,
		Credit:      course.Credit,
		Department:  course.Department,
		MainTeacher: domain.Teacher{Name: course.MainTeacherName, ID: course.MainTeacherID},
		Categories:  categories,
	}
}

func ConvertCourseDomainToListDTO(course domain.Course) dto.CourseListDTO {
	return dto.CourseListDTO{
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
	result := make([]dto.CourseListDTO, 0, len(courses))
	if len(courses) == 0 {
		return result
	}
	for _, course := range courses {
		result = append(result, ConvertCourseDomainToListDTO(course))
	}
	return result
}

func ConvertOfferedCoursePOToDomain(offeredCourse po.OfferedCoursePO) domain.OfferedCourse {
	return domain.OfferedCourse{}
}
