package converter

import (
	"jcourse_go/model/model"
	"jcourse_go/model/po"
)

func ConvertTeacherDetailFromPO(po po.TeacherPO) model.TeacherDetail {
	return model.TeacherDetail{
		TeacherSummary: ConvertTeacherSummaryFromPO(po),
		Email:          po.Email,
		Code:           po.Code,
		Title:          po.Title,
		ProfileURL:     po.ProfileURL,
		Biography:      po.Biography,
	}
}

func PackTeacherWithOfferedCourses(t *model.TeacherDetail, courses []model.OfferedCourse) {
	t.Courses = courses
}

func ConvertTeacherSummaryFromPO(po po.TeacherPO) model.TeacherSummary {
	return model.TeacherSummary{
		ID:         int64(po.ID),
		Name:       po.Name,
		Department: po.Department,
		Picture:    po.Picture,
	}
}
