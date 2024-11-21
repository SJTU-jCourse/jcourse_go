package mock

import (
	"jcourse_go/model/po"
	"strconv"

	"gorm.io/gorm/clause"

	"syreclabs.com/go/faker"
)

// 课程号,课程名称,学时,合上教师,任课教师,开课院系,课程安排,教学班名称,选课人数,学分,教室,授课语言,是否通识课,通识课归属模块,年级
// CS1010,计算机基础,48,08123/李华/副教授[计算机科学与技术学院],08123|李华,计算机科学与技术学院,星期一第1-2节{1-16周};星期三第3-4节{1-16周}|西教101;西教101|Monday Section 1-2{Week 1-16};Wednesday Section 3-4{Week 1-16}|西教101;西教101,(2024-2025-1)-CS1010-01,30,3,西教101,中文,否,,2023
// ME2015,机械设计基础,64,08124/王强/副教授[机械工程学院],08124|王强,机械工程学院,星期二第5-6节{1-16周};星期四第3-4节{1-16周}|南教楼205;南教楼205|Tuesday Section 5-6{Week 1-16};Thursday Section 3-4{Week 1-16}|南教楼205;南教楼205,(2024-2025-1)-ME2015-01,25,4,南教楼205,中文,否,,2022
func MockBaseCourses(gen MockDBGenerator, n int) ([]po.BaseCoursePO, error) {
	baseCourses := make([]po.BaseCoursePO, n)
	codes, err := GenerateUniqueSet(n, func() string {
		return faker.Code().Ean8()
	})
	if err != nil {
		return nil, err
	}
	for i := 0; i < n; i++ {
		course := po.BaseCoursePO{
			Name:   faker.App().Name(),            // 2-3 words
			Code:   codes[i],                      // 8 digits
			Credit: float64(gen.Rand.Int()%6 + 1), // 1-6
		}
		baseCourses[i] = course
	}
	err = gen.db.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(baseCourses, gen.batchSize).Error
	return baseCourses, err
}

type Department struct {
	Name   string
	Weight int
}

// MockDepartmentDis 2024-2025-1 的课程数目分布数据
var MockDepartmentDis = []Department{
	{"体育系", 176},
	{"共青团上海交通大学委员会", 4},
	{"军事教研室", 19},
	{"农业与生物学院", 70},
	{"凯原法学院", 41},
	{"化学化工学院", 92},
	{"医学院", 18},
	{"国际与公共事务学院", 48},
	{"图书馆", 2},
	{"外国语学院", 179},
	{"媒体与传播学院", 73},
	{"学指委（学生处、团委、人武部）合署", 5},
	{"学生创新中心", 22},
	{"学生就业服务与职业发展中心", 5},
	{"安泰经济与管理学院", 123},
	{"密西根学院", 94},
	{"巴黎卓越工程师学院", 71},
	{"心理与行为科学研究院", 1},
	{"教务处", 8},
	{"教育学院", 4},
	{"数学科学学院", 94},
	{"智慧能源创新学院", 14},
	{"机械与动力工程学院", 180},
	{"材料科学与工程学院", 56},
	{"校医院", 1},
	{"海洋学院", 13},
	{"溥渊未来技术学院", 2},
	{"物理与天文学院", 66},
	{"环境科学与工程学院", 26},
	{"生命科学技术学院", 87},
	{"生物医学工程学院", 43},
	{"电子信息与电气工程学院", 365},
	{"研究生院", 718},
	{"继续教育学院", 1},
	{"致远学院", 163},
	{"航空航天学院", 26},
	{"船舶海洋与建筑工程学院", 117},
	{"药学院", 24},
	{"设计学院", 87},
	{"马克思主义学院", 132},
}

func GenDepartment(gen MockDBGenerator) string {
	total := 0
	for _, department := range MockDepartmentDis {
		total += department.Weight
	}
	idx := gen.Rand.Int() % total
	for _, department := range MockDepartmentDis {
		idx -= department.Weight
		if idx <= 0 {
			return department.Name
		}
	}
	panic("unreachable")
}

func MockCourses(gen MockDBGenerator, n int, baseCourses []po.BaseCoursePO, teachers []po.TeacherPO) ([]po.CoursePO, error) {
	Courses := make([]po.CoursePO, n)
	type uniqueCourse struct {
		Code       string
		CourseIdx  int
		TeacherIdx int
	}
	uniqueTuples, err := GenerateUniqueSet(n, func() uniqueCourse {
		courseIdx := gen.Rand.Int() % len(baseCourses)
		teacherIdx := gen.Rand.Int() % len(teachers)
		course := baseCourses[courseIdx]
		return uniqueCourse{Code: course.Code, TeacherIdx: teacherIdx, CourseIdx: courseIdx}
	})
	for i := 0; i < n; i++ {
		mainTeacher := teachers[uniqueTuples[i].TeacherIdx]
		course := po.CoursePO{
			Code:            uniqueTuples[i].Code,
			Name:            baseCourses[uniqueTuples[i].CourseIdx].Name,
			Credit:          baseCourses[uniqueTuples[i].CourseIdx].Credit,
			Department:      mainTeacher.Department,
			MainTeacherID:   int64(mainTeacher.ID),
			MainTeacherName: mainTeacher.Name,
		}
		Courses[i] = course
	}
	err = gen.db.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(Courses, gen.batchSize).Error
	return Courses, err
}

func MockOfferedCourses(gen MockDBGenerator, n int, courses []po.CoursePO, teachers []po.TeacherPO) ([]po.OfferedCoursePO, error) {
	offeredCourses := make([]po.OfferedCoursePO, n)
	type uniqueOfferedCourse struct {
		Semester  string
		CourseIdx int
	}
	uniqueOfferedCourses, err := GenerateUniqueSet(n, func() uniqueOfferedCourse {
		courseIdx := gen.Rand.Int() % len(courses)
		grade := gen.Rand.Intn(4) + 2020
		year := grade + gen.Rand.Intn(4)
		semester := gen.Rand.Intn(3) + 1
		return uniqueOfferedCourse{
			CourseIdx: courseIdx,
			Semester:  strconv.Itoa(year) + "-" + strconv.Itoa(year+1) + strconv.Itoa(semester),
		}
	})
	for i := 0; i < n; i++ {
		course := courses[uniqueOfferedCourses[i].CourseIdx]
		grade := gen.Rand.Intn(4) + 2020
		offeredCourse := po.OfferedCoursePO{
			CourseID:      int64(course.ID),
			MainTeacherID: course.MainTeacherID,
			Language:      "中文",
			Grade:         strconv.Itoa(grade),
			Semester:      uniqueOfferedCourses[i].Semester,
		}
		offeredCourses[i] = offeredCourse
	}
	err = gen.db.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(offeredCourses, gen.batchSize).Error
	return offeredCourses, err
}
