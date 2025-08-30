package converter

import (
	model2 "jcourse_go/internal/model/model"
	"jcourse_go/internal/model/po"
)

func ConvertTrainingPlanSummaryFromPO(po *po.TrainingPlanPO) model2.TrainingPlanSummary {
	return model2.TrainingPlanSummary{
		ID:         po.ID,
		Code:       po.MajorCode,
		MajorName:  po.Major,
		EntryYear:  po.EntryYear,
		Department: po.Department,
		Degree:     po.Degree,
	}
}

func ConvertTrainingPlanDetailFromPO(po *po.TrainingPlanPO) model2.TrainingPlanDetail {
	detail := model2.TrainingPlanDetail{
		TrainingPlanSummary: ConvertTrainingPlanSummaryFromPO(po),
		TotalYear:           po.TotalYear,
		MinCredits:          po.MinCredits,
		MajorClass:          po.MajorClass,
		Courses:             make([]model2.TrainingPlanCourse, 0),
	}
	for _, baseCourse := range po.BaseCourses {
		detail.Courses = append(detail.Courses, ConvertTrainingPlanCourseFromPO(&baseCourse))
	}
	return detail
}

func ConvertTrainingPlanCourseFromPO(po *po.TrainingPlanCoursePO) model2.TrainingPlanCourse {
	tpCourse := model2.TrainingPlanCourse{
		BaseCourse: model2.BaseCourse{
			ID: po.BaseCourseID,
		},
		ID:              po.ID,
		SuggestSemester: po.SuggestSemester,
		Category:        po.Category,
	}
	if po.BaseCourse != nil {
		tpCourse.BaseCourse = ConvertBaseCourseFromPO(po.BaseCourse)
	}
	return tpCourse
}

func PackTrainingPlanWithRatingInfo(t *model2.TrainingPlanSummary, rating model2.RatingInfo) {
	t.RatingInfo = rating
}
