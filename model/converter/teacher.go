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
		ID:          int64(teacher.ID),
		Name:        teacher.Name,
		Email:       teacher.Email,
		Code:        teacher.Code,
		Department:  teacher.Department,
		Title:       teacher.Title,
		Picture:     teacher.Picture,
		ProfileURL:  teacher.ProfileURL,
		ProfileDesc: teacher.ProfileDesc,
	}
}

// FIXME: BaseCourse可能包含了某学期未开课的课程, 是否需要换成OfferedCourse
func PackTeacherWithCourses(teacher *domain.Teacher, courses []po.BaseCoursePO) {
	if teacher == nil {
		return
	}
	teacherCourses := make([]domain.BaseCourse, 0)
	for _, baseCoursePO := range courses {
		baseCourse := ConvertBaseCoursePOToDomain(baseCoursePO)
		teacherCourses = append(teacherCourses, baseCourse)
	}
	teacher.Courses = teacherCourses
}

func ConvertTeacherDomainToDTO(teacher domain.Teacher) dto.TeacherDTO {
	teacherDTO := dto.TeacherDTO{
		ID:          teacher.ID,
		Email:       teacher.Email,
		Code:        teacher.Code,
		Name:        teacher.Name,
		Department:  teacher.Department,
		Title:       teacher.Title,
		Picture:     teacher.Picture,
		ProfileURL:  teacher.ProfileURL,
		ProfileDesc: teacher.ProfileDesc,
		Courses:     make([]dto.BaseCourseDTO, 0),
	}
	for _, baseCourse := range teacher.Courses {
		baseCourseDTO := ConvertBaseCourseDomainToDTO(baseCourse)
		teacherDTO.Courses = append(teacherDTO.Courses, baseCourseDTO)
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
