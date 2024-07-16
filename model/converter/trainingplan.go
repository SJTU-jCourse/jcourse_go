package converter

import (
	"jcourse_go/model/domain"
	"jcourse_go/model/dto"
	"strconv"
)

func ConvertTrainingPlanDomainToDTO(domain domain.TrainingPlanDetail) dto.TrainingPlanListItemDTO{
	entryYear, err := strconv.Atoi(domain.EntryYear)
	if err != nil {
		entryYear=0
	}
	courses := make([]dto.BaseCourseDTO, 0)
	for _,c := range domain.Courses {
		courses = append(courses, ConvertBaseCourseDomainToDTO(c))
	}
	return dto.TrainingPlanListItemDTO{
		ID: domain.ID,
		Department: domain.Department,
		EntryYear: int64(entryYear),
		MajorName: domain.Major,
		Courses: courses,
	}
}

func ConvertTrainingPlanDomainListToDTO(domains []domain.TrainingPlanDetail) []dto.TrainingPlanListItemDTO{
	result := make([]dto.TrainingPlanListItemDTO, 0)
	for _,t := range domains{
		result = append(result, ConvertTrainingPlanDomainToDTO(t))
	}
	return result
}