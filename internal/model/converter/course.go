package converter

import (
	"strings"

	model2 "jcourse_go/internal/model/model"
	"jcourse_go/internal/model/po"
)

func ConvertBaseCourseFromPO(po *po.BaseCoursePO) model2.BaseCourse {
	return model2.BaseCourse{
		ID:     po.ID,
		Code:   po.Code,
		Name:   po.Name,
		Credit: po.Credit,
	}
}

func ConvertCourseMinimalFromPO(po *po.CoursePO) model2.CourseMinimal {
	course := model2.CourseMinimal{
		ID: po.ID,
		BaseCourse: model2.BaseCourse{
			Code:   po.Code,
			Name:   po.Name,
			Credit: po.Credit,
		},
		MainTeacher: model2.TeacherSummary{
			ID:   po.MainTeacherID,
			Name: po.MainTeacherName,
		},
	}
	if po.MainTeacher != nil {
		course.MainTeacher = ConvertTeacherSummaryFromPO(po.MainTeacher)
	}
	return course
}

func ConvertCourseSummaryFromPO(po *po.CoursePO) model2.CourseSummary {
	courseSummary := model2.CourseSummary{
		CourseMinimal: ConvertCourseMinimalFromPO(po),
		Categories:    make([]string, 0),
		Department:    po.Department,
		RatingInfo:    model2.RatingInfo{},
	}
	if len(po.Categories) > 0 {
		for _, v := range po.Categories {
			courseSummary.Categories = append(courseSummary.Categories, v.Category)
		}
	}
	return courseSummary
}

func ConvertCourseDetailFromPO(po *po.CoursePO) model2.CourseDetail {
	courseDetail := model2.CourseDetail{
		CourseSummary: ConvertCourseSummaryFromPO(po),
		OfferedCourse: make([]model2.OfferedCourse, 0),
	}
	if len(po.OfferedCourses) > 0 {
		for _, v := range po.OfferedCourses {
			courseDetail.OfferedCourse = append(courseDetail.OfferedCourse, ConvertOfferedCourseFromPO(v))
		}
	}
	return courseDetail
}

func PackCourseWithRatingInfo(c *model2.CourseSummary, rating model2.RatingInfo) {
	c.RatingInfo = rating
}

func ConvertOfferedCourseFromPO(po po.OfferedCoursePO) model2.OfferedCourse {
	grade := strings.Split(po.Grade, ",")
	return model2.OfferedCourse{
		ID:       po.ID,
		Semester: po.Semester,
		Grade:    grade,
		Language: po.Language,
	}
}
