package seleniumget

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"jcourse_go/model/po"

	"gorm.io/gorm"
)

type LoadedCourse struct {
	Code            string
	Name            string
	Credit          float32
	Department      string
	SuggestSemester string
}
type LoadedTrainingPlan struct {
	Code       string
	Name       string
	TotalYear  int
	MinCredits float64
	MajorClass string
	Department string
	EntryYear  int
	Degree     string
	Courses    []LoadedCourse
}

func LoadedCourse2PO(course LoadedCourse) (po.BaseCoursePO, po.TrainingPlanCoursePO) {
	return po.BaseCoursePO{
			Code:   course.Code,
			Name:   course.Name,
			Credit: float64(course.Credit),
		}, po.TrainingPlanCoursePO{
			SuggestSemester: course.SuggestSemester,
			Department:      course.Department,
		}
}
func Line2Course(line string) LoadedCourse {
	// e.g. CH118,初级汉语精读（1）,8.0,2018-2019,1,巴黎卓越工程师学院
	meta := strings.Split(line, ",")
	if len(meta) != 6 {
		panic("Invalid line2course: " + line)
	}
	credit, err := strconv.ParseFloat(meta[2], 32)
	if err != nil {
		credit = 0.0
	}
	suggest_semester := fmt.Sprintf("%s-%s", meta[3], meta[4]) // 2018-2019-1
	return LoadedCourse{
		Code:            meta[0],
		Name:            meta[1],
		Credit:          float32(credit),
		SuggestSemester: suggest_semester,
		Department:      meta[5],
	}
}
func Line2TrainingPlanMeta(line string) LoadedTrainingPlan {
	plan := LoadedTrainingPlan{
		Courses: make([]LoadedCourse, 0),
	}
	metaInfo := strings.Split(line, ",")
	if len(metaInfo) != 10 {
		panic("Invalid line2trainingplan: " + line)
	}
	plan.Name = metaInfo[2]
	plan.Code = metaInfo[0]
	plan.TotalYear, _ = strconv.Atoi(metaInfo[4])
	plan.MinCredits, _ = strconv.ParseFloat(metaInfo[5], 64)
	plan.MajorClass = metaInfo[6]
	plan.Department = metaInfo[7]
	plan.Degree = metaInfo[8]
	plan.EntryYear, _ = strconv.Atoi(metaInfo[9])
	return plan
}
func Lines2TrainingPlan(lines []string) LoadedTrainingPlan {
	plan := Line2TrainingPlanMeta(lines[0])
	for _, line := range lines[1:] {
		plan.Courses = append(plan.Courses, Line2Course(line))
	}
	return plan
}
func TrainingPlan2PO(plan LoadedTrainingPlan) po.TrainingPlanPO {
	return po.TrainingPlanPO{
		Degree:     plan.Degree,
		Major:      plan.Name,
		Department: plan.Department,
		EntryYear:  strconv.Itoa(plan.EntryYear),
		TotalYear:  int64(plan.TotalYear),
		MinCredits: plan.MinCredits,
		MajorCode:  plan.Code,
		MajorClass: plan.MajorClass,
	}
}
func SaveTrainingPlanCourses(plan LoadedTrainingPlan, db *gorm.DB, tid int64) {
	if db == nil {
		for _, c := range plan.Courses {
			println("saved: " + c.Code)
		}
	}
	if len(plan.Courses) == 0 {
		return
	}
	baseCourseIDs := make([]int64, 0)
	codes := make([]string, 0)
	for _, c := range plan.Courses {
		codes = append(codes, c.Code)
	}

	result := db.Model(&po.BaseCoursePO{}).Where("code in ?", codes).Pluck("id", &baseCourseIDs)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return
	}
	if result.Error != nil {
		log.Fatalf("Failed to load courses: %v", result.Error)
	}
	trainingPlanCourses := make([]po.TrainingPlanCoursePO, 0)
	for _, id := range baseCourseIDs {
		trainingPlanCourses = append(trainingPlanCourses, po.TrainingPlanCoursePO{
			TrainingPlanID: tid,
			BaseCourseID:   id,
		})
	}
	db.CreateInBatches(&trainingPlanCourses, 100)
}
func SaveTrainingPlan(plans []LoadedTrainingPlan, db *gorm.DB) {
	if db == nil {
		// HINT:DEBUG only
		for _, plan := range plans[:10] {
			jsonData, _ := json.Marshal(TrainingPlan2PO(plan))
			println("saved: " + string(jsonData))
		}
		return
	}

	for _, plan := range plans {
		trainingPlanPO := TrainingPlan2PO(plan)
		db.Save(&trainingPlanPO)
		SaveTrainingPlanCourses(plan, db, int64(trainingPlanPO.ID))
	}
}
func LoadTrainingPLans(from string) []LoadedTrainingPlan {
	file, err := os.Open(from)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var lines []string
	var allTrainingPlans []LoadedTrainingPlan
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			if len(lines) > 0 {
				allTrainingPlans = append(allTrainingPlans, Lines2TrainingPlan(lines))
				lines = lines[:0]
			}
			continue
		}
		lines = append(lines, text)
	}
	return allTrainingPlans
}
func LoadTrainingPlan2DB(from string, db *gorm.DB) {
	allTrainingPlans := LoadTrainingPLans(from)
	// HINT: debug only
	SaveTrainingPlan(allTrainingPlans, db)
}

