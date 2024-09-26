package repository

import "gorm.io/gorm"

type DBOption func(*gorm.DB) *gorm.DB

func WithUserIDs(userIDs []int64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id in ?", userIDs)
	}
}

func WithUserID(id int64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", id)
	}
}

func WithLimit(limit int64) DBOption {
	return func(db *gorm.DB) *gorm.DB { return db.Limit(int(limit)) }
}

func WithOffset(offset int64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(int(offset))
	}
}

func WithNotAnonymous() DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("is_anonymous = ?", false)
	}
}

func WithEmail(email string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("email = ?", email)
	}
}

func WithID(id int64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func WithIDs(ids []int64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id in ?", ids)
	}
}

func WithPassword(password string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("password = ?", password)
	}
}

func WithCode(code string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("code = ?", code)
	}
}

func WithName(name string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("name = ?", name)
	}
}

func WithCredit(credit float64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("credit = ?", credit)
	}
}

func WithCredits(credits []float64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("credit in ?", credits)
	}
}

func WithDepartment(department string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("department = ?", department)
	}
}

func WithDepartments(departments []string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("department in ?", departments)
	}
}

func WithCategories(categories []string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Joins("inner join course_categories on course_categories.course_id = courses.id").Where("category in ?", categories)
	}
}

func WithMainTeacherName(name string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("main_teacher_name = ?", name)
	}
}

func WithMainTeacherID(id int64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("main_teacher_id = ?", id)
	}
}

func WithCourseID(courseID int64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("course_id = ?", courseID)
	}
}

func WithCourseIDs(courseIDs []int64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("course_id in ?", courseIDs)
	}
}

func WithSemester(semester string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("semester = ?", semester)
	}
}

func WithTitle(title string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("title = ?", title)
	}
}

func WithTitles(titles []string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("title in ?", titles)
	}
}

func WithOrderBy(orderBy string, ascending bool) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		if ascending {
			orderBy = orderBy + " asc"
		} else {
			orderBy = orderBy + " desc"
		}
		return db.Order(orderBy)
	}
}

func WithSearch(query string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("search_index @@ to_tsquery('simple', ?)",
			userQueryToTsQuery(query),
		)
	}
}

// WithNameSearch 目前只搜用户名
func WithNameSearch(query string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("name LIKE ?", query+"%")
	}
}

func WithPinyin(pinyin string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("pinyin = ?", pinyin)
	}
}

func WithPinyinAbbr(pinyin string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("pinyin_abbr = ?", pinyin)
	}
}

func WithPaginate(page int64, pageSize int64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 || pageSize <= 0 {
			return db.Where("1 = 0")
		}
		return db.Offset(int((page - 1) * pageSize)).Limit(int(pageSize))
	}
}

func WithMajor(major string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("major = ?", major)
	}
}

func WithEntryYears(entryYears []string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("entry_year in ?", entryYears)
	}
}

func WithDegrees(degrees []string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("degree in ?", degrees)
	}
}

func WithSuggestSemester(semester string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("suggest_semester = ?", semester)
	}
}

func WithTrainingPlanID(trainingPlanID int64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("training_plan_id = ?", trainingPlanID)
	}
}
