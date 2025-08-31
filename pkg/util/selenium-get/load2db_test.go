package seleniumget

import (
	"fmt"
	"reflect"
	"testing"

	"jcourse_go/internal/dal"
	entity2 "jcourse_go/internal/infrastructure/entity"
	"jcourse_go/pkg/util"
)

func migrate() {
	db := dal.GetDBClient()
	err := db.AutoMigrate(&entity2.UserPO{},
		&entity2.BaseCourse{}, &entity2.Course{}, &entity2.TeacherPO{}, &entity2.CourseCategoryPO{},
		&entity2.OfferedCoursePO{}, &entity2.OfferedCourseTeacherPO{},
		&entity2.ReviewPO{}, &entity2.RatingPO{}, &entity2.TrainingPlanPO{}, &entity2.TrainingPlanCoursePO{})
	if err != nil {
		panic(err)
	}
}

func TestLoadTrainingPlan2DB(t *testing.T) {
	t.Run("simple print", func(t *testing.T) {
		LoadTrainingPlan2DB("../../data/trainingPlan.txt", nil)
	})
	t.Run("mem db", func(t *testing.T) {
		dal.InitTestMemDBClient()
		_ = util.InitSegWord()
		migrate()
		db := dal.GetDBClient()
		LoadTrainingPlan2DB("../../data/trainingPlan.txt", db)
	})
}
func TestLoadTeacherProfile2DB(t *testing.T) {
	t.Run("simple print", func(t *testing.T) {
		LoadTeacherProfile2DB("../../data/teachers.json", nil)
	})

	t.Run("duplicate teachers", func(t *testing.T) {
		teachers := LoadTeacherProfiles("../../data/teachers.json")
		codeMap := make(map[int64]string)
		for _, teacher := range teachers {
			if _, ok := codeMap[teacher.Code]; ok {
				fmt.Printf("duplicate teacher code: %d, name: (%s, %s)",
					teacher.Code, teacher.Name, codeMap[teacher.Code])
			}
			codeMap[teacher.Code] = teacher.Name
		}
	})

	t.Run("mem db", func(t *testing.T) {
		dal.InitTestMemDBClient()
		_ = util.InitSegWord()
		migrate()
		db := dal.GetDBClient()
		LoadTeacherProfile2DB("../../data/teachers.json", db)
	})
}

func TestLine2TrainingPlanMeta(t *testing.T) {
	lines := []string{
		"13050243,14761201813050243,视觉传达设计,13050243,4,28.5,艺术学,设计学院,本科,2018",
		"LXSGKPT026,147712018LXSGKPT026,留学生工科平台(中法),LXSGKPT026,4,127,未知类别,上海交大-巴黎高科卓越工程师学院,本科,2018",
		"LKSYB071,147812018LKSYB071,理科试验班类,LKSYB071,4,0,未知类别,数学科学学院,本科,2018",
		"110205TYMX,150612016110205TYMX,人力资源管理(体育明星班),110205TYMX,4,149,军事学,安泰经济与管理学院,本科,2016",
		"120402TY,150712016120402TY,行政管理...,120402TY,4,150,管理学,国际与公共事务学院,本科,2016",
	}
	example_result := LoadedTrainingPlan{
		Code:       "13050243",
		Name:       "视觉传达设计",
		TotalYear:  4,
		MinCredits: 28.5,
		MajorClass: "艺术学",
		Department: "设计学院",
		EntryYear:  2018,
		Degree:     "本科",
		Courses:    make([]LoadedCourse, 0),
	}
	t.Run("first full", func(t *testing.T) {
		if !reflect.DeepEqual(example_result, Line2TrainingPlanMeta(lines[0])) {
			t.Errorf("Expected %v, got %v", example_result, Line2TrainingPlanMeta(lines[0]))
		}
	})
	metas := make([]LoadedTrainingPlan, 0)
	for idx, l := range lines {
		name := fmt.Sprintf("test %d", idx)
		t.Run(name, func(t *testing.T) {
			metas = append(metas, Line2TrainingPlanMeta(l))
		})
	}

}

func TestLine2Course(t *testing.T) {
	lines := []string{
		"AD102,素描（1）,4.0,2018-2019,1,设计学院",
		"AD104,色彩（Ⅰ）,4.0,2018-2019,1,设计学院",
		"AD108,设计基础（A类）（1）,3.0,2018-2019,1,设计学院",
		"EN061,大学英语（1）,3.0,2018-2019,1,外国语学院",
		"PE001,体育（1）,1.0,2018-2019,1,体育系",
		"TH020,形势与政策,0.5,2018-2019,1,马克思主义学院",
		"PH247,工程物理与化学基础（1）,4.0,2019-2020,1,巴黎卓越工程师学院",
		"CH112,中国文化史（B类）（1）,3.0,2016-2017,1,马克思主义学院",
		"CH114,东、西方社交礼仪,3.0,2016-2017,1,马克思主义学院",
		"BS475,毕业设计（论文）（人力资源管理体育明星班）,16.0,2019-2020,2,体育系",
		"TY001,“UTJS”体验式教育:大学生演讲与沟通训练,2.0,2016-2017,1,学指委（学生处、团委、人武部）合署",
		"XXYZZ123, “select * from users”“”“”“SOME随机课程”“&&%……￥%##,13.5,2030-2031,3,测试测试测试",
	}

	example_result := LoadedCourse{
		Code:            "AD102",
		Name:            "素描（1）",
		Credit:          4.0,
		SuggestSemester: "2018-2019-1",
		Department:      "设计学院",
	}
	t.Run("first full", func(t *testing.T) {
		if !reflect.DeepEqual(example_result, Line2Course(lines[0])) {
			t.Errorf("Expected %v, got %v", example_result, Line2Course(lines[0]))
		}
	})
	courses := make([]LoadedCourse, 0)
	for idx, l := range lines {
		name := fmt.Sprintf("test %d", idx)
		t.Run(name, func(t *testing.T) {
			courses = append(courses, Line2Course(l))
		})
	}
}
