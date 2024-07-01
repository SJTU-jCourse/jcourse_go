package converter

import (
	"jcourse_go/model/domain"
	"jcourse_go/model/po"
)

func ConvertTeacherPOToDomain(teacher *po.TeacherPO) *domain.Teacher {
	if teacher == nil {
		return nil
	}
	return &domain.Teacher{
		ID:         int64(teacher.ID),
		Code:       teacher.Code,
		Name:       teacher.Name,
		Department: teacher.Department,
		Title:      teacher.Title,
	}
}
