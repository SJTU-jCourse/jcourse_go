package seleniumget

import (
	"bufio"
	"encoding/json"
	"jcourse_go/model/po"
	"os"
	"strconv"
	"strings"

	"gorm.io/gorm"
)
type LoadedCourse struct {
	Code string
	Name string
	Credit float32
	Department string
	SuggestYear int
	SuggestSemester int
}
type LoadedTrainingPlan struct {
	Code string
	Name string
	TotalYear int
	MinPoints float64
	MajorClass string
	Department string
	EntryYear int
	Courses []LoadedCourse
}
func Line2Course(line string) LoadedCourse {
	// AD102,素描（1）,4.0,2018-2019,1,
	meta := strings.Split(line, ",")
	if len(meta) != 6 {
		panic("Invalid line2course: " + line)
	}
	credit, err := strconv.ParseFloat(meta[2], 32)
	if err != nil {
		credit = 0.0
	}
	suggest_year, err := strconv.ParseInt(meta[3][0:4], 10, 32)
	if err != nil {
		suggest_year = 0
	}
	suggest_semester, err := strconv.Atoi(meta[4])
	if err != nil {
		suggest_semester = 0
	}
	return LoadedCourse{
		Code: meta[0],
		Name: meta[1],
		Credit: float32(credit),
		SuggestYear: int(suggest_year),
		SuggestSemester: int(suggest_semester),
		Department: meta[5],
	}
}
func Lines2TrainingPlan(lines []string) LoadedTrainingPlan{
	plan := LoadedTrainingPlan{
		Courses: make([]LoadedCourse, 0),
	}
	// 13050243,14761201813050243,视觉传达设计,13050243,4,28.5,艺术学,未知院系,2018
	metaInfo := strings.Split(lines[0], ",")
	if len(metaInfo) != 9 {
		panic("Invalid line2trainingplan: " + lines[0])
	}
	plan.Name = metaInfo[0]
	plan.Code = metaInfo[1]
	plan.TotalYear, _ = strconv.Atoi(metaInfo[4])
	plan.MinPoints, _ = strconv.ParseFloat(metaInfo[5], 64)
	plan.MajorClass = metaInfo[6]
	plan.Department = metaInfo[7]
	plan.EntryYear, _ = strconv.Atoi(metaInfo[8])
	for _, line := range lines[1:] {
		plan.Courses = append(plan.Courses, Line2Course(line))
	}
	return plan
}
func TrainingPlan2PO(plan LoadedTrainingPlan) po.TrainingPlanPO{
	return po.TrainingPlanPO{
		Degree: "",
		Major: plan.Name,
		Department: plan.Department,
		EntryYear: strconv.Itoa(plan.EntryYear),
		TotalYear: plan.TotalYear,
		MinPoints: plan.MinPoints,
		MajorCode: plan.Code,
		MajorClass: plan.MajorClass,
	}
	// TODO: change po
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
		db.Save(TrainingPlan2PO(plan))
	}
}
func LoadTrainingPlan2DB(from string, db *gorm.DB) {
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
		println(text)
		if text == "" {
			if len(lines) > 0{
				allTrainingPlans = append(allTrainingPlans, Lines2TrainingPlan(lines))
				lines = lines[:0]
			}
			continue
		}
		lines = append(lines, text)
	}
	// HINT: debug only
	SaveTrainingPlan(allTrainingPlans, db)
}
type LoadedTeacher struct{
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
	Name string `json:"name"`
	Code string `json:"code"`
	Department string `json:"department"`
	Title string `json:"title"`
	Pinyin string `json:"pinyin"`
	PinyinAbbr string `json:"pinyin_abbr"`
	ProfileUrl string `json:"profile_url"`
	HeadImage string `json:"head_image"`
	Mail string `json:"mail"`
	ProfileDescription string `json:"profile_description"`
}
func Teacher2PO(teacher LoadedTeacher) po.TeacherPO{
	return po.TeacherPO{
		Name: teacher.Name,
		Code: teacher.Code,
		Email: teacher.Mail,
		Department: teacher.Department,
		Title: teacher.Title,
		Pinyin: teacher.Pinyin,
		PinyinAbbr: teacher.PinyinAbbr,
		// TODO: add more
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
	for _, teacher := range teachers {
		db.Save(Teacher2PO(teacher))
	}
}
func LoadTeacherProfile2DB(from string, db *gorm.DB){
	file, err := os.Open(from)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	data, _ := os.ReadFile(from)
	var teachers []LoadedTeacher
	json.Unmarshal(data, &teachers)
	SaveTeacher(teachers, db)
}