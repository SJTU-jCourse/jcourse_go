package converter

import (
	"strings"

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

func ConvertBaseCourseDomainToDTO(course domain.BaseCourse) dto.BaseCourseDTO {
	return dto.BaseCourseDTO{
		Code:   course.Code,
		Name:   course.Name,
		Credit: course.Credit,
	}

}

func ConvertCoursePOToDomain(course po.CoursePO) domain.Course {
	return domain.Course{
		ID:          int64(course.ID),
		Code:        course.Code,
		Name:        course.Name,
		Credit:      course.Credit,
		Department:  course.Department,
		MainTeacher: domain.Teacher{Name: course.MainTeacherName, ID: course.MainTeacherID},
	}
}

func PackCourseWithCategories(course *domain.Course, categories []string) {
	if categories == nil {
		return
	}
	if len(categories) == 0 {
		categories = make([]string, 0)
	}
	course.Categories = categories
}

func PackCourseWithMainTeacher(course *domain.Course, mainTeacherPO po.TeacherPO) {
	if course == nil {
		return
	}
	course.MainTeacher = *ConvertTeacherPOToDomain(&mainTeacherPO)
}

func PackCourseWithOfferedCourse(course *domain.Course, offeredCoursePOs []po.OfferedCoursePO) {
	if course == nil {
		return
	}
	offeredCourses := make([]domain.OfferedCourse, 0)
	for _, offeredCoursePO := range offeredCoursePOs {
		offeredCourse := ConvertOfferedCoursePOToDomain(offeredCoursePO)
		offeredCourses = append(offeredCourses, offeredCourse)
	}
	course.OfferedCourses = offeredCourses
}

func PackCourseWithReviewInfo(course *domain.Course, info po.CourseReviewInfo) {
	if course == nil {
		return
	}
	course.RatingInfo = domain.RatingInfo{
		Average: info.Average,
		Count:   info.Count,
	}
}

func ConvertCourseDomainToListDTO(course domain.Course) dto.CourseListItemDTO {
	mainTeacherDTO := dto.TeacherDTO{
		ID:         course.MainTeacher.ID,
		Code:       course.MainTeacher.Code,
		Name:       course.MainTeacher.Name,
		Department: course.MainTeacher.Department,
		Title:      course.MainTeacher.Title,
	}
	return dto.CourseListItemDTO{
		ID:          course.ID,
		Code:        course.Code,
		Name:        course.Name,
		Credit:      course.Credit,
		MainTeacher: mainTeacherDTO,
		Categories:  course.Categories,
		Department:  course.Department,
		ReviewInfo:  course.RatingInfo,
	}
}

func ConvertCourseListDomainToDTO(courses []domain.Course) []dto.CourseListItemDTO {
	result := make([]dto.CourseListItemDTO, 0, len(courses))
	if len(courses) == 0 {
		return result
	}
	for _, course := range courses {
		result = append(result, ConvertCourseDomainToListDTO(course))
	}
	return result
}

func ConvertOfferedCoursePOToDomain(offeredCourse po.OfferedCoursePO) domain.OfferedCourse {
	return domain.OfferedCourse{
		ID:       int64(offeredCourse.ID),
		Semester: offeredCourse.Semester,
		Language: offeredCourse.Language,
		Grade:    strings.Split(offeredCourse.Grade, ","),
	}
}

func ConvertOfferedCourseDomainToDTO(offeredCourse domain.OfferedCourse) dto.OfferedCourseDTO {
	offeredCourseDTO := dto.OfferedCourseDTO{
		ID:           offeredCourse.ID,
		Semester:     offeredCourse.Semester,
		Grade:        offeredCourse.Grade,
		Language:     offeredCourse.Language,
		TeacherGroup: make([]dto.TeacherDTO, 0),
	}
	for _, teacher := range offeredCourse.TeacherGroup {
		teacherDTO := ConvertTeacherDomainToDTO(teacher)
		offeredCourseDTO.TeacherGroup = append(offeredCourseDTO.TeacherGroup, teacherDTO)
	}
	return offeredCourseDTO
}

func ConvertCourseDomainToDetailDTO(course domain.Course) dto.CourseDetailDTO {
	courseDetailDTO := dto.CourseDetailDTO{
		ID:            course.ID,
		Code:          course.Code,
		Name:          course.Name,
		Credit:        course.Credit,
		MainTeacher:   ConvertTeacherDomainToDTO(course.MainTeacher),
		OfferedCourse: make([]dto.OfferedCourseDTO, 0),
		ReviewInfo:    course.RatingInfo,
	}
	for _, offeredCourse := range course.OfferedCourses {
		offeredCourseDTO := ConvertOfferedCourseDomainToDTO(offeredCourse)
		courseDetailDTO.OfferedCourse = append(courseDetailDTO.OfferedCourse, offeredCourseDTO)
	}
	return courseDetailDTO
}

func ConvertVectorToString(vector []float32) string {
	return "hello"
}
