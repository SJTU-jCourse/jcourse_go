package converter

import (
	"jcourse_go/internal/domain/course"
	"jcourse_go/internal/domain/rating"
	"jcourse_go/internal/infrastructure/entity"
)

func ConvertTrainingPlanSummaryFromPO(po *entity.TrainingPlanPO) course.TrainingPlanSummary {
	return course.TrainingPlanSummary{
		ID:         po.ID,
		Code:       po.MajorCode,
		MajorName:  po.Major,
		EntryYear:  po.EntryYear,
		Department: po.Department,
		Degree:     po.Degree,
	}
}

func ConvertTrainingPlanDetailFromPO(po *entity.TrainingPlanPO) course.TrainingPlanDetail {
	detail := course.TrainingPlanDetail{
		TrainingPlanSummary: ConvertTrainingPlanSummaryFromPO(po),
		TotalYear:           po.TotalYear,
		MinCredits:          po.MinCredits,
		MajorClass:          po.MajorClass,
		Courses:             make([]course.TrainingPlanCourse, 0),
	}
	for _, baseCourse := range po.BaseCourses {
		detail.Courses = append(detail.Courses, ConvertTrainingPlanCourseFromPO(&baseCourse))
	}
	return detail
}

func ConvertTrainingPlanCourseFromPO(po *entity.TrainingPlanCoursePO) course.TrainingPlanCourse {
	tpCourse := course.TrainingPlanCourse{
		BaseCourse: course.Curriculum{
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

func PackTrainingPlanWithRatingInfo(t *course.TrainingPlanSummary, rating rating.RatingInfo) {
	t.RatingInfo = rating
}
