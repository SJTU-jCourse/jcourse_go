package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"jcourse_go/model/model"
	"jcourse_go/model/po"
)

type IBaseCourseQuery interface {
	GetBaseCourse(ctx context.Context, opts ...DBOption) ([]po.BaseCoursePO, error)
	GetBaseCourseCount(ctx context.Context, opts ...DBOption) (int64, error)
	GetBaseCoursesByIDs(ctx context.Context, ids []int64) (map[int64]po.BaseCoursePO, error)
}

type BaseCourseQuery struct {
	db *gorm.DB
}

func (b *BaseCourseQuery) GetBaseCourseCount(ctx context.Context, opts ...DBOption) (int64, error) {
	db := b.optionDB(ctx, opts...)
	var count int64
	result := db.Model(&po.BaseCoursePO{}).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

func (b *BaseCourseQuery) GetBaseCoursesByIDs(ctx context.Context, ids []int64) (map[int64]po.BaseCoursePO, error) {
	db := b.optionDB(ctx)
	courses := make([]po.BaseCoursePO, 0)
	coursesMap := make(map[int64]po.BaseCoursePO)
	result := db.Where("id in ?", ids).Find(&courses)
	if result.Error != nil {
		return coursesMap, result.Error
	}
	for _, course := range courses {
		coursesMap[int64(course.ID)] = course
	}
	return coursesMap, nil
}

func (b *BaseCourseQuery) optionDB(ctx context.Context, opts ...DBOption) *gorm.DB {
	db := b.db.WithContext(ctx).Model(&po.BaseCoursePO{})
	for _, opt := range opts {
		db = opt(db)
	}
	return db
}

func (b *BaseCourseQuery) GetBaseCourse(ctx context.Context, opts ...DBOption) ([]po.BaseCoursePO, error) {
	db := b.optionDB(ctx, opts...)
	coursePOs := make([]po.BaseCoursePO, 0)
	result := db.Find(&coursePOs)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return coursePOs, nil
}

func NewBaseCourseQuery(db *gorm.DB) IBaseCourseQuery {
	return &BaseCourseQuery{db: db}
}

type ICourseQuery interface {
	GetCourse(ctx context.Context, opts ...DBOption) ([]po.CoursePO, error)
	GetCourseCount(ctx context.Context, opts ...DBOption) (int64, error)
	GetCourseCategories(ctx context.Context, courseIDs []int64) (map[int64][]string, error)
	GetCourseByIDs(ctx context.Context, courseIDs []int64) (map[int64]po.CoursePO, error)
	GetCourseFilter(ctx context.Context) (model.CourseFilter, error)
}

type CourseQuery struct {
	db *gorm.DB
}

func (c *CourseQuery) GetCourseFilter(ctx context.Context) (model.CourseFilter, error) {
	filter := model.CourseFilter{
		Departments: make([]model.FilterItem, 0),
		Credits:     make([]model.FilterItem, 0),
		Semesters:   make([]model.FilterItem, 0),
		Categories:  make([]model.FilterItem, 0),
	}

	err := c.db.WithContext(ctx).Model(&po.CoursePO{}).
		Select("department as value, count(*) as count").
		Group("department").Find(&filter.Departments).Error
	if err != nil {
		return filter, err
	}

	err = c.db.WithContext(ctx).Model(&po.CoursePO{}).
		Select("cast(credit as varchar) as value, count(*) as count").
		Group("credit").Find(&filter.Credits).Error
	if err != nil {
		return filter, err
	}

	err = c.db.WithContext(ctx).Model(&po.CourseCategoryPO{}).
		Select("category as value, count(course_id) as count").
		Group("category").Find(&filter.Categories).Error
	if err != nil {
		return filter, err
	}

	err = c.db.WithContext(ctx).Model(&po.OfferedCoursePO{}).
		Select("semester as value, count(course_id) as count").
		Group("semester").Find(&filter.Semesters).Error
	if err != nil {
		return filter, err
	}
	return filter, nil
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

func (c *CourseQuery) optionDB(ctx context.Context, opts ...DBOption) *gorm.DB {
	db := c.db.WithContext(ctx).Model(po.CoursePO{})
	for _, opt := range opts {
		db = opt(db)
	}
	return db
}

func (c *CourseQuery) GetCourseCategories(ctx context.Context, courseIDs []int64) (map[int64][]string, error) {
	courseCategoryMap := make(map[int64][]string)
	for _, id := range courseIDs {
		courseCategoryMap[id] = make([]string, 0)
	}

	db := c.db.WithContext(ctx).Model(po.CourseCategoryPO{})
	courseCategoryPOs := make([]po.CourseCategoryPO, 0)
	result := db.Where("course_id in ?", courseIDs).Find(&courseCategoryPOs)
	if result.Error != nil {
		return courseCategoryMap, result.Error
	}

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

func (c *CourseQuery) GetCourse(ctx context.Context, opts ...DBOption) ([]po.CoursePO, error) {
	db := c.optionDB(ctx, opts...)
	coursePOs := make([]po.CoursePO, 0)
	err := db.Find(&coursePOs).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return coursePOs, nil
	}
	if err != nil {
		return coursePOs, err
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

func NewCourseQuery(db *gorm.DB) ICourseQuery {
	return &CourseQuery{db: db}
}

type IOfferedCourseQuery interface {
	GetOfferedCourse(ctx context.Context, opts ...DBOption) ([]po.OfferedCoursePO, error)
	GetOfferedCourseTeacherGroup(ctx context.Context, offeredCourseIDs []int64) (map[int64][]po.TeacherPO, error)
	// 获取教过courseIDs之中所有课的主教师的ID List
	GetMainTeacherIDsWithOfferedCourseIDs(ctx context.Context, courseIDs []int64) ([]int64, error)
}

type OfferedCourseQuery struct {
	db *gorm.DB
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

func (o *OfferedCourseQuery) GetOfferedCourse(ctx context.Context, opts ...DBOption) ([]po.OfferedCoursePO, error) {
	db := o.optionDB(ctx, opts...)
	coursePOs := make([]po.OfferedCoursePO, 0)
	result := db.Find(&coursePOs)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return coursePOs, nil
}
func (o *OfferedCourseQuery) GetMainTeacherIDsWithOfferedCourseIDs(ctx context.Context, courseIDs []int64) ([]int64, error) {
	mainTeacherIDs := make([]int64, 0)
	db := o.optionDB(ctx)
	if len(courseIDs) == 0 {
		result := db.Distinct("main_teacher_id").Pluck("main_teacher_id", &mainTeacherIDs)
		if result.Error != nil {
			return nil, result.Error
		}
		return mainTeacherIDs, nil
	}
	result := db.Where("id IN ?", courseIDs).
		Group("training_plan_id").
		Having("count(DISTINCT id) = ?", len(courseIDs)).Find(&mainTeacherIDs)
	if result.Error != nil {
		return nil, result.Error
	}
	return mainTeacherIDs, nil
}

func NewOfferedCourseQuery(db *gorm.DB) IOfferedCourseQuery {
	return &OfferedCourseQuery{db: db}
}
