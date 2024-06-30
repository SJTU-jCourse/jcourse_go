package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	pinyin2 "github.com/mozillazg/go-pinyin"
	"gorm.io/gorm"
	"jcourse_go/dal"
	"jcourse_go/model/po"
)

const Semester = "2024-2025-1"

var (
	db                         *gorm.DB
	baseCourseKeyMap           = make(map[string]po.BaseCoursePO)
	baseCourseIDMap            = make(map[uint]po.BaseCoursePO)
	courseKeyMap               = make(map[string]po.CoursePO)
	courseIDMap                = make(map[uint]po.CoursePO)
	teacherKeyMap              = make(map[string]po.TeacherPO)
	teacherIDMap               = make(map[uint]po.TeacherPO)
	offeredCourseKeyMap        = make(map[string]po.OfferedCoursePO)
	offeredCourseIDMap         = make(map[uint]po.OfferedCoursePO)
	offeredCourseCategoryMap   = make(map[string]po.OfferedCourseCategoryPO)
	offeredCourseTeacherKeyMap = make(map[string]po.OfferedCourseTeacherPO)
)

func main() {
	_ = godotenv.Load()
	dal.InitDBClient()
	db = dal.GetDBClient()

	queryAllBaseCourse()
	readCSV("./data/base_course.csv", importBaseCourse)
	queryAllTeacher()
	readCSV("./data/teacher.csv", importTeacher)
	// queryAllCourse()
}

func makeBaseCourseKey(courseCode string) string {
	return courseCode
}
func queryAllBaseCourse() {
	baseCourses := make([]po.BaseCoursePO, 0)
	result := db.Model(&po.BaseCoursePO{}).Find(&baseCourses)
	if result.Error != nil {
		return
	}
	for _, baseCourse := range baseCourses {
		baseCourseKeyMap[makeBaseCourseKey(baseCourse.Code)] = baseCourse
		baseCourseIDMap[baseCourse.ID] = baseCourse
	}
}

func makeCourseKey(courseCode, mainTeacherName string) string {
	return fmt.Sprintf("%s:%s", courseCode, mainTeacherName)
}

func queryAllCourse() {
	courses := make([]po.CoursePO, 0)
	result := db.Model(&po.CoursePO{}).Find(&courses)
	if result.Error != nil {
		return
	}
	for _, course := range courses {
		courseKeyMap[makeCourseKey(course.Code, course.MainTeacherName)] = course
		courseIDMap[course.ID] = course
	}
	return
}

func makeTeacherKey(teacherCode string) string {
	return teacherCode
}

func queryAllTeacher() {
	teachers := make([]po.TeacherPO, 0)

	result := db.Model(&po.TeacherPO{}).Find(&teachers)
	if result.Error != nil {
		return
	}
	for _, teacher := range teachers {
		teacherKeyMap[makeTeacherKey(teacher.Code)] = teacher
		teacherIDMap[teacher.ID] = teacher
	}
	return
}

func makeOfferedCourseKey(courseCode string, mainTeacherName string, semester string) string {
	return fmt.Sprintf("%s:%s:%s", courseCode, mainTeacherName, semester)
}

func queryAllOfferedCourse() {
	offeredCourses := make([]po.OfferedCoursePO, 0)

	result := db.Model(&po.OfferedCoursePO{}).Find(&offeredCourses)
	if result.Error != nil {
		return
	}
	for _, offeredCourse := range offeredCourses {
		offeredCourseIDMap[offeredCourse.ID] = offeredCourse
		course, ok := courseIDMap[uint(offeredCourse.CourseID)]
		if !ok {
			continue
		}
		teacher, ok := teacherIDMap[uint(offeredCourse.MainTeacherID)]
		if !ok {
			continue
		}
		offeredCourseKeyMap[makeOfferedCourseKey(course.Code, teacher.Name, Semester)] = offeredCourse
	}
	return
}

