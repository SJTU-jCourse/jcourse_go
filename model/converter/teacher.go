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
	}
}

func ConvertTeacherDomainToDTO(teacher domain.Teacher) dto.TeacherDTO {
	return dto.TeacherDTO{
		ID:         teacher.ID,
		Email:      teacher.Email,
		Name:       teacher.Name,
		Code:       teacher.Code,
		Department: teacher.Department,
		Title:      teacher.Title,
	}
}
