package main

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"jcourse_go/internal/config"
	"jcourse_go/internal/domain/notification"
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
	newOfferings := make([]entity.CourseOffering, 0)
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
		newOfferings = append(newOfferings, newOffering)
		newCourses = append(newCourses, newCourse)
	}

	if err = newDB.Model(&entity.Course{}).
		Clauses(clause.OnConflict{DoNothing: true}).
		CreateInBatches(&newCourses, 100).Error; err != nil {
		panic(err)
	}

	if err = newDB.Model(&entity.Curriculum{}).
		Clauses(clause.OnConflict{DoNothing: true}).
		CreateInBatches(&newCurriculums, 100).Error; err != nil {
		panic(err)
	}

	if err = newDB.Model(&entity.CourseOffering{}).
		Clauses(clause.OnConflict{DoNothing: true}).
		CreateInBatches(&newOfferings, 100).Error; err != nil {
		panic(err)
	}

	clear(newOfferings)
	if err = newDB.Model(&entity.CourseOffering{}).Find(&newOfferings).Error; err != nil {
		panic(err)
	}

	lastOfferingMap := make(map[int64]int64)
	for _, o := range newOfferings {
		lastOfferingMap[o.CourseID] = o.ID
	}

	newOfferingCategories := make([]entity.CourseOfferingCategory, 0)
	newOfferingTeacherGroup := make([]entity.CourseOfferingTeacher, 0)
	for _, c := range courses {
		offeringID := lastOfferingMap[c.ID]
		for _, ct := range c.Categories {
			category := entity.CourseOfferingCategory{
				CourseOfferingID: offeringID,
				Category:         ct.Name,
				CourseID:         c.ID,
			}
			newOfferingCategories = append(newOfferingCategories, category)
		}
		for _, t := range c.TeacherGroup {
			teacherGroup := entity.CourseOfferingTeacher{
				CourseOfferingID: offeringID,
				TeacherID:        t.ID,
				CourseID:         c.ID,
			}
			newOfferingTeacherGroup = append(newOfferingTeacherGroup, teacherGroup)
		}
	}

	if err = newDB.Model(&entity.CourseOfferingTeacher{}).
		Clauses(clause.OnConflict{DoNothing: true}).
		CreateInBatches(&newOfferingTeacherGroup, 100).Error; err != nil {
		panic(err)
	}

	if err = newDB.Model(&entity.CourseOfferingCategory{}).
		Clauses(clause.OnConflict{DoNothing: true}).
		CreateInBatches(&newOfferingCategories, 100).Error; err != nil {
		panic(err)
	}

	for i := range newCourses {
		newCourses[i].LastOfferingID = lastOfferingMap[newCourses[i].ID]
	}

	if err = newDB.Model(&entity.Course{}).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"last_offering_id"}),
		}).
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
	if err = oldDB.Model(&Review{}).
		Joins("User").Joins("Course").Joins("Semester").
		Preload("Reactions").Preload("Revisions.Semester").
		Find(&reviews).Error; err != nil {
		panic(err)
	}

	newReviews := make([]entity.Review, 0)
	for _, r := range reviews {
		newReview := entity.Review{
			ID:        r.ID,
			CourseID:  r.CourseID,
			UserID:    r.UserID,
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
				ID:        rv.ID,
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

	if err = newDB.Table("course c").Where("review_count = 0").Updates(map[string]interface{}{
		"review_count": newDB.Table("review r").
			Select("count(r.id)").
			Where("r.course_id = c.id"),
		"review_avg": newDB.Table("review r").
			Select("avg(r.rating)").
			Where("r.course_id = c.id"),
	}).Error; err != nil {
		panic(err)
	}

	notifications := make([]CourseNotificationLevel, 0)
	if err = oldDB.Model(&CourseNotificationLevel{}).Find(&notifications).Error; err != nil {
		panic(err)
	}

	newNotifications := make([]entity.CourseNotification, 0)
	for _, n := range notifications {
		if n.NotificationLevel == int64(notification.LevelNormal) {
			continue
		}
		newNotification := entity.CourseNotification{
			ID:        n.ID,
			CourseID:  n.CourseID,
			UserID:    n.UserID,
			Level:     n.NotificationLevel,
			CreatedAt: n.ModifiedAt,
			UpdatedAt: n.ModifiedAt,
		}
		newNotifications = append(newNotifications, newNotification)
	}

	if err = newDB.Model(&entity.CourseNotification{}).
		Clauses(clause.OnConflict{DoNothing: true}).
		CreateInBatches(&newNotifications, 100).Error; err != nil {
		panic(err)
	}

	enrolls := make([]EnrollCourse, 0)
	if err = oldDB.Model(&EnrollCourse{}).
		Joins("Semester").
		Find(&enrolls).Error; err != nil {
		panic(err)
	}

	newEnrollments := make([]entity.UserCourseEnrollment, 0)
	for _, enroll := range enrolls {
		if enroll.SemesterID == 0 {
			continue
		}
		newEnrollment := entity.UserCourseEnrollment{
			ID:        enroll.ID,
			CourseID:  enroll.CourseID,
			UserID:    enroll.UserID,
			Semester:  enroll.Semester.Name,
			CreatedAt: enroll.CreatedAt,
		}
		newEnrollments = append(newEnrollments, newEnrollment)
	}

	if err = newDB.Model(&entity.UserCourseEnrollment{}).
		Clauses(clause.OnConflict{DoNothing: true}).
		CreateInBatches(&newEnrollments, 100).Error; err != nil {
		panic(err)
	}
}
