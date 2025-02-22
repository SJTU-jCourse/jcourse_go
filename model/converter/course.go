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

func ConvertCourseMinimalFromPO(po po.CoursePO) model.CourseMinimal {
	course := model.CourseMinimal{
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
	}
	if po.MainTeacher.ID != 0 {
		course.MainTeacher = ConvertTeacherSummaryFromPO(po.MainTeacher)
	}
	return course
}

func ConvertCourseSummaryFromPO(po po.CoursePO) model.CourseSummary {
	courseSummary := model.CourseSummary{
		CourseMinimal: ConvertCourseMinimalFromPO(po),
		Categories:    make([]string, 0),
		Department:    po.Department,
		RatingInfo:    model.RatingInfo{},
	}
	if len(po.Categories) > 0 {
		for _, v := range po.Categories {
			courseSummary.Categories = append(courseSummary.Categories, v.Category)
		}
	}
	return courseSummary
}

func ConvertCourseSummariesFromPO(pos []po.CoursePO) []model.CourseSummary {
	res := make([]model.CourseSummary, 0)
	for _, v := range pos {
		res = append(res, ConvertCourseSummaryFromPO(v))
	}
	return res
}
func ConvertCourseDetailFromPO(po po.CoursePO) model.CourseDetail {
	courseDetail := model.CourseDetail{
		CourseSummary: ConvertCourseSummaryFromPO(po),
		OfferedCourse: make([]model.OfferedCourse, 0),
	}
	if len(po.OfferedCourses) > 0 {
		for _, v := range po.OfferedCourses {
			courseDetail.OfferedCourse = append(courseDetail.OfferedCourse, ConvertOfferedCourseFromPO(v))
		}
	}
	return courseDetail
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
