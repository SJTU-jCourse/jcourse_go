package converter

import (
	model2 "jcourse_go/internal/model/model"
	"jcourse_go/internal/model/po"
)

func ConvertTeacherDetailFromPO(po *po.TeacherPO) model2.TeacherDetail {
	return model2.TeacherDetail{
		TeacherSummary: ConvertTeacherSummaryFromPO(po),
		Email:          po.Email,
		Code:           po.Code,
		Title:          po.Title,
		ProfileURL:     po.ProfileURL,
		Biography:      po.Biography,
	}
}

func PackTeacherWithCourses(t *model2.TeacherDetail, courses []model2.CourseSummary) {
	t.Courses = courses
}

func ConvertTeacherSummaryFromPO(po *po.TeacherPO) model2.TeacherSummary {
	return model2.TeacherSummary{
		ID:         po.ID,
		Name:       po.Name,
		Department: po.Department,
		Picture:    po.Picture,
	}
}

func PackTeacherWithRatingInfo(t *model2.TeacherSummary, rating model2.RatingInfo) {
	t.RatingInfo = rating
}
