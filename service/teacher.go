package service

import (
	"context"
	"jcourse_go/model/domain"
	"jcourse_go/repository"
)
func BuildFilterFromText(c context.Context, text string) domain.TeacherListFilter {
	// HINT: 教师名/拼音/拼音缩写
	teacherQuery := repository.NewTeacherQuery()
	filter := domain.TeacherListFilter{}
	_, err := teacherQuery.GetTeacher(c, teacherQuery.WithName(text))
	if err != nil {
		filter.Name = text
		return filter
	}
	_, err = teacherQuery.GetTeacher(c, teacherQuery.WithPinyin(text))
	if err != nil {
		filter.Pinyin = text
		return filter
	}
	_, err = teacherQuery.GetTeacher(c, teacherQuery.WithPinyinAbbr(text))
	if err != nil {
		filter.PinyinAbbr = text
		return filter
	}
	_, err = teacherQuery.GetTeacher(c, teacherQuery.WithDepartment(text))
	if err != nil {
		filter.Department = text
		return filter
	}
	return filter
}
func buildTeacherDBOptionsFromFilter(query repository.ITeacherQuery, filter domain.TeacherListFilter) []repository.DBOption {
	options := make([]repository.DBOption, 0)

	if filter.Name != "" {
		options = append(options, query.WithName(filter.Name))
	}
	if filter.Pinyin != "" {
		options = append(options, query.WithPinyin(filter.Pinyin))
	}
	if filter.PinyinAbbr != "" {
		options = append(options, query.WithPinyinAbbr(filter.PinyinAbbr))
	}
	if filter.Department != "" {
		options = append(options, query.WithDepartment(filter.Department))
	}
	return options
}

func GetTeacherList(c context.Context, filter domain.TeacherListFilter) ([]domain.Teacher, error) {
	return nil, nil
}

func GetTeacherDetail(c context.Context, teacherID int64) (*domain.Teacher, error) {
	return nil, nil
}

