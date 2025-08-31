package converter

import (
	"strings"

	"jcourse_go/internal/domain/course"
	"jcourse_go/internal/domain/rating"
	entity2 "jcourse_go/internal/infrastructure/entity"
)

func ConvertBaseCourseFromPO(po *entity2.BaseCourse) course.Curriculum {
	return course.Curriculum{
		ID:     po.ID,
		Code:   po.Code,
		Name:   po.Name,
		Credit: po.Credit,
	}
}

func ConvertCourseMinimalFromPO(po *entity2.Course) course.CourseMinimal {
	course := course.CourseMinimal{
		ID: po.ID,
		BaseCourse: course.Curriculum{
			Code:   po.Code,
			Name:   po.Name,
			Credit: po.Credit,
		},
		MainTeacher: course.TeacherSummary{
			ID:   po.MainTeacherID,
			Name: po.MainTeacherName,
		},
	}
	if po.MainTeacher != nil {
		course.MainTeacher = ConvertTeacherSummaryFromPO(po.MainTeacher)
	}
	return course
}

func ConvertCourseSummaryFromPO(po *entity2.Course) course.CourseSummary {
	courseSummary := course.CourseSummary{
		CourseMinimal: ConvertCourseMinimalFromPO(po),
		Categories:    make([]string, 0),
		Department:    po.Department,
		RatingInfo:    rating.RatingInfo{},
	}
	if len(po.Categories) > 0 {
		for _, v := range po.Categories {
			courseSummary.Categories = append(courseSummary.Categories, v.Category)
		}
	}
	return courseSummary
}

func ConvertCourseDetailFromPO(po *entity2.Course) course.CourseDetail {
	courseDetail := course.CourseDetail{
		CourseSummary: ConvertCourseSummaryFromPO(po),
		OfferedCourse: make([]course.CourseOffering, 0),
	}
	if len(po.OfferedCourses) > 0 {
		for _, v := range po.OfferedCourses {
			courseDetail.OfferedCourse = append(courseDetail.OfferedCourse, ConvertOfferedCourseFromPO(v))
		}
	}
	return courseDetail
}

func PackCourseWithRatingInfo(c *course.CourseSummary, rating rating.RatingInfo) {
	c.RatingInfo = rating
}

func ConvertOfferedCourseFromPO(po entity2.OfferedCoursePO) course.CourseOffering {
	grade := strings.Split(po.Grade, ",")
	return course.CourseOffering{
		ID:       po.ID,
		Semester: po.Semester,
		Grade:    grade,
		Language: po.Language,
	}
}