type LoadedTeacher struct {
	/**
		    {
	        "name": "\u5434\u7d2b\u4e91",
	        "code": 15036,
	        "department": "\u519c\u4e1a\u4e0e\u751f\u7269\u5b66\u9662",
	        "title": "\u526f\u6559\u6388",
	        "pinyin": "wuziyun",
	        "pinyin_abbr": "wzy",
	        "profile_url": "http://faculty.sjtu.edu.cn/wuziyun/zh_CN/index.htm",
	        "head_image": "https://faculty.sjtu.edu.cn//__local/7/F7/22/4B543D5D54790BC2EFC0192E923_DFEF7F92_1D826.png",
	        "mail": "wuziyun@vip.qq.com",
	        "profile_description": "PI\u7b80\u4ecb\uff1a\u5434\u7d2b\u4e91\u535a\u58eb\uff0c\u535a\u5bfc\uff0c\u6e56\u5357\u5e38\u5fb7\u4eba\uff0c2014\u5e74\u6bd5\u4e1a\u4e8e\u65b0\u52a0\u5761..."
	    },*/
	Name       string `json:"name"`
	Code       int64  `json:"code"`
	Department string `json:"department"`
	Title      string `json:"title"`
	Pinyin     string `json:"pinyin"`
	PinyinAbbr string `json:"pinyin_abbr"`
	ProfileUrl string `json:"profile_url"`
	HeadImage  string `json:"head_image"`
	Mail       string `json:"mail"`
	Biography  string `json:"biography"`
}

func Teacher2PO(teacher LoadedTeacher) po.TeacherPO {
	return po.TeacherPO{
		Name:       teacher.Name,
		Code:       strconv.Itoa(int(teacher.Code)),
		Email:      teacher.Mail,
		Department: teacher.Department,
		Title:      teacher.Title,
		Pinyin:     teacher.Pinyin,
		PinyinAbbr: teacher.PinyinAbbr,
		Biography:  teacher.Biography,
		ProfileURL: teacher.ProfileUrl,
		Picture:    teacher.HeadImage,
	}
}
func SaveTeacher(teachers []LoadedTeacher, db *gorm.DB) {
	if db == nil {
		// HINT:DEBUG only
		for _, teacher := range teachers[:10] {
			jsonData, _ := json.Marshal(Teacher2PO(teacher))
			println("saved: " + string(jsonData))
		}
		return
	}
	teacherPOs := make([]po.TeacherPO, 0)
	for _, teacher := range teachers {
		teacherPOs = append(teacherPOs, Teacher2PO(teacher))
	}
	db.CreateInBatches(teacherPOs, 100)

}
func LoadTeacherProfile2DB(from string, db *gorm.DB) {
	teachers := LoadTeacherProfiles(from)
	SaveTeacher(teachers, db)
}

func LoadTeacherProfiles(from string) []LoadedTeacher {
	file, err := os.Open(from)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	data, _ := os.ReadFile(from)
	var teachers []LoadedTeacher
	err = json.Unmarshal(data, &teachers)
	if err != nil {
		log.Fatalf("In read teacher profile from file %#v:%#v", from, err)
	}
	return teachers
}