func makeOfferedCourseTeacherKey(courseCode string, teacherName string, semester string) string {
	return fmt.Sprintf("%s:%s:%s", courseCode, teacherName, semester)
}

func queryAllOfferedCourseTeacherGroup() {
	offeredCourseTeachers := make([]po.OfferedCourseTeacherPO, 0)

	result := db.Model(&po.OfferedCourseTeacherPO{}).Find(&offeredCourseTeachers)
	if result.Error != nil {
		return
	}
	for _, offeredCourseTeacher := range offeredCourseTeachers {
		course, ok := courseIDMap[uint(offeredCourseTeacher.CourseID)]
		if !ok {
			continue
		}
		teacher, ok := teacherIDMap[uint(offeredCourseTeacher.TeacherID)]
		if !ok {
			continue
		}
		offeredCourseTeacherKeyMap[makeOfferedCourseTeacherKey(course.Code, teacher.Name, Semester)] = offeredCourseTeacher
	}
	return
}

func makeOfferedCourseCategoryKey(courseCode string, teacherName string, category string, semester string) string {
	return fmt.Sprintf("%s:%s:%s:%s", courseCode, teacherName, category, semester)
}

func queryAllOfferedCourseCategory() {
	offeredCourseCategories := make([]po.OfferedCourseCategoryPO, 0)
	result := db.Model(&po.OfferedCourseCategoryPO{}).Find(&offeredCourseCategories)
	if result.Error != nil {
		return
	}
	for _, offeredCourseCategory := range offeredCourseCategories {
		course, ok := courseIDMap[uint(offeredCourseCategory.CourseID)]
		if !ok {
			continue
		}
		teacher, ok := teacherIDMap[uint(offeredCourseCategory.MainTeacherID)]
		if !ok {
			continue
		}
		offeredCourseCategoryMap[makeOfferedCourseCategoryKey(course.Code, teacher.Name, offeredCourseCategory.Category, Semester)] = offeredCourseCategory
	}
	return
}

func readCSV(filename string, readerFunc func([][]string)) {
	fs, err := os.Open(filename)
	defer fs.Close()
	if err != nil {
		return
	}
	reader := csv.NewReader(fs)
	lines, err := reader.ReadAll()
	if err != nil {
		return
	}
	readerFunc(lines)
}

func importBaseCourse(data [][]string) {
	newBaseCourse := make([]po.BaseCoursePO, 0)
	for _, line := range data {
		code := line[0]
		name := line[1]
		credit, _ := strconv.ParseFloat(line[2], 64)
		if _, exists := baseCourseKeyMap[makeBaseCourseKey(code)]; exists {
			continue
		}
		newBaseCourse = append(newBaseCourse, po.BaseCoursePO{Code: code, Name: name, Credit: credit})
	}
	println("new base course length:", len(newBaseCourse))
	db.Model(&po.BaseCoursePO{}).CreateInBatches(&newBaseCourse, 100)
}

func importTeacher(data [][]string) {
	newTeacher := make([]po.TeacherPO, 0)
	for _, line := range data {
		code := line[0]
		name := line[1]
		department := line[2]
		title := line[3]
		if _, exists := teacherKeyMap[makeTeacherKey(code)]; exists {
			continue
		}
		pinyin := generatePinyin(name)
		pinyinAbbr := generatePinyinAbbr(name)
		newTeacher = append(newTeacher, po.TeacherPO{Code: code, Name: name, Department: department, Title: title, Pinyin: pinyin, PinyinAbbr: pinyinAbbr})
	}
	println("new teacher length:", len(newTeacher))
	db.Model(&po.TeacherPO{}).CreateInBatches(&newTeacher, 100)

}

func generatePinyin(name string) string {
	result := pinyin2.LazyPinyin(name, pinyin2.NewArgs())
	return strings.Join(result, "")
}

func generatePinyinAbbr(name string) string {
	result := pinyin2.LazyPinyin(name, pinyin2.Args{Style: pinyin2.FirstLetter})
	return strings.Join(result, "")
}
