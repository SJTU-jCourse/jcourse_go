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

func main() {
	_ = godotenv.Load()
	dal.InitDBClient()
	db = dal.GetDBClient()

	// import first level
	queryAllBaseCourse()
	readCSV("./data/base_course.csv", importBaseCourse)
	queryAllTeacher()
	readCSV("./data/teacher.csv", importTeacher)

	// refresh
	queryAllBaseCourse()
	queryAllTeacher()

	// import second level
	queryAllCourse()
	readCSV("./data/course.csv", importCourse)

	// refresh
	queryAllCourse()
	queryAllOfferedCourse()

	// import
	readCSV("./data/offered_course.csv", importOfferedCourse)

	// refresh
	queryAllOfferedCourse()
	queryAllOfferedCourseCategory()
	queryAllOfferedCourseTeacherGroup()

	// import
	readCSV("./data/offered_course_category.csv", importOfferedCourseCategory)
	readCSV("./data/offered_course_teacher_group.csv", importOfferedCourseTeacherGroup)
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
		offeredCourseKeyMap[makeOfferedCourseKey(course.Code, teacher.Name, offeredCourse.Semester)] = offeredCourse
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
		offeredCourseTeacherKeyMap[makeOfferedCourseTeacherKey(course.Code, teacher.Name, offeredCourseTeacher.Semester)] = offeredCourseTeacher
	}
	return
}

func makeCourseCategoryKey(courseCode string, teacherName string, category string) string {
	return fmt.Sprintf("%s:%s:%s:%s", courseCode, teacherName, category)
}

func queryAllOfferedCourseCategory() {
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
		teacher, ok := teacherIDMap[uint(courseCategory.MainTeacherID)]
		if !ok {
			continue
		}
		courseCategoryMap[makeCourseCategoryKey(course.Code, teacher.Name, courseCategory.Category)] = courseCategory
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

func importCourse(data [][]string) {
	newCourses := make([]po.CoursePO, 0)
	for _, line := range data {
		courseCode := line[0]
		teacherCode := line[1]
		baseCourse, ok := baseCourseKeyMap[makeBaseCourseKey(courseCode)]
		if !ok {
			continue
		}
		teacher, ok := teacherKeyMap[makeTeacherKey(teacherCode)]
		if !ok {
			continue
		}
		if _, exists := courseKeyMap[makeCourseKey(courseCode, teacher.Name)]; exists {
			continue
		}
		course := po.CoursePO{
			Code:            baseCourse.Code,
			Name:            baseCourse.Name,
			Credit:          baseCourse.Credit,
			BaseCourseID:    int64(baseCourse.ID),
			MainTeacherID:   int64(teacher.ID),
			MainTeacherName: teacher.Name,
		}
		newCourses = append(newCourses, course)
	}
	println("new course length:", len(newCourses))
	db.Model(&po.CoursePO{}).CreateInBatches(&newCourses, 100)
}

func importOfferedCourse(data [][]string) {
	newOfferedCourses := make([]po.OfferedCoursePO, 0)
	for _, line := range data {
		courseCode := line[0]
		teacherCode := line[1]
		semester := line[2]
		teacher, ok := teacherKeyMap[makeTeacherKey(teacherCode)]
		if !ok {
			continue
		}

		course, ok := courseKeyMap[makeCourseKey(courseCode, teacher.Name)]
		if !ok {
			continue
		}

		if _, exists := offeredCourseKeyMap[makeOfferedCourseKey(courseCode, teacher.Name, semester)]; exists {
			continue
		}

		offeredCourse := po.OfferedCoursePO{
			CourseID:      int64(course.ID),
			BaseCourseID:  course.BaseCourseID,
			MainTeacherID: int64(teacher.ID),
			Semester:      semester,
			Department:    line[3],
			Grade:         line[4],
			Language:      line[5],
		}
		newOfferedCourses = append(newOfferedCourses, offeredCourse)
	}
	println("new offered course length:", len(newOfferedCourses))
	db.Model(&po.OfferedCoursePO{}).CreateInBatches(&newOfferedCourses, 100)
}

func importOfferedCourseCategory(data [][]string) {
	newCourseCategories := make([]po.CourseCategoryPO, 0)
	for _, line := range data {
		courseCode := line[0]
		teacherCode := line[1]
		semester := line[2]
		category := line[3]
		teacher, ok := teacherKeyMap[makeTeacherKey(teacherCode)]
		if !ok {
			continue
		}

		offeredCourse, ok := offeredCourseKeyMap[makeOfferedCourseKey(courseCode, teacher.Name, semester)]
		if !ok {
			continue
		}

		if _, exists := courseCategoryMap[makeCourseCategoryKey(courseCode, teacher.Name, category)]; exists {
			continue
		}

		offeredCourseCategory := po.CourseCategoryPO{
			CourseID:      offeredCourse.CourseID,
			BaseCourseID:  offeredCourse.BaseCourseID,
			MainTeacherID: int64(teacher.ID),
			Category:      category,
		}
		newCourseCategories = append(newCourseCategories, offeredCourseCategory)
	}
	println("new course category length:", len(newCourseCategories))
	db.Model(&po.CourseCategoryPO{}).CreateInBatches(&newCourseCategories, 100)
}

func importOfferedCourseTeacherGroup(data [][]string) {
	newOfferedCourseTeachers := make([]po.OfferedCourseTeacherPO, 0)
	for _, line := range data {
		courseCode := line[0]
		mainTeacherCode := line[1]
		semester := line[2]

		mainTeacher, ok := teacherKeyMap[makeTeacherKey(mainTeacherCode)]
		if !ok {
			continue
		}

		offeredCourse, ok := offeredCourseKeyMap[makeOfferedCourseKey(courseCode, mainTeacher.Name, semester)]
		if !ok {
			continue
		}

		thisTeacherCode := line[3]
		thisTeacher, ok := teacherKeyMap[makeTeacherKey(thisTeacherCode)]
		if !ok {
			continue
		}

		if _, exists := offeredCourseTeacherKeyMap[makeOfferedCourseTeacherKey(courseCode, thisTeacher.Name, semester)]; exists {
			continue
		}

		offeredCourseTeacher := po.OfferedCourseTeacherPO{
			CourseID:        offeredCourse.CourseID,
			BaseCourseID:    offeredCourse.BaseCourseID,
			MainTeacherID:   int64(mainTeacher.ID),
			OfferedCourseID: int64(offeredCourse.ID),
			TeacherID:       int64(thisTeacher.ID),
			TeacherName:     thisTeacher.Name,
			Semester:        semester,
		}
		newOfferedCourseTeachers = append(newOfferedCourseTeachers, offeredCourseTeacher)
	}
	println("new offered course teacher length:", len(newOfferedCourseTeachers))
	db.Model(&po.OfferedCourseTeacherPO{}).CreateInBatches(&newOfferedCourseTeachers, 100)
}
