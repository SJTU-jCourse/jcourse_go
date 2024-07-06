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
	"gorm.io/gorm/clause"
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
	courseCategoryMap          = make(map[string]po.CourseCategoryPO)
	offeredCourseKeyMap        = make(map[string]po.OfferedCoursePO)
	offeredCourseIDMap         = make(map[uint]po.OfferedCoursePO)
	offeredCourseTeacherKeyMap = make(map[string]po.OfferedCourseTeacherPO)
)

func initDB() {
	_ = godotenv.Load()
	dal.InitDBClient()
	db = dal.GetDBClient()
}

func readRawCSV(filename string) [][]string {
	fs, err := os.Open(filename)
	defer fs.Close()
	if err != nil {
		panic(err)
	}
	reader := csv.NewReader(fs)
	lines, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	return lines
}

func main() {
	initDB()
	data := readRawCSV(fmt.Sprintf("./data/%s.csv", Semester))

	// init
	queryAllBaseCourse()
	queryAllTeacher()
	queryAllCourse()
	queryAllOfferedCourse()
	queryAllOfferedCourseTeacherGroup()
	queryAllCourseCategory()

	// first import
	importBaseCourse(data)
	importTeacher(data)

	// refresh
	queryAllBaseCourse()
	queryAllTeacher()

	importCourse(data)
	queryAllCourse()

	importCourseCategory(data)

	importOfferedCourse(data)
	queryAllOfferedCourse()

	importOfferedCourseTeacher(data)
}

func importBaseCourse(data [][]string) {
	baseCourses := make([]po.BaseCoursePO, 0)
	for _, line := range data[1:] {
		baseCourses = append(baseCourses, parseBaseCourseFromLine(line))
	}
	db.Model(&po.BaseCoursePO{}).Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(&baseCourses, 100)
}

func importTeacher(data [][]string) {
	teachers := make([]po.TeacherPO, 0)
	teacherSet := make(map[string]bool)
	for _, line := range data[1:] {
		for _, t := range parseTeacherGroupFromLine(line) {
			if _, ok := teacherSet[t.Code]; ok {
				continue
			}
			teachers = append(teachers, t)
			teacherSet[t.Code] = true
		}
	}
	db.Model(&po.TeacherPO{}).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "code"}},
		DoUpdates: clause.AssignmentColumns([]string{"department", "title"}),
	}).CreateInBatches(&teachers, 100)
}

func importCourse(data [][]string) {
	courses := make([]po.CoursePO, 0)
	for _, line := range data[1:] {
		courses = append(courses, parseCourseFromLine(line))
	}
	db.Model(&po.CoursePO{}).Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(&courses, 100)
}

func importOfferedCourse(data [][]string) {
	offeredCourses := make([]po.OfferedCoursePO, 0)
	for _, line := range data[1:] {
		offeredCourses = append(offeredCourses, parseOfferedCourseFromLine(line))
	}
	db.Model(&po.OfferedCoursePO{}).Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(&offeredCourses, 100)
}

func importCourseCategory(data [][]string) {
	categories := make([]po.CourseCategoryPO, 0)
	for _, line := range data[1:] {
		categories = append(categories, parseCourseCategories(line)...)
	}
	db.Model(&po.CourseCategoryPO{}).Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(&categories, 100)
}

func importOfferedCourseTeacher(data [][]string) {
	offeredCourseTeachers := make([]po.OfferedCourseTeacherPO, 0)
	for _, line := range data[1:] {
		offeredCourseTeachers = append(offeredCourseTeachers, parseOfferedCourseTeacherGroup(line)...)
	}
	db.Model(&po.OfferedCourseTeacherPO{}).Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(&offeredCourseTeachers, 100)
}

