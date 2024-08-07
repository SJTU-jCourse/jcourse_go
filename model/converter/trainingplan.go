package converter

import (
	"jcourse_go/model/domain"
	"jcourse_go/model/dto"
	"jcourse_go/model/po"
	"strconv"
)

func ConvertTrainingPlanCourseDomainToDTO(courseDomain domain.TrainingPlanCourse) dto.TrainingPlanCourseDTO {
	return dto.TrainingPlanCourseDTO{
		ID:              courseDomain.ID,
		Name:            courseDomain.Name,
		Code:            courseDomain.Code,
		Credit:          courseDomain.Credit,
		SuggestYear:     courseDomain.SuggestYear,
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
		MinPoints:  domain.MinPoints,
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
		SuggestYear:     coursePO.SuggestYear,
		Department:      coursePO.Department,
	}
}

func ConvertTrainingPlanRatePOToDomain(ratePO po.TrainingPlanRatePO) domain.TrainingPlanRate {
	return domain.TrainingPlanRate{
		TrainingPlanID: ratePO.TrainingPlanID,
		Rate:           ratePO.Rate,
		UserID:         ratePO.UserID,
	}
}

func ConvertTrainingPlanRateDomainToDTO(rate domain.TrainingPlanRate) dto.TrainingPlanRateDTO {
	return dto.TrainingPlanRateDTO{
		TrainingPlanID: rate.TrainingPlanID,
		Rate:           rate.Rate,
		UserID:         rate.UserID,
	}
}

func ConvertTrainingPlanRateInfoDomainToDTO(rateInfo domain.TrainingPlanRateInfo) dto.TrainingPlanRateInfoDTO {
	rateMap := map[int64][]domain.TrainingPlanRate{}
	for _, r := range rateInfo.Rates {
		rateMap[r.Rate] = append(rateMap[r.Rate], r)
	}
	rateDist := make([]dto.TrainingPlanRateItem, 0)
	for k, v := range rateMap {
		rateDist = append(rateDist, dto.TrainingPlanRateItem{
			Rate:  float64(k),
			Count: int64(len(v)),
		})
	}
	return dto.TrainingPlanRateInfoDTO{
		Avg:      rateInfo.Avg,
		Count:    rateInfo.Count,
		RateDist: rateDist,
	}
}
