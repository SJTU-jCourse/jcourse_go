package converter

import (
	"jcourse_go/model/model"
	"jcourse_go/model/po"
)

func ConvertTrainingPlanSummaryFromPO(po po.TrainingPlanPO) model.TrainingPlanSummary {
	return model.TrainingPlanSummary{
		ID:         int64(po.ID),
		Code:       po.MajorCode,
		MajorName:  po.Major,
		EntryYear:  po.EntryYear,
		Department: po.Department,
		Degree:     po.Degree,
	}
}

func ConvertTrainingPlanDetailFromPO(po po.TrainingPlanPO) model.TrainingPlanDetail {
	detail := model.TrainingPlanDetail{
		TrainingPlanSummary: ConvertTrainingPlanSummaryFromPO(po),
		TotalYear:           po.TotalYear,
		MinCredits:          po.MinCredits,
		MajorClass:          po.MajorClass,
		Courses:             make([]model.TrainingPlanCourse, 0),
	}
	for _, baseCourse := range po.BaseCourses {
		detail.Courses = append(detail.Courses, ConvertTrainingPlanCourseFromPO(baseCourse))
	}
	return detail
}

func PackTrainingPlanDetailWithCourse(tp *model.TrainingPlanDetail, courses []model.TrainingPlanCourse) {
	tp.Courses = courses
}

func ConvertTrainingPlanCourseFromPO(po po.TrainingPlanCoursePO) model.TrainingPlanCourse {
	tpCourse := model.TrainingPlanCourse{
		BaseCourse: model.BaseCourse{
			ID: po.BaseCourseID,
		},
		ID:              int64(po.ID),
		SuggestSemester: po.SuggestSemester,
		Category:        po.Category,
	}
	if po.BaseCourse.ID != 0 {
		tpCourse.BaseCourse = ConvertBaseCourseFromPO(po.BaseCourse)
	}
	return tpCourse
}

func PackTrainingPlanCourseWithBaseCourse(c *model.TrainingPlanCourse, course model.BaseCourse) {
	c.BaseCourse = course
}

func PackTrainingPlanWithRatingInfo(t *model.TrainingPlanSummary, rating model.RatingInfo) {
	t.RatingInfo = rating
}
