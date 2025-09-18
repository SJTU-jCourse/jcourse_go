package main

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"jcourse_go/internal/config"
	"jcourse_go/internal/infrastructure/dal"
	"jcourse_go/internal/infrastructure/entity"
)

var (
	oldDB *gorm.DB
	newDB *gorm.DB
)

func main() {
	var err error
	oldDB, err = dal.NewPostgresSQL(config.DBConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "jcourse",
		Password: "jcourse",
		DBName:   "jcourse",
		Debug:    false,
	})
	if err != nil {
		panic(err)
	}
	newDB, err = dal.NewPostgresSQL(config.DBConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "jcourse",
		Password: "jcourse",
		DBName:   "jcourse_v2",
		Debug:    false,
	})
	if err != nil {
		panic(err)
	}

	semesters := make([]Semester, 0)
	if err = oldDB.Model(&Semester{}).Find(&semesters).Error; err != nil {
		panic(err)
	}

	newSemesters := make([]entity.Semester, 0)
	for _, s := range semesters {
		newSemester := entity.Semester{
			ID:        s.ID,
			Name:      s.Name,
			Available: s.Available,
		}
		newSemesters = append(newSemesters, newSemester)
	}

	if err = newDB.Model(&entity.Semester{}).
		Clauses(clause.OnConflict{DoNothing: true}).
		CreateInBatches(&newSemesters, 100).Error; err != nil {
		panic(err)
	}

	fallbackSemester := semesters[len(semesters)-1]

	teachers := make([]Teacher, 0)
	if err = oldDB.Model(&Teacher{}).Joins("Department").Find(&teachers).Error; err != nil {
		panic(err)
	}

	newTeachers := make([]entity.Teacher, 0)
	for _, t := range teachers {
		newTeacher := entity.Teacher{
			ID:         t.ID,
			Name:       t.Name,
			Pinyin:     t.Pinyin,
			PinyinAbbr: t.AbbrPinyin,
			Code:       t.Tid,
			Department: t.Department.Name,
			Title:      t.Title,
		}
		newTeachers = append(newTeachers, newTeacher)
	}

	if err = newDB.Model(&entity.Teacher{}).
		Clauses(clause.OnConflict{DoNothing: true}).
		CreateInBatches(&newTeachers, 100).Error; err != nil {
		panic(err)
	}

	courses := make([]Course, 0)
	if err = oldDB.Model(&Course{}).
		Joins("MainTeacher").Joins("Department").Joins("LastSemester").
		Preload("TeacherGroup").Preload("Categories").
		Find(&courses).Error; err != nil {
		panic(err)
	}

	newCourses := make([]entity.Course, 0)
	newCurriculums := make([]entity.Curriculum, 0)
	curriculumMap := make(map[string]entity.Curriculum)
	for _, c := range courses {
		if c.LastSemester == nil {
			c.LastSemester = &fallbackSemester
		}
		newCourse := entity.Course{
			ID:            c.ID,
			Code:          c.Code,
			Name:          c.Name,
			Credit:        c.Credit,
			MainTeacherID: c.MainTeacherID,
			Offerings:     make([]*entity.CourseOffering, 0),
		}
		if _, ok := curriculumMap[c.Code]; !ok {
			curriculum := entity.Curriculum{
				Code:   c.Code,
				Name:   c.Name,
				Credit: c.Credit,
			}
			curriculumMap[c.Code] = curriculum
			newCurriculums = append(newCurriculums, curriculum)
		}
		newOffering := entity.CourseOffering{
			CourseID:      c.ID,
			MainTeacherID: c.MainTeacherID,
			Semester:      c.LastSemester.Name,
			Department:    c.Department.Name,
			Categories:    make([]entity.CourseOfferingCategory, 0),
			TeacherGroup:  make([]entity.CourseOfferingTeacher, 0),
		}
		for _, ct := range c.Categories {
			category := entity.CourseOfferingCategory{
				Category: ct.Name,
				CourseID: c.ID,
			}
			newOffering.Categories = append(newOffering.Categories, category)
		}
		for _, t := range c.TeacherGroup {
			teacherGroup := entity.CourseOfferingTeacher{
				TeacherID: t.ID,
				CourseID:  c.ID,
			}
			newOffering.TeacherGroup = append(newOffering.TeacherGroup, teacherGroup)
		}
		newCourse.Offerings = append(newCourse.Offerings, &newOffering)
		newCourses = append(newCourses, newCourse)
	}

	if err = newDB.Model(&entity.Curriculum{}).
		Clauses(clause.OnConflict{DoNothing: true}).
		CreateInBatches(&newCurriculums, 100).Error; err != nil {
		panic(err)
	}

	if err = newDB.Model(&entity.Course{}).
		Clauses(clause.OnConflict{DoNothing: true}).
		CreateInBatches(&newCourses, 100).Error; err != nil {
		panic(err)
	}

	users := make([]User, 0)
	if err = oldDB.Model(&User{}).
		Joins("UserProfile").Preload("Points").
		Find(&users).Error; err != nil {
		panic(err)
	}

	newUsers := make([]entity.User, 0)
	newUserPoints := make([]entity.UserPoint, 0)
	for _, u := range users {
		newUser := entity.User{
			ID:        u.ID,
			Username:  u.Username,
			Password:  u.Password,
			UserRole:  "normal",
			CreatedAt: u.DateJoined,
		}
		if u.UserProfile != nil {
			newUser.LowerCase = u.UserProfile.LowerCase
			newUser.LastSeenAt = u.UserProfile.LastSeenAt
			newUser.SuspendedTill = u.UserProfile.SuspendedTill
		}
		for _, p := range u.Points {
			newPoint := entity.UserPoint{
				ID:          p.ID,
				UserID:      p.UserID,
				Value:       p.Value,
				Description: p.Description,
				Type:        "migrate",
				CreatedAt:   p.Time,
			}
			newUserPoints = append(newUserPoints, newPoint)
		}
		newUsers = append(newUsers, newUser)
	}
	if err = newDB.Model(&entity.User{}).
		Clauses(clause.OnConflict{DoNothing: true}).
		CreateInBatches(&newUsers, 100).Error; err != nil {
		panic(err)
	}
	if err = newDB.Model(&entity.UserPoint{}).
		Clauses(clause.OnConflict{DoNothing: true}).
		CreateInBatches(&newUserPoints, 100).Error; err != nil {
		panic(err)
	}

	reviews := make([]Review, 0)
	err = oldDB.Model(&Review{}).
		Joins("User").Joins("Course").Joins("Semester").
		Preload("Reactions").Preload("Revisions.Semester").
		Find(&reviews).Error
	if err != nil {
		panic(err)
	}

	newReviews := make([]entity.Review, 0)
	for _, r := range reviews {
		newReview := entity.Review{
			ID:        r.ID,
			CourseID:  r.CourseID,
			Course:    nil,
			UserID:    r.UserID,
			User:      nil,
			Comment:   r.Comment,
			Rating:    r.Rating,
			Semester:  r.Semester.Name,
			IsPublic:  false,
			Score:     r.Score,
			Revisions: make([]*entity.ReviewRevision, 0),
			Reactions: make([]*entity.ReviewReaction, 0),
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.ModifiedAt,
		}

		for _, rv := range r.Revisions {
			newRv := entity.ReviewRevision{
				UserID:    rv.UserID,
				ReviewID:  rv.ReviewID,
				Comment:   rv.Comment,
				Rating:    rv.Rating,
				Semester:  rv.Semester.Name,
				IsPublic:  false,
				Score:     rv.Score,
				UpdatedBy: rv.UserID,
				CreatedAt: rv.CreatedAt,
			}
			newReview.Revisions = append(newReview.Revisions, &newRv)
		}

		for _, reaction := range r.Reactions {
			reactionVal := "+1"
			if reaction.Reaction == 0 {
				continue
			}
			if reaction.Reaction == 1 {
				reactionVal = "+1"
			} else if reaction.Reaction == -1 {
				reactionVal = "-1"
			}
			newReaction := entity.ReviewReaction{
				ID:        reaction.ID,
				ReviewID:  reaction.ReviewID,
				UserID:    reaction.UserID,
				Reaction:  reactionVal,
				CreatedAt: reaction.ModifiedAt,
			}
			newReview.Reactions = append(newReview.Reactions, &newReaction)
		}
		newReviews = append(newReviews, newReview)
	}

	if err = newDB.Model(&entity.Review{}).
		Clauses(clause.OnConflict{DoNothing: true}).
		CreateInBatches(&newReviews, 100).Error; err != nil {
		panic(err)
	}
}
