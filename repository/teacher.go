package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"jcourse_go/dal"
	"jcourse_go/model/po"
)

// TODO: 暂时没有添加从ProfileDesc中搜索的功能
type ITeacherQuery interface {
	GetTeacher(ctx context.Context, opts ...DBOption) (*po.TeacherPO, error)
	GetTeacherList(ctx context.Context, opts ...DBOption) ([]po.TeacherPO, error)
	GetTeacherCount(ctx context.Context, opts ...DBOption) (int64, error)
	WithID(id int64) DBOption
	WithCode(code string) DBOption
	WithName(name string) DBOption
	WithPinyin(pinyin string) DBOption
	WithPinyinAbbr(pinyin string) DBOption
	WithDepartment(department string) DBOption
	WithTitle(title string) DBOption
	WithPicture(picture string) DBOption
	WithProfileURL(profileURL string) DBOption
	WithIDs(ids []int64) DBOption
	WithPaginate(page int64, pageSize int64) DBOption
}

type TeacherQuery struct {
	db *gorm.DB
}

func NewTeacherQuery() ITeacherQuery {
	return &TeacherQuery{
		db: dal.GetDBClient(),
	}
}

func (q *TeacherQuery) optionDB(ctx context.Context, opts ...DBOption) *gorm.DB {
	db := q.db.WithContext(ctx).Model(po.TeacherPO{})
	for _, opt := range opts {
		db = opt(db)
	}
	return db
}

func (q *TeacherQuery) GetTeacher(ctx context.Context, opts ...DBOption) (*po.TeacherPO, error) {
	db := q.optionDB(ctx, opts...)
	teacher := po.TeacherPO{}
	result := db.First(&teacher)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &teacher, nil
}

func (q *TeacherQuery) GetTeacherList(ctx context.Context, opts ...DBOption) ([]po.TeacherPO, error) {
	db := q.optionDB(ctx, opts...)
	teacherPOs := make([]po.TeacherPO, 0)
	result := db.Find(&teacherPOs)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return teacherPOs, nil
}

func (q *TeacherQuery) GetTeacherCount(ctx context.Context, opts ...DBOption) (int64, error) {
	db := q.optionDB(ctx, opts...)
	var count int64
	result := db.Model(&po.TeacherPO{}).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

func (q *TeacherQuery) WithID(id int64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func (q *TeacherQuery) WithCode(code string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("code = ?", code)
	}
}

func (q *TeacherQuery) WithName(name string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("name = ?", name)
	}
}

func (q *TeacherQuery) WithPinyin(pinyin string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("pinyin = ?", pinyin)
	}
}

func (q *TeacherQuery) WithPinyinAbbr(pinyin string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("pinyin_abbr = ?", pinyin)
	}
}

func (q *TeacherQuery) WithDepartment(department string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("department = ?", department)
	}
}

func (q *TeacherQuery) WithTitle(title string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("title = ?", title)
	}
}

func (q *TeacherQuery) WithPicture(picture string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("picture = ?", picture)
	}
}

func (q *TeacherQuery) WithProfileURL(profileURL string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("profile_url = ?", profileURL)
	}
}

func (q *TeacherQuery) WithIDs(ids []int64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id IN ?", ids)
	}
}

func (q *TeacherQuery) WithPaginate(page int64, pageSize int64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 || pageSize <= 0 {
			return db.Where("1 = 0")
		}
		return db.Offset(int((page - 1) * pageSize)).Limit(int(pageSize))
	}

}