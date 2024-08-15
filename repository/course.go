package repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"jcourse_go/dal"
	"jcourse_go/model/po"
)

type IBaseCourseQuery interface {
	GetBaseCourse(ctx context.Context, opts ...DBOption) (*po.BaseCoursePO, error)
	GetBaseCourseList(ctx context.Context, opts ...DBOption) ([]po.BaseCoursePO, error)
	WithCode(code string) DBOption
	WithName(name string) DBOption
	WithCredit(credit float64) DBOption
	WithIDs(IDs []int64) DBOption
	WithID(id int64) DBOption
}

type BaseCourseQuery struct {
	db *gorm.DB
}

func (b *BaseCourseQuery) WithID(id int64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}
func (b *BaseCourseQuery) optionDB(ctx context.Context, opts ...DBOption) *gorm.DB {
	db := b.db.WithContext(ctx).Model(&po.BaseCoursePO{})
	for _, opt := range opts {
		db = opt(db)
	}
	return db
}
func (b *BaseCourseQuery) WithIDs(IDs []int64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id in ?", IDs)
	}
}
func (b *BaseCourseQuery) GetBaseCourse(ctx context.Context, opts ...DBOption) (*po.BaseCoursePO, error) {
	db := b.optionDB(ctx, opts...)
	course := po.BaseCoursePO{}
	result := db.Debug().First(&course)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &course, nil
}

func (b *BaseCourseQuery) GetBaseCourseList(ctx context.Context, opts ...DBOption) ([]po.BaseCoursePO, error) {
	db := b.optionDB(ctx, opts...)
	coursePOs := make([]po.BaseCoursePO, 0)
	result := db.Debug().Find(&coursePOs)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return coursePOs, nil
}

func (b *BaseCourseQuery) WithCode(code string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("code = ?", code)
	}
}

func (b *BaseCourseQuery) WithName(name string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("name = ?", name)
	}
}

func (b *BaseCourseQuery) WithCredit(credit float64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("credit = ?", credit)
	}
}

func NewBaseCourseQuery() IBaseCourseQuery {
	return &BaseCourseQuery{db: dal.GetDBClient()}
}

type ICourseQuery interface {
	GetCourse(ctx context.Context, opts ...DBOption) (*po.CoursePO, error)
	GetCourseList(ctx context.Context, opts ...DBOption) ([]po.CoursePO, error)
	GetCourseCount(ctx context.Context, opts ...DBOption) (int64, error)
	GetCourseCategories(ctx context.Context, courseIDs []int64) (map[int64][]string, error)
	GetCourseByIDs(ctx context.Context, courseIDs []int64) (map[int64]po.CoursePO, error)
	WithID(id int64) DBOption
	WithCode(code string) DBOption
	WithName(name string) DBOption
	WithCredits(credits []float64) DBOption
	WithCategories(categories []string) DBOption
	WithDepartments(departments []string) DBOption
	WithMainTeacherName(name string) DBOption
	WithMainTeacherID(id int64) DBOption
	WithLimit(limit int64) DBOption
	WithOffset(offset int64) DBOption
}

type CourseQuery struct {
	db *gorm.DB
}

func (c *CourseQuery) GetCourseByIDs(ctx context.Context, courseIDs []int64) (map[int64]po.CoursePO, error) {
	db := c.optionDB(ctx)
	courses := make([]po.CoursePO, 0)
	coursesMap := make(map[int64]po.CoursePO)
	result := db.Where("id in ?", courseIDs).Find(&courses)
	if result.Error != nil {
		return coursesMap, result.Error
	}
	for _, course := range courses {
		coursesMap[int64(course.ID)] = course
	}
	return coursesMap, nil
}

func (c *CourseQuery) WithLimit(limit int64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(int(limit))
	}
}

func (c *CourseQuery) WithOffset(offset int64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(int(offset))
	}
}

func (c *CourseQuery) optionDB(ctx context.Context, opts ...DBOption) *gorm.DB {
	db := c.db.WithContext(ctx).Model(po.CoursePO{})
	for _, opt := range opts {
		db = opt(db)
	}
	return db
}

func (c *CourseQuery) GetCourseCategories(ctx context.Context, courseIDs []int64) (map[int64][]string, error) {
	db := c.db.WithContext(ctx).Model(po.CourseCategoryPO{})
	courseCategoryPOs := make([]po.CourseCategoryPO, 0)
	result := db.Where("course_id in ?", courseIDs).Find(&courseCategoryPOs)
	if result.Error != nil {
		return nil, result.Error
	}
	courseCategoryMap := make(map[int64][]string)
	for _, courseCategoryPO := range courseCategoryPOs {
		categories, ok := courseCategoryMap[courseCategoryPO.CourseID]
		if !ok {
			categories = make([]string, 0)
		}
		categories = append(categories, courseCategoryPO.Category)
		courseCategoryMap[courseCategoryPO.CourseID] = categories
	}
	return courseCategoryMap, nil
}

func (c *CourseQuery) GetCourse(ctx context.Context, opts ...DBOption) (*po.CoursePO, error) {
	db := c.optionDB(ctx, opts...)
	course := po.CoursePO{}
	result := db.First(&course)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &course, nil
}

func (c *CourseQuery) GetCourseList(ctx context.Context, opts ...DBOption) ([]po.CoursePO, error) {
	db := c.optionDB(ctx, opts...)
	coursePOs := make([]po.CoursePO, 0)
	result := db.Find(&coursePOs)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return coursePOs, nil
}

func (c *CourseQuery) GetCourseCount(ctx context.Context, opts ...DBOption) (int64, error) {
	db := c.optionDB(ctx, opts...)
	var count int64
	result := db.Model(&po.CoursePO{}).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

func (c *CourseQuery) WithID(id int64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func (c *CourseQuery) WithCode(code string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("code = ?", code)
	}
}

func (c *CourseQuery) WithName(name string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("name = ?", name)
	}
}

func (c *CourseQuery) WithCredits(credits []float64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("credit in ?", credits)
	}
}

func (c *CourseQuery) WithDepartments(departments []string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("department in ?", departments)
	}
}

func (c *CourseQuery) WithCategories(categories []string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Joins("inner join course_categories on course_categories.course_id = courses.id").Where("category in ?", categories)
	}
}

func (c *CourseQuery) WithMainTeacherName(name string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("main_teacher_name = ?", name)
	}
}

func (c *CourseQuery) WithMainTeacherID(id int64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("main_teacher_id = ?", id)
	}
}

func NewCourseQuery() ICourseQuery {
	return &CourseQuery{db: dal.GetDBClient()}
}

type IOfferedCourseQuery interface {
	GetOfferedCourse(ctx context.Context, opts ...DBOption) (*po.OfferedCoursePO, error)
	GetOfferedCourseList(ctx context.Context, opts ...DBOption) ([]po.OfferedCoursePO, error)
	GetOfferedCourseTeacherGroup(ctx context.Context, offeredCourseIDs []int64) (map[int64][]po.TeacherPO, error)
	// 获取教过courseIDs之中所有课的主教师的ID List
	GetMainTeacherIDsWithOfferedCourseIDs(ctx context.Context, courseIDs []int64) ([]int64, error)
	WithID(id int64) DBOption
	WithIDs(ids []int64) DBOption
	WithCourseID(id int64) DBOption
	WithMainTeacherID(id int64) DBOption
	WithSemester(semester string) DBOption
	WithOrderBy(field string, ascending bool) DBOption
}

type OfferedCourseQuery struct {
	db *gorm.DB
}

func (o *OfferedCourseQuery) WithOrderBy(field string, ascending bool) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		if ascending {
			field = fmt.Sprintf("%s %s", field, "asc")
		} else {
			field = fmt.Sprintf("%s %s", field, "desc")
		}
		return db.Order(field)
	}
}

func (o *OfferedCourseQuery) GetOfferedCourseTeacherGroup(ctx context.Context, offeredCourseIDs []int64) (map[int64][]po.TeacherPO, error) {
	db := o.db.WithContext(ctx).Model(&po.OfferedCourseTeacherPO{})
	courseTeacherPOs := make([]po.OfferedCourseTeacherPO, 0)
	result := db.Where("offered_course_id in ?", offeredCourseIDs).Find(&courseTeacherPOs)
	if result.Error != nil {
		return nil, result.Error
	}
	courseTeacherMap := make(map[int64][]po.TeacherPO)
	for _, courseTeacher := range courseTeacherPOs {
		val, ok := courseTeacherMap[courseTeacher.OfferedCourseID]
		if !ok {
			val = make([]po.TeacherPO, 0)
		}
		teacher := po.TeacherPO{Name: courseTeacher.TeacherName}
		teacher.ID = uint(courseTeacher.TeacherID)
		val = append(val, teacher)
		courseTeacherMap[courseTeacher.OfferedCourseID] = val
	}
	return courseTeacherMap, nil
}

func (o *OfferedCourseQuery) optionDB(ctx context.Context, opts ...DBOption) *gorm.DB {
	db := o.db.WithContext(ctx).Model(po.OfferedCoursePO{})
	for _, opt := range opts {
		db = opt(db)
	}
	return db
}

func (o *OfferedCourseQuery) GetOfferedCourse(ctx context.Context, opts ...DBOption) (*po.OfferedCoursePO, error) {
	db := o.optionDB(ctx, opts...)
	course := po.OfferedCoursePO{}
	result := db.First(&course)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &course, nil
}

func (o *OfferedCourseQuery) GetOfferedCourseList(ctx context.Context, opts ...DBOption) ([]po.OfferedCoursePO, error) {
	db := o.optionDB(ctx, opts...)
	coursePOs := make([]po.OfferedCoursePO, 0)
	result := db.Find(&coursePOs)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return coursePOs, nil
}
func (o *OfferedCourseQuery) GetMainTeacherIDsWithOfferedCourseIDs(ctx context.Context, courseIDs []int64) ([]int64, error) {
	main_teacher_ids := make([]int64, 0)
	db := o.optionDB(ctx)
	if courseIDs == nil || len(courseIDs) == 0 {
		result := db.Distinct("main_teacher_id").Pluck("main_teacher_id", &main_teacher_ids)
		if result.Error != nil {
			return nil, result.Error
		}
		return main_teacher_ids, nil
	}
	result := db.Where("id IN ?", courseIDs).
		Group("training_plan_id").
		Having("count(DISTINCT id) = ?", len(courseIDs)).Find(&main_teacher_ids)
	if result.Error != nil {
		return nil, result.Error
	}
	return main_teacher_ids, nil
}

func (o *OfferedCourseQuery) WithID(id int64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}
func (o *OfferedCourseQuery) WithIDs(ids []int64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id in ?", ids)
	}
}
func (o *OfferedCourseQuery) WithCourseID(id int64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("course_id = ?", id)
	}
}

func (o *OfferedCourseQuery) WithMainTeacherID(id int64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("main_teacher_id = ?", id)
	}
}

func (o *OfferedCourseQuery) WithSemester(semester string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("semester = ?", semester)
	}
}

func NewOfferedCourseQuery() IOfferedCourseQuery {
	return &OfferedCourseQuery{db: dal.GetDBClient()}
}
