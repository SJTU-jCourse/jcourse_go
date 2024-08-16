package converter

import (
	"jcourse_go/model/domain"
	"jcourse_go/model/dto"
	"jcourse_go/model/po"
)

func ConvertTeacherPOToDomain(teacher *po.TeacherPO) *domain.Teacher {
	if teacher == nil {
		return nil
	}

	return &domain.Teacher{
		ID:         int64(teacher.ID),
		Name:       teacher.Name,
		Email:      teacher.Email,
		Code:       teacher.Code,
		Department: teacher.Department,
		Title:      teacher.Title,
		Picture:    teacher.Picture,
		ProfileURL: teacher.ProfileURL,
		Biography:  teacher.Biography,
	}
}

func PackTeacherWithCourses(teacher *domain.Teacher, courses []po.OfferedCoursePO) {
	if teacher == nil {
		return
	}
	teacherCourses := make([]domain.OfferedCourse, 0)
	for _, offeredCoursePO := range courses {
		offeredCourse := ConvertOfferedCoursePOToDomain(offeredCoursePO)
		teacherCourses = append(teacherCourses, offeredCourse)
	}
	teacher.Courses = teacherCourses
}

func ConvertTeacherDomainToDTO(teacher domain.Teacher) dto.TeacherDTO {
	teacherDTO := dto.TeacherDTO{
		ID:         teacher.ID,
		Email:      teacher.Email,
		Code:       teacher.Code,
		Name:       teacher.Name,
		Department: teacher.Department,
		Title:      teacher.Title,
		Picture:    teacher.Picture,
		ProfileURL: teacher.ProfileURL,
		Biography:  teacher.Biography,
		Courses:    make([]dto.OfferedCourseDTO, 0),
	}
	for _, offeredCourse := range teacher.Courses {
		offeredCourseDTO := ConvertOfferedCourseDomainToDTO(offeredCourse)
		teacherDTO.Courses = append(teacherDTO.Courses, offeredCourseDTO)
	}
	return teacherDTO
}

func ConvertTeacherListDomainToDTO(teachers []domain.Teacher) []dto.TeacherDTO {
	result := make([]dto.TeacherDTO, 0, len(teachers))
	if len(teachers) == 0 {
		return result
	}
	for _, teacher := range teachers {
		result = append(result, ConvertTeacherDomainToDTO(teacher))
	}
	return result
}
