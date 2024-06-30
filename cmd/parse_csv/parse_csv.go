package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"jcourse_go/model/domain"
)

const Semester = "2024-2025-1"

func parseBaseCourseFromLine(line []string) domain.BaseCourse {
	credit, _ := strconv.ParseFloat(line[9], 32)
	return domain.BaseCourse{
		Code:   line[0],
		Name:   line[1],
		Credit: float32(credit),
	}
}

func parseCourseFromLine(line []string) domain.Course {
	return domain.Course{
		BaseCourse:  parseBaseCourseFromLine(line),
		MainTeacher: parseTeacherFromString(line[4]),
	}
}

func parseTeacherFromString(str string) domain.Teacher {
	l := strings.Split(str, "/")
	if len(l) == 3 {
		s := strings.Split(l[2], "[")
		dept, _ := strings.CutSuffix(s[1], "]")
		return domain.Teacher{
			Name:       l[1],
			Department: dept,
			Title:      s[0],
			Code:       l[0],
		}
	}
	l = strings.Split(str, "|")
	return domain.Teacher{
		Name: l[1],
		Code: l[0],
	}
}

func parseTeacherGroupFromString(str string) []domain.Teacher {
	replaced := strings.ReplaceAll(str, "THIERRY; Fine; VAN CHUNG", "THIERRY, Fine, VAN CHUNG")
	teacherGroup := make([]domain.Teacher, 0)
	l := strings.Split(replaced, ";")
	for _, v := range l {
		teacherGroup = append(teacherGroup, parseTeacherFromString(v))
	}
	return teacherGroup
}

func parseOfferedCourseFromLine(line []string) domain.OfferedCourse {
	course := parseCourseFromLine(line)
	mainTeacher := parseTeacherFromString(line[4])
	teacherGroup := parseTeacherGroupFromString(line[3])
	return domain.OfferedCourse{
		Course:       course,
		MainTeacher:  mainTeacher,
		TeacherGroup: teacherGroup,
		Semester:     Semester,
		Department:   line[5],
		Categories:   strings.Split(line[13], ","),
	}
}

func makeOfferedCourseKey(offeredCourse domain.OfferedCourse) string {
	return fmt.Sprintf("%s:%s", offeredCourse.Course.BaseCourse.Code, offeredCourse.MainTeacher.Code)
}

func makeCourseKey(course domain.Course) string {
	return fmt.Sprintf("%s:%s", course.BaseCourse.Code, course.MainTeacher.Code)
}

func main() {
	filename := fmt.Sprintf("./data/%s.csv", Semester)
	fs, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer fs.Close()

	r := csv.NewReader(fs)
	content, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	baseCourseMap := make(map[string]domain.BaseCourse)
	courseMap := make(map[string]domain.Course)
	offeredCourseMap := make(map[string]domain.OfferedCourse)
	teacherMap := make(map[string]domain.Teacher)

	for _, line := range content[1:] {
		baseCourse := parseBaseCourseFromLine(line)
		course := parseCourseFromLine(line)
		offeredCourse := parseOfferedCourseFromLine(line)

		baseCourseMap[baseCourse.Code] = baseCourse
		courseMap[makeCourseKey(course)] = course
		offeredCourseMap[makeOfferedCourseKey(offeredCourse)] = offeredCourse
		teacherMap[offeredCourse.MainTeacher.Code] = offeredCourse.MainTeacher
		for _, teacher := range offeredCourse.TeacherGroup {
			teacherMap[teacher.Code] = teacher
		}
	}

	println(len(offeredCourseMap), len(teacherMap), len(baseCourseMap), len(courseMap))

	outputBaseCourse("./data/base_course.csv", baseCourseMap)
	outputTeacher("./data/teacher.csv", teacherMap)
	outputCourse("./data/course.csv", courseMap)
	outputOfferedCourse("./data/offered_course.csv", offeredCourseMap)
}

func outputBaseCourse(filename string, courses map[string]domain.BaseCourse) {
	fs, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	defer fs.Close()
	if err != nil {
		return
	}
	writer := csv.NewWriter(fs)
	for _, course := range courses {
		err = writer.Write([]string{course.Code, course.Name, strconv.FormatFloat(float64(course.Credit), 'g', 2, 32)})
		if err != nil {
			continue
		}
	}
	writer.Flush()
}

func outputCourse(filename string, courses map[string]domain.Course) {
	fs, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	defer fs.Close()
	if err != nil {
		return
	}
	writer := csv.NewWriter(fs)
	for _, course := range courses {
		err = writer.Write([]string{course.Code, course.MainTeacher.Code})
		if err != nil {
			continue
		}
	}
	writer.Flush()
}

func outputTeacher(filename string, teachers map[string]domain.Teacher) {
	fs, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	defer fs.Close()
	if err != nil {
		return
	}
	writer := csv.NewWriter(fs)
	for _, teacher := range teachers {
		err = writer.Write([]string{teacher.Code, teacher.Name, teacher.Department, teacher.Title})
		if err != nil {
			continue
		}
	}
	writer.Flush()
}

func outputOfferedCourse(filename string, offeredCourse map[string]domain.OfferedCourse) {
	fs, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	defer fs.Close()
	if err != nil {
		return
	}
	writer := csv.NewWriter(fs)
	for _, course := range offeredCourse {
		err = writer.Write([]string{course.Code, course.MainTeacher.Code, course.Department, course.Semester, course.Grade, course.Language, course.Location})
		if err != nil {
			continue
		}
	}
	writer.Flush()
}
