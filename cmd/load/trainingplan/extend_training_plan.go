package main

import (
	"errors"
	"jcourse_go/dal"
	"jcourse_go/model/po"
	seleniumget "jcourse_go/util/selenium-get"
	"log"
	"os"
	"strconv"

	"gorm.io/gorm"
)

func main() {
	dal.InitDBClient()
	db := dal.GetDBClient()
	data_path := "./data/trainingPlan.txt"
	log_file, err := os.OpenFile("./log/logfile.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer log_file.Close()
	defer log.SetOutput(os.Stdout)
	log.SetOutput(log_file)
	allTrainingPlans := seleniumget.LoadTrainingPLans(data_path)

	// 对齐trainingplan course 和 basecourse
	var all_courses []po.BaseCoursePO
	db.Model(po.BaseCoursePO{}).Find(&all_courses)

	for _, tp := range allTrainingPlans {
		tp_po := po.TrainingPlanPO{
			Degree:     tp.Degree,
			Major:      tp.Name,
			Department: tp.Department,
			EntryYear:  strconv.Itoa(tp.EntryYear),
			TotalYear:  tp.TotalYear,
			MinCredits: tp.MinCredits,
			MajorCode:  tp.Code,
			MajorClass: tp.MajorClass,
		}
		tp_po.GenerateSearchIndex()
		result := db.Model(po.TrainingPlanPO{}).Create(&tp_po)
		if result.Error != nil {
			log.Fatalf("In create training plan %#v:%#v", tp, result.Error)
		}
		for _, c := range tp.Courses {
			course := po.BaseCoursePO{}
			cresult := db.Model(po.BaseCoursePO{}).Where("code = ?", c.Code).First(&course)
			if cresult.Error != nil {
				if !errors.Is(cresult.Error, gorm.ErrRecordNotFound) {
					log.Fatalf("In bind course %#v totraining plan %#v:%#v", c, tp, cresult.Error)
				}
				// HINT:for production
				log.Printf("In bind course %#v totraining plan %#v:course not found", c, tp)
				continue
			}
			tpc_po := po.TrainingPlanCoursePO{
				TrainingPlanID:  int64(tp_po.ID),
				CourseID:        int64(course.ID),
				SuggestSemester: c.SuggestSemester,
				Department:      c.Department,
			}
			cresult = db.Model(po.TrainingPlanCoursePO{}).Create(&tpc_po)
			if cresult.Error != nil {
				if !errors.Is(cresult.Error, gorm.ErrRecordNotFound) {
					log.Fatalf("In bind course %#v totraining plan %#v:%#v", c, tp, cresult.Error)
				}
				log.Printf("In bind course %#v totraining plan %#v:course not found", c, tp)
			}
		}
	}
}
