package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/mozillazg/go-pinyin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"jcourse_go/internal/config"
	"jcourse_go/internal/infrastructure/dal"
	"jcourse_go/internal/infrastructure/entity"
)

const Semester = "2025-2026-1"

var (
	db *gorm.DB
)

type teacherGroupItem struct {
	Code       string
	Name       string
	Department string
	Title      string
}

func (t teacherGroupItem) uniqueKey() string {
	return t.Code
}

var teacherItemPattern = regexp.MustCompile(`([^/]+)/([^/]+)/([^[]+)\[([^]]+)\]`)

func parseTeacherGroup(s string) []teacherGroupItem {
	match := teacherItemPattern.FindAllStringSubmatch(s, -1)
	result := make([]teacherGroupItem, 0, len(match))
	for _, m := range match {
		result = append(result, teacherGroupItem{
			Code:       strings.Trim(m[1], ";"),
			Name:       m[2],
			Title:      m[3],
			Department: m[4],
		})
	}
	return result
}

type mainTeacherItem struct {
	Code string
	Name string
}

func parseMainTeacher(s string) mainTeacherItem {
	slice := strings.Split(s, "|")
	if len(slice) < 2 {
		return mainTeacherItem{}
	}
	return mainTeacherItem{
		Code: slice[0],
		Name: slice[1],
	}
}

type teachCourse struct {
	Code         string             // 课程号
	Name         string             // 课程名称
	TeachHour    int                // 学时
	TeacherGroup []teacherGroupItem // 合上教师
	MainTeacher  mainTeacherItem    // 任课教师
	Department   string             // 开课院系
	AttendTime   string             // 课程安排
	ClassName    string             // 教学班名称
	EnrollCount  int                // 选课人数
	Credit       float64            // 学分
	ClassRoom    string             // 教室
	Language     string             // 授课语言
	IsGeneral    bool               // 是否通识课
	Categories   []string           // 通识课归属模块
	Grade        []string           // 年级
}

func (c teachCourse) uniqueKey() string {
	return fmt.Sprintf("%s:%s", c.Code, c.MainTeacher.Code)
}

func (c teachCourse) toCurriculum() entity.Curriculum {
	return entity.Curriculum{
		Code:   c.Code,
		Name:   c.Name,
		Credit: c.Credit,
	}
}

func (c teachCourse) getTeachers() []entity.Teacher {
	teachers := make([]entity.Teacher, 0, len(c.TeacherGroup))
	for _, teacher := range c.TeacherGroup {
		teachers = append(teachers, entity.Teacher{
			Code:       teacher.Code,
			Name:       teacher.Name,
			Pinyin:     generatePinyin(teacher.Name),
			PinyinAbbr: generatePinyinAbbr(teacher.Name),
			Department: teacher.Department,
			Title:      teacher.Title,
		})
	}
	return teachers
}

func parseLine(line []string) (teachCourse, bool) {
	teachHour, _ := strconv.Atoi(line[2])
	teacherGroup := parseTeacherGroup(line[3])
	mainTeacher := parseMainTeacher(line[4])
	if mainTeacher.Code == "" {
		if len(teacherGroup) > 0 {
			mainTeacher = mainTeacherItem{teacherGroup[0].Code, teacherGroup[0].Name}
		} else {
			return teachCourse{}, false
		}
	}
	categories := make([]string, 0)
	if line[13] != "" {
		categories = strings.Split(line[13], ",")
	}
	credit, _ := strconv.ParseFloat(line[9], 64)
	return teachCourse{
		Code:         line[0],
		Name:         line[1],
		TeachHour:    teachHour,
		TeacherGroup: teacherGroup,
		MainTeacher:  mainTeacher,
		Department:   line[5],
		AttendTime:   line[6],
		ClassName:    line[7],
		EnrollCount:  0,
		Credit:       credit,
		ClassRoom:    line[10],
		Language:     line[11],
		IsGeneral:    line[12] == "是",
		Categories:   categories,
		Grade:        strings.Split(line[14], ","),
	}, true
}

func readRawCSV(filename string) []teachCourse {
	fs, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer fs.Close()

	reader := csv.NewReader(fs)
	lines, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	courses := make([]teachCourse, 0, len(lines))
	for _, line := range lines[1:] {
		course, ok := parseLine(line)
		if !ok {
			continue
		}
		courses = append(courses, course)
	}
	return courses
}

