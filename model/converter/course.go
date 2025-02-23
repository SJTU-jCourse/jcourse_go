package converter

import (
	"strings"

	"jcourse_go/model/model"
	"jcourse_go/model/po"
)

func ConvertBaseCourseFromPO(po *po.BaseCoursePO) model.BaseCourse {
	return model.BaseCourse{
		ID:     po.ID,
		Code:   po.Code,
		Name:   po.Name,
		Credit: po.Credit,
	}
}

func ConvertCourseMinimalFromPO(po *po.CoursePO) model.CourseMinimal {
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
	if po.MainTeacher != nil {
		course.MainTeacher = ConvertTeacherSummaryFromPO(po.MainTeacher)
	}
	return course
}

func ConvertCourseSummaryFromPO(po *po.CoursePO) model.CourseSummary {
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

func ConvertCourseDetailFromPO(po *po.CoursePO) model.CourseDetail {
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

func PackCourseWithRatingInfo(c *model.CourseSummary, rating model.RatingInfo) {
	c.RatingInfo = rating
}

func ConvertOfferedCourseFromPO(po po.OfferedCoursePO) model.OfferedCourse {
	grade := strings.Split(po.Grade, ",")
	return model.OfferedCourse{
		ID:       po.ID,
		Semester: po.Semester,
		Grade:    grade,
		Language: po.Language,
	}
}
