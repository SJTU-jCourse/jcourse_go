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
		Credit: credit,
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
		Language:     line[11],
		Grade:        strings.Split(line[14], ","),
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

	outputCSV("./data/base_course.csv", outputBaseCourse, baseCourseMap)
	outputCSV("./data/teacher.csv", outputTeacher, teacherMap)
	outputCSV("./data/course.csv", outputCourse, courseMap)
	outputCSV("./data/offered_course.csv", outputOfferedCourse, offeredCourseMap)
	outputCSV("./data/offered_course_category.csv", outputOfferedCourseCategory, offeredCourseMap)
	outputCSV("./data/offered_course_teacher_group.csv", outputOfferedCourseTeacherGroup, offeredCourseMap)
}

func outputCSV(filename string, writeFunc func(writer *csv.Writer, data any), data any) {
	fs, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	defer fs.Close()
	if err != nil {
		return
	}
	writer := csv.NewWriter(fs)
	writeFunc(writer, data)
	writer.Flush()
}

func outputBaseCourse(writer *csv.Writer, data any) {
	courses := data.(map[string]domain.BaseCourse)
	for _, course := range courses {
		err := writer.Write([]string{course.Code, course.Name, strconv.FormatFloat(float64(course.Credit), 'g', 2, 32)})
		if err != nil {
			continue
		}
	}

}

func outputCourse(writer *csv.Writer, data any) {
	courses := data.(map[string]domain.Course)
	for _, course := range courses {
		err := writer.Write([]string{course.Code, course.MainTeacher.Code})
		if err != nil {
			continue
		}
	}
}

func outputTeacher(writer *csv.Writer, data any) {
	teachers := data.(map[string]domain.Teacher)
	for _, teacher := range teachers {
		err := writer.Write([]string{teacher.Code, teacher.Name, teacher.Department, teacher.Title})
		if err != nil {
			continue
		}
	}

}

func outputOfferedCourse(writer *csv.Writer, data any) {
	offeredCourse := data.(map[string]domain.OfferedCourse)
	for _, course := range offeredCourse {
		err := writer.Write([]string{course.Code, course.MainTeacher.Code, course.Semester, course.Department, strings.Join(course.Grade, ","), course.Language})
		if err != nil {
			continue
		}
	}

}

func outputOfferedCourseCategory(writer *csv.Writer, data any) {
	offeredCourse := data.(map[string]domain.OfferedCourse)
	for _, course := range offeredCourse {
		for _, category := range course.Categories {
			if category == "" {
				continue
			}
			err := writer.Write([]string{course.Code, course.MainTeacher.Code, course.Semester, category})
			if err != nil {
				continue
			}
		}
	}
}

func outputOfferedCourseTeacherGroup(writer *csv.Writer, data any) {
	offeredCourse := data.(map[string]domain.OfferedCourse)
	for _, course := range offeredCourse {
		for _, teacher := range course.TeacherGroup {
			err := writer.Write([]string{course.Code, course.MainTeacher.Code, course.Semester, teacher.Code})
			if err != nil {
				continue
			}
		}
	}
}