func main() {
	var err error
	conf := config.InitConfig("./config/config.yaml")
	db, err = dal.NewPostgresSQL(conf.DB)
	if err != nil {
		panic(err)
	}
	rawCourses := readRawCSV(fmt.Sprintf("./data/%s.csv", Semester))
	rawCourses = removeDuplicatedCourses(rawCourses)

	// 1. 写入没有任何依赖项的两个模型

	curriculumMaps := make(map[string]entity.Curriculum)
	curriculums := make([]entity.Curriculum, 0)
	for _, course := range rawCourses {
		if _, ok := curriculumMaps[course.Code]; ok {
			continue
		}
		curriculum := course.toCurriculum()
		curriculumMaps[course.Code] = curriculum
		curriculums = append(curriculums, curriculum)
	}

	teachers := make([]entity.Teacher, 0)
	teacherMap := make(map[string]entity.Teacher)
	for _, course := range rawCourses {
		for _, teacher := range course.getTeachers() {
			if _, ok := teacherMap[teacher.Code]; ok {
				continue
			}
			teacherMap[teacher.Code] = teacher
			teachers = append(teachers, teacher)
		}
	}

	if err = db.Model(&entity.Curriculum{}).
		Clauses(clause.OnConflict{UpdateAll: true}).
		CreateInBatches(&curriculums, 100).Error; err != nil {
		log.Fatal(err)
	}
	if err = db.Model(&entity.Teacher{}).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "code"}}, // 以 code 为唯一键
		DoUpdates: clause.AssignmentColumns([]string{"name", "department", "title"}),
	}).CreateInBatches(&teachers, 100).Error; err != nil {
		log.Fatal(err)
	}

	// 2. 查询现有的所有 teacher，获取 teacher_id
	clear(teachers)
	if err = db.Model(&entity.Teacher{}).Find(&teachers).Error; err != nil {
		log.Fatal(err)
	}
	teacherCode2IDMap := make(map[string]entity.Teacher)
	for _, teacher := range teachers {
		teacherCode2IDMap[teacher.Code] = teacher
	}

	// 3. 构建 course 模型

	courses := make([]entity.Course, 0)
	for _, c := range rawCourses {
		t, ok := teacherCode2IDMap[c.MainTeacher.Code]
		if !ok {
			log.Printf("teacher %s not found", c.MainTeacher.Code)
			continue
		}
		course := entity.Course{
			Code:          c.Code,
			Name:          c.Name,
			Credit:        c.Credit,
			MainTeacherID: t.ID,
			Offerings:     nil,
		}
		courses = append(courses, course)
	}

	if err := db.Model(&entity.Course{}).
		Clauses(clause.OnConflict{DoNothing: true}).
		CreateInBatches(&courses, 100).Error; err != nil {
		log.Fatal(err)
	}
	clear(courses)
	if err := db.Model(&entity.Course{}).
		Joins("MainTeacher").
		Find(&courses).Error; err != nil {
		log.Fatal(err)
	}

	courseMaps := make(map[string]entity.Course)
	for _, course := range courses {
		courseMaps[fmt.Sprintf("%s:%s", course.Code, course.MainTeacher.Code)] = course
	}

	// 写入 course_offering 模型
	courseOfferings := make([]entity.CourseOffering, 0)
	for _, c := range rawCourses {
		course, ok := courseMaps[fmt.Sprintf("%s:%s", c.Code, c.MainTeacher.Code)]
		if !ok {
			log.Printf("course %s not found", c.Code)
			continue
		}
		co := entity.CourseOffering{
			CourseID:      course.ID,
			MainTeacherID: course.MainTeacherID,
			Semester:      Semester,
			Department:    c.Department,
			Language:      c.Language,
			Categories:    make([]entity.CourseOfferingCategory, 0),
			TeacherGroup:  make([]entity.CourseOfferingTeacher, 0),
		}
		for _, teacher := range c.getTeachers() {
			t, ok := teacherCode2IDMap[teacher.Code]
			if !ok {
				log.Printf("teacher %s not found", teacher.Code)
				continue
			}
			co.TeacherGroup = append(co.TeacherGroup, entity.CourseOfferingTeacher{
				TeacherID: t.ID,
				CourseID:  course.ID,
			})
		}
		for _, category := range c.Categories {
			co.Categories = append(co.Categories, entity.CourseOfferingCategory{
				Category: category,
				CourseID: course.ID,
			})
		}
		courseOfferings = append(courseOfferings, co)
	}

	if err := db.Model(&entity.CourseOffering{}).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "course_id"}, {Name: "semester"}},
			UpdateAll: true,
		}).
		CreateInBatches(&courseOfferings, 100).Error; err != nil {
		log.Fatal(err)
	}
}

func removeDuplicatedCourses(courses []teachCourse) []teachCourse {
	m := make(map[string]teachCourse)
	result := make([]teachCourse, 0, len(m))

	for _, course := range courses {
		if _, ok := m[course.uniqueKey()]; ok {
			continue
		}
		m[course.uniqueKey()] = course
		result = append(result, course)
	}
	return result
}

func generatePinyin(name string) string {
	result := pinyin.LazyPinyin(name, pinyin.NewArgs())
	return strings.Join(result, "")
}

func generatePinyinAbbr(name string) string {
	result := pinyin.LazyPinyin(name, pinyin.Args{Style: pinyin.FirstLetter})
	return strings.Join(result, "")
}
