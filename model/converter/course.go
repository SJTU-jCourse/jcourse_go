package converter

import (
	"strings"

	"jcourse_go/model/model"
	"jcourse_go/model/po"
)

func ConvertBaseCourseFromPO(po po.BaseCoursePO) model.BaseCourse {
	return model.BaseCourse{
		ID:     po.ID,
		Code:   po.Code,
		Name:   po.Name,
		Credit: po.Credit,
	}
}

func ConvertCourseSummaryFromPO(po po.CoursePO) model.CourseSummary {
	return model.CourseSummary{
		CourseMinimal: model.CourseMinimal{
			ID: po.ID,
			BaseCourse: model.BaseCourse{
				Code:   po.Code,
				Name:   po.Name,
				Credit: po.Credit,
			},
			MainTeacher: model.TeacherSummary{
				ID:   po.MainTeacherID,
				Name: po.MainTeacherName,
			},
		},
		Categories: nil,
		Department: po.Department,
		RatingInfo: model.RatingInfo{},
	}
}

func ConvertCourseSummariesFromPO(pos []po.CoursePO) []model.CourseSummary {
	res := make([]model.CourseSummary, 0)
	for _, v := range pos {
		res = append(res, ConvertCourseSummaryFromPO(v))
	}
	return res
}
func ConvertCourseDetailFromPO(po po.CoursePO) model.CourseDetail {
	return model.CourseDetail{
		CourseSummary: ConvertCourseSummaryFromPO(po),
	}
}

func PackCourseWithMainTeacher(c *model.CourseMinimal, teacher model.TeacherSummary) {
	c.MainTeacher = teacher
}

func PackCourseWithCategories(c *model.CourseSummary, categories []string) {
	c.Categories = categories
}

func PackCourseWithRatingInfo(c *model.CourseSummary, rating model.RatingInfo) {
	c.RatingInfo = rating
}

func PackCourseWithOfferedCourse(c *model.CourseDetail, offered []model.OfferedCourse) {
	c.OfferedCourse = offered
}

func ConvertOfferedCourseFromPO(po po.OfferedCoursePO) model.OfferedCourse {
	grade := strings.Split(po.Grade, ",")
	return model.OfferedCourse{
		ID:       int64(po.ID),
		Semester: po.Semester,
		Grade:    grade,
		Language: po.Language,
	}
}

func ConvertOfferedCoursesFromPOs(pos []po.OfferedCoursePO) []model.OfferedCourse {
	res := make([]model.OfferedCourse, 0, len(pos))
	for _, p := range pos {
		res = append(res, ConvertOfferedCourseFromPO(p))
	}
	return res
}

func PackOfferedCourseWithTeacher(c *model.OfferedCourse, teacher []model.TeacherDetail) {
	c.TeacherGroup = teacher
}
