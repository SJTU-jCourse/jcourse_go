package converter

import (
	"strconv"

	"jcourse_go/model/domain"
	"jcourse_go/model/dto"
	"jcourse_go/model/po"
)

func ConvertTrainingPlanCourseDomainToDTO(courseDomain domain.TrainingPlanCourse) dto.TrainingPlanCourseDTO {
	return dto.TrainingPlanCourseDTO{
		ID:              courseDomain.ID,
		Name:            courseDomain.Name,
		Code:            courseDomain.Code,
		Credit:          courseDomain.Credit,
		SuggestSemester: courseDomain.SuggestSemester,
	}
}
func ConvertTrainingPlanDomainToDTO(domain domain.TrainingPlanDetail) dto.TrainingPlanListItemDTO {
	entryYear, err := strconv.Atoi(domain.EntryYear)
	if err != nil {
		entryYear = 0
	}
	courses := make([]dto.TrainingPlanCourseDTO, 0)
	for _, c := range domain.Courses {
		courses = append(courses, ConvertTrainingPlanCourseDomainToDTO(c))
	}
	return dto.TrainingPlanListItemDTO{
		ID:         domain.ID,
		Department: domain.Department,
		EntryYear:  int64(entryYear),
		MajorName:  domain.Major,
		MinCredits: domain.MinCredits,
		MajorClass: domain.MajorClass,
		TotalYear:  int64(domain.TotalYear),
		Courses:    courses,
	}
}

func ConvertTrainingPlanDomainListToDTO(domains []domain.TrainingPlanDetail) []dto.TrainingPlanListItemDTO {
	result := make([]dto.TrainingPlanListItemDTO, 0)
	for _, t := range domains {
		result = append(result, ConvertTrainingPlanDomainToDTO(t))
	}
	return result
}

func ConvertTrainingPlanCoursePOToDomain(coursePO po.TrainingPlanCoursePO, baseCoursePO po.BaseCoursePO) domain.TrainingPlanCourse {
	return domain.TrainingPlanCourse{
		Code:            baseCoursePO.Code,
		Name:            baseCoursePO.Name,
		Credit:          baseCoursePO.Credit,
		SuggestSemester: coursePO.SuggestSemester,
		Department:      coursePO.Department,
	}
}

func PackTrainingPlanWithCourses(trainingPlan *domain.TrainingPlan, courses []domain.BaseCourse) {
	if courses == nil {
		return
	}
	if len(courses) == 0 {
		courses = make([]domain.BaseCourse, 0)
	}
	trainingPlan.Courses = courses
}
func PackTrainingPlanDetailWithCourses(trainingPlan *domain.TrainingPlanDetail, courses []domain.TrainingPlanCourse) {
	if courses == nil {
		return
	}
	if len(courses) == 0 {
		courses = make([]domain.TrainingPlanCourse, 0)
	}
	trainingPlan.Courses = courses
}

func ConvertTrainingPlanPOToDomain(trainingPlan po.TrainingPlanPO) domain.TrainingPlanDetail {
	return domain.TrainingPlanDetail{
		ID:         int64(trainingPlan.ID),
		Major:      trainingPlan.Major,
		Department: trainingPlan.Department,
		EntryYear:  trainingPlan.EntryYear,
		MajorCode:  trainingPlan.MajorCode,
		MajorClass: trainingPlan.MajorClass,
		MinCredits: trainingPlan.MinCredits,
		TotalYear:  trainingPlan.TotalYear,
	}
}
