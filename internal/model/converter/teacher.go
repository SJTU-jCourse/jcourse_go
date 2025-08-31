package converter

import (
	"jcourse_go/internal/domain/course"
	"jcourse_go/internal/domain/rating"
	"jcourse_go/internal/infrastructure/entity"
)

func ConvertTeacherDetailFromPO(po *entity.TeacherPO) course.TeacherDetail {
	return course.TeacherDetail{
		TeacherSummary: ConvertTeacherSummaryFromPO(po),
		Email:          po.Email,
		Code:           po.Code,
		Title:          po.Title,
		ProfileURL:     po.ProfileURL,
		Biography:      po.Biography,
	}
}

func PackTeacherWithCourses(t *course.TeacherDetail, courses []course.CourseSummary) {
	t.Courses = courses
}

func ConvertTeacherSummaryFromPO(po *entity.TeacherPO) course.TeacherSummary {
	return course.TeacherSummary{
		ID:         po.ID,
		Name:       po.Name,
		Department: po.Department,
		Picture:    po.Picture,
	}
}

func PackTeacherWithRatingInfo(t *course.TeacherSummary, rating rating.RatingInfo) {
	t.RatingInfo = rating
}