func parseBaseCourseFromLine(line []string) po.BaseCoursePO {
	credit, _ := strconv.ParseFloat(line[9], 32)
	baseCourse := po.BaseCoursePO{
		Code:   line[0],
		Name:   line[1],
		Credit: credit,
	}
	if baseCourseFromDB, ok := baseCourseKeyMap[makeBaseCourseKey(baseCourse.Code)]; ok {
		baseCourse.ID = baseCourseFromDB.ID
	}
	return baseCourse
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

func parseMainTeacherFromLine(line []string) po.TeacherPO {
	teacherInfo := strings.Split(line[4], "|")
	teacher := po.TeacherPO{
		Name:       teacherInfo[1],
		Code:       teacherInfo[0],
		Pinyin:     generatePinyin(teacherInfo[1]),
		PinyinAbbr: generatePinyinAbbr(teacherInfo[1]),
	}
	if teacherFromDB, ok := teacherKeyMap[makeTeacherKey(teacher.Code)]; ok {
		teacher.Department = teacherFromDB.Department
		teacher.Title = teacherFromDB.Title
		teacher.ID = teacherFromDB.ID
	}
	return teacher
}

func parseSingleTeacherFromLine(teacherInfo string) po.TeacherPO {
	l := strings.Split(teacherInfo, "/")
	s := strings.Split(l[2], "[")
	dept, _ := strings.CutSuffix(s[1], "]")
	teacher := po.TeacherPO{
		Name:       l[1],
		Code:       l[0],
		Department: dept,
		Title:      s[0],
		Pinyin:     generatePinyin(l[1]),
		PinyinAbbr: generatePinyinAbbr(l[1]),
	}
	if teacherFromDB, ok := teacherKeyMap[makeTeacherKey(teacher.Code)]; ok {
		teacher.ID = teacherFromDB.ID
	}
	return teacher
}

func parseTeacherGroupFromLine(line []string) []po.TeacherPO {
	replaced := strings.ReplaceAll(line[3], "THIERRY; Fine; VAN CHUNG", "THIERRY, Fine, VAN CHUNG")
	teacherInfos := strings.Split(replaced, ";")

	teachers := make([]po.TeacherPO, len(teacherInfos))
	for _, teacherInfo := range teacherInfos {
		teachers = append(teachers, parseSingleTeacherFromLine(teacherInfo))
	}
	return teachers
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

func parseCourseFromLine(line []string) po.CoursePO {
	baseCourse := parseBaseCourseFromLine(line)
	mainTeacher := parseMainTeacherFromLine(line)
	course := po.CoursePO{
		Code:            baseCourse.Code,
		Name:            baseCourse.Name,
		Credit:          baseCourse.Credit,
		MainTeacherID:   int64(mainTeacher.ID),
		MainTeacherName: mainTeacher.Name,
	}
	if courseFromDB, ok := courseKeyMap[makeCourseKey(course.Code, mainTeacher.Name)]; ok {
		course.ID = courseFromDB.ID
	}
	return course
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

func parseOfferedCourseFromLine(line []string) po.OfferedCoursePO {
	course := parseCourseFromLine(line)
	mainTeacher := parseMainTeacherFromLine(line)
	offeredCourse := po.OfferedCoursePO{
		CourseID:      int64(course.ID),
		MainTeacherID: int64(mainTeacher.ID),
		Semester:      Semester,
		Department:    line[5],
		Language:      line[11],
		Grade:         line[14],
	}
	if offeredCourseFromDB, ok := offeredCourseKeyMap[makeOfferedCourseKey(int64(course.ID), Semester)]; ok {
		offeredCourse.ID = offeredCourseFromDB.ID
	}
	return offeredCourse
}

func makeOfferedCourseKey(courseID int64, semester string) string {
	return fmt.Sprintf("%d:%s", courseID, semester)
}

func queryAllOfferedCourse() {
	offeredCourses := make([]po.OfferedCoursePO, 0)
	result := db.Model(&po.OfferedCoursePO{}).Find(&offeredCourses)
	if result.Error != nil {
		return
	}
	for _, offeredCourse := range offeredCourses {
		offeredCourseIDMap[offeredCourse.ID] = offeredCourse
		offeredCourseKeyMap[makeOfferedCourseKey(offeredCourse.CourseID, offeredCourse.Semester)] = offeredCourse
	}
	return
}

func parseOfferedCourseTeacherGroup(line []string) []po.OfferedCourseTeacherPO {
	teacherGroup := parseTeacherGroupFromLine(line)
	offeredCourse := parseOfferedCourseFromLine(line)
	teachers := make([]po.OfferedCourseTeacherPO, 0)
	for _, teacher := range teacherGroup {
		teachers = append(teachers, po.OfferedCourseTeacherPO{
			CourseID:        offeredCourse.CourseID,
			OfferedCourseID: int64(offeredCourse.ID),
			MainTeacherID:   offeredCourse.MainTeacherID,
			TeacherID:       int64(teacher.ID),
			TeacherName:     teacher.Name,
		})
	}
	return teachers
}

func makeOfferedCourseTeacherKey(offeredCourseID int64, teacherID int64) string {
	return fmt.Sprintf("%d:%d", offeredCourseID, teacherID)
}

func queryAllOfferedCourseTeacherGroup() {
	offeredCourseTeachers := make([]po.OfferedCourseTeacherPO, 0)
	result := db.Model(&po.OfferedCourseTeacherPO{}).Find(&offeredCourseTeachers)
	if result.Error != nil {
		return
	}
	for _, offeredCourseTeacher := range offeredCourseTeachers {
		offeredCourseTeacherKeyMap[makeOfferedCourseTeacherKey(offeredCourseTeacher.OfferedCourseID, offeredCourseTeacher.TeacherID)] = offeredCourseTeacher
	}
	return
}

func parseCourseCategories(line []string) []po.CourseCategoryPO {
	course := parseCourseFromLine(line)
	categories := strings.Split(line[13], ",")
	courseCategories := make([]po.CourseCategoryPO, 0)
	for _, category := range categories {
		if category == "" {
			continue
		}
		courseCategories = append(courseCategories, po.CourseCategoryPO{
			CourseID: int64(course.ID),
			Category: category,
		})
	}
	return courseCategories
}
func makeCourseCategoryKey(courseID int64, category string) string {
	return fmt.Sprintf("%d:%s", courseID, category)
}

func queryAllCourseCategory() {
	courseCategories := make([]po.CourseCategoryPO, 0)
	result := db.Model(&po.CourseCategoryPO{}).Find(&courseCategories)
	if result.Error != nil {
		return
	}
	for _, courseCategory := range courseCategories {
		course, ok := courseIDMap[uint(courseCategory.CourseID)]
		if !ok {
			continue
		}
		courseCategoryMap[makeCourseCategoryKey(int64(course.ID), courseCategory.Category)] = courseCategory
	}
	return
}

func generatePinyin(name string) string {
	result := pinyin2.LazyPinyin(name, pinyin2.NewArgs())
	return strings.Join(result, "")
}

func generatePinyinAbbr(name string) string {
	result := pinyin2.LazyPinyin(name, pinyin2.Args{Style: pinyin2.FirstLetter})
	return strings.Join(result, "")
}
