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
	return model.TrainingPlanDetail{
		TrainingPlanSummary: ConvertTrainingPlanSummaryFromPO(po),
		TotalYear:           po.TotalYear,
		MinCredits:          po.MinCredits,
		MajorClass:          po.MajorClass,
	}
}

func PackTrainingPlanDetailWithCourse(tp *model.TrainingPlanDetail, courses []model.TrainingPlanCourse) {
	tp.Courses = courses
}

func ConvertTrainingPlanCourseFromPO(po po.TrainingPlanCoursePO) model.TrainingPlanCourse {
	return model.TrainingPlanCourse{
		BaseCourse: model.BaseCourse{
			ID: po.BaseCourseID,
		},
		ID:              int64(po.ID),
		SuggestSemester: po.SuggestSemester,
		Department:      po.Department,
	}
}

func PackTrainingPlanCourseWithBaseCourse(c *model.TrainingPlanCourse, course model.BaseCourse) {
	c.BaseCourse = course
}
