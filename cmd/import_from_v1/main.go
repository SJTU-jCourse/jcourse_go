package main

import (
	"fmt"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"jcourse_go/dal"
	"jcourse_go/model/po"
	"jcourse_go/util"
)

func InitOldDB() *gorm.DB {
	host := util.GetPostgresHost()
	port := "5433"
	user := util.GetPostgresUser()
	password := util.GetPostgresPassword()
	dbname := util.GetPostgresDBName()
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", host, user, password, dbname, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

var (
	oldDB             *gorm.DB
	newDB             *gorm.DB
	semesterMap       map[int64]string
	departmentMap     map[int64]string
	courseMap         map[int64]CourseV1
	teacherMap        map[int64]TeacherV1
	userMap           map[int64]UserV1
	userProfileMap    map[int64]UserProfileV1
	reviewMap         map[int64]ReviewV1
	reviewRevisionMap map[int64]ReviewRevisionV1
	userPointMap      map[int64]UserPointV1
	reviewReactionMap map[int64]ReviewReactionV1

	oldToNewCourseMap  map[int64]po.CoursePO  // old id -> new model
	oldToNewTeacherMap map[int64]po.TeacherPO // old id -> new model

	newBaseCourseMap     map[string]po.BaseCoursePO     // code -> new
	newCourseMap         map[string]po.CoursePO         // course key -> new
	newTeacherMap        map[string]po.TeacherPO        // teacher code -> new
	newUserMap           map[int64]po.UserPO            // uid -> new
	newReviewMap         map[int64]po.ReviewPO          // uid -> new
	newReviewRevisionMap map[int64]po.ReviewRevisionPO  // uid -> new
	newUserPointMap      map[int64]po.UserPointDetailPO // uid -> new
	newReviewReactionMap map[int64]po.ReviewReactionPO  // uid -> new
)

func makeCourseKey(course po.CoursePO) string {
	return fmt.Sprintf("%s:%s", course.Code, course.MainTeacherName)
}

func loadOldSemester() {
	semesters := make([]SemesterV1, 0)
	oldDB.Model(&SemesterV1{}).Find(&semesters)
	semesterMap = make(map[int64]string)
	for _, semester := range semesters {
		semesterMap[semester.ID] = semester.Name
	}
}

func loadOldDepartment() {
	departments := make([]DepartmentV1, 0)
	oldDB.Model(&DepartmentV1{}).Find(&departments)
	departmentMap = make(map[int64]string)
	for _, department := range departments {
		departmentMap[department.ID] = department.Name
	}
}

func loadOldTeacher() {
	teachers := make([]TeacherV1, 0)
	oldDB.Model(&TeacherV1{}).Find(&teachers)
	teacherMap = make(map[int64]TeacherV1)
	for _, teacher := range teachers {
		teacherMap[teacher.ID] = teacher
	}
}

func loadOldCourse() {
	courses := make([]CourseV1, 0)
	oldDB.Model(&CourseV1{}).Find(&courses)
	courseMap = make(map[int64]CourseV1)
	for _, course := range courses {
		courseMap[course.ID] = course
	}
}

func loadOldUser() {
	users := make([]UserV1, 0)
	oldDB.Model(&UserV1{}).Find(&users)
	userMap = make(map[int64]UserV1)
	for _, user := range users {
		userMap[user.ID] = user
	}
}

func loadOldUserProfile() {
	userProfiles := make([]UserProfileV1, 0)
	oldDB.Model(&UserProfileV1{}).Find(&userProfiles)
	userProfileMap = make(map[int64]UserProfileV1)
	for _, userProfile := range userProfiles {
		userProfileMap[userProfile.UserID] = userProfile
	}
}

func loadOldUserPoint() {
	userPoints := make([]UserPointV1, 0)
	oldDB.Model(&UserPointV1{}).Find(&userPoints)
	userPointMap = make(map[int64]UserPointV1)
	for _, userPoint := range userPoints {
		userPointMap[userPoint.ID] = userPoint
	}
}

func loadOldReview() {
	reviews := make([]ReviewV1, 0)
	oldDB.Model(&ReviewV1{}).Find(&reviews)
	reviewMap = make(map[int64]ReviewV1)
	for _, review := range reviews {
		reviewMap[review.ID] = review
	}
}

func loadOldReviewRevision() {
	reviewRevisions := make([]ReviewRevisionV1, 0)
	oldDB.Model(&ReviewRevisionV1{}).Find(&reviewRevisions)
	reviewRevisionMap = make(map[int64]ReviewRevisionV1)
	for _, reviewRevision := range reviewRevisions {
		reviewRevisionMap[reviewRevision.ID] = reviewRevision
	}
}

func loadOldReviewReaction() {
	reviewReactions := make([]ReviewReactionV1, 0)
	oldDB.Model(&ReviewReactionV1{}).Where("reaction != 0").Find(&reviewReactions)
	reviewReactionMap = make(map[int64]ReviewReactionV1)
	for _, reviewReaction := range reviewReactions {
		reviewReactionMap[reviewReaction.ID] = reviewReaction
	}
}

func loadNewBaseCourse() {
	baseCourses := make([]po.BaseCoursePO, 0)
	newDB.Model(&po.BaseCoursePO{}).Find(&baseCourses)
	newBaseCourseMap = make(map[string]po.BaseCoursePO)
	for _, baseCourse := range baseCourses {
		newBaseCourseMap[baseCourse.Code] = baseCourse
	}
}

func loadNewCourse() {
	courses := make([]po.CoursePO, 0)
	newDB.Model(&po.CoursePO{}).Find(&courses)
	newCourseMap = make(map[string]po.CoursePO)
	for _, course := range courses {
		newCourseMap[makeCourseKey(course)] = course
	}
}

func loadNewTeacher() {
	teachers := make([]po.TeacherPO, 0)
	newDB.Model(&po.TeacherPO{}).Find(&teachers)
	newTeacherMap = make(map[string]po.TeacherPO)
	for _, teacher := range teachers {
		newTeacherMap[teacher.Code] = teacher
	}
}

func loadNewUser() {
	users := make([]po.UserPO, 0)
	newDB.Model(&po.UserPO{}).Find(&users)
	newUserMap = make(map[int64]po.UserPO)
	for _, user := range users {
		newUserMap[int64(user.ID)] = user
	}
}

func loadNewUserPoint() {
	userPoints := make([]po.UserPointDetailPO, 0)
	newDB.Model(&po.UserPointDetailPO{}).Find(&userPoints)
	newUserPointMap = make(map[int64]po.UserPointDetailPO)
	for _, userPoint := range userPoints {
		newUserPointMap[int64(userPoint.ID)] = userPoint
	}
}

func loadNewReview() {
	reviews := make([]po.ReviewPO, 0)
	newDB.Model(&po.ReviewPO{}).Find(&reviews)
	newReviewMap = make(map[int64]po.ReviewPO)
	for _, review := range reviews {
		newReviewMap[int64(review.ID)] = review
	}
}

func loadNewReviewRevision() {
	reviewRevisions := make([]po.ReviewRevisionPO, 0)
	newDB.Model(&po.ReviewRevisionPO{}).Find(&reviewRevisions)
	newReviewRevisionMap = make(map[int64]po.ReviewRevisionPO)
	for _, reviewRevision := range reviewRevisions {
		newReviewRevisionMap[int64(reviewRevision.ID)] = reviewRevision
	}
}

func loadNewReviewReaction() {
	reviewReactions := make([]po.ReviewReactionPO, 0)
	newDB.Model(&po.ReviewReactionPO{}).Find(&reviewReactions)
	newReviewReactionMap = make(map[int64]po.ReviewReactionPO)
	for _, reviewReaction := range reviewReactions {
		newReviewReactionMap[int64(reviewReaction.ID)] = reviewReaction
	}
}

func BuildNewTeacherFromOld(teacher TeacherV1) po.TeacherPO {
	return po.TeacherPO{
		Code:       teacher.Tid,
		Name:       teacher.Name,
		Title:      teacher.Title,
		Pinyin:     teacher.Pinyin,
		PinyinAbbr: teacher.AbbrPinyin,
		Department: departmentMap[teacher.DepartmentID],
	}
}

func BuildNewCourseFromOld(course CourseV1) (po.BaseCoursePO, po.CoursePO) {
	newBaseCourse := po.BaseCoursePO{
		Code:   course.Code,
		Name:   course.Name,
		Credit: course.Credit,
	}

	teacher := oldToNewTeacherMap[course.MainTeacherID]
	newCourse := po.CoursePO{
		Code:            course.Code,
		Name:            course.Name,
		Credit:          course.Credit,
		MainTeacherID:   int64(teacher.ID),
		MainTeacherName: teacher.Name,
		Department:      departmentMap[course.DepartmentID],
	}
	return newBaseCourse, newCourse
}

func BuildNewUserFromOld(user UserV1) po.UserPO {
	newUser := po.UserPO{
		Model: gorm.Model{
			ID:        uint(user.ID),
			CreatedAt: user.DateJoined,
		},
		Username: user.Username,
		Password: user.Password,
		Email:    user.Username,
	}

	profile := userProfileMap[user.ID]
	newUser.LastSeenAt = profile.LastSeenAt
	newUser.Type = profile.UserType
	return newUser
}

func BuildNewReviewFormOld(review ReviewV1) po.ReviewPO {
	newReview := po.ReviewPO{
		Model: gorm.Model{
			ID:        uint(review.ID),
			CreatedAt: review.CreatedAt,
			UpdatedAt: review.ModifiedAt,
		},
		Comment:     review.Comment,
		Grade:       review.Score,
		Rating:      review.Rating,
		Semester:    semesterMap[review.SemesterID],
		UserID:      review.UserID,
		IsAnonymous: true,
	}
	course := oldToNewCourseMap[review.CourseID]
	newReview.CourseID = int64(course.ID)
	return newReview
}

func BuildNewReviewRevisionFromOld(revision ReviewRevisionV1) po.ReviewRevisionPO {
	return po.ReviewRevisionPO{
		Model: gorm.Model{
			ID:        uint(revision.ID),
			CreatedAt: revision.CreatedAt,
		},
		CourseID: revision.CourseID,
		ReviewID: revision.ReviewID,
		UserID:   revision.UserID,
		Comment:  revision.Comment,
		Rating:   revision.Rating,
		Grade:    revision.Score,
		Semester: semesterMap[revision.SemesterID],
	}
}

func BuildUserPointFromOld(point UserPointV1) po.UserPointDetailPO {
	return po.UserPointDetailPO{
		Model: gorm.Model{
			ID:        uint(point.ID),
			CreatedAt: point.Time,
		},
		UserID:      point.UserID,
		Description: point.Description,
		Value:       point.Value,
	}
}

func BuildNewReviewReactionFromOld(reaction ReviewReactionV1) po.ReviewReactionPO {
	reactionMapping := map[int64]string{
		1:  "like",
		-1: "dislike",
	}
	return po.ReviewReactionPO{
		Model: gorm.Model{
			ID: uint(reaction.ID),
		},
		UserID:   reaction.UserID,
		ReviewID: reaction.ReviewID,
		Reaction: reactionMapping[reaction.Reaction],
	}
}

func main() {
	_ = godotenv.Load()
	oldDB = InitOldDB()
	dal.InitDBClient()
	newDB = dal.GetDBClient()

	loadOldSemester()
	loadOldDepartment()
	loadOldTeacher()
	loadOldCourse()
	loadOldUser()
	loadOldUserProfile()
	loadOldUserPoint()
	loadOldReview()
	loadOldReviewRevision()
	loadOldReviewReaction()

	loadNewBaseCourse()
	loadNewCourse()
	loadNewTeacher()
	loadNewUser()
	loadNewUserPoint()
	loadNewReview()
	loadNewReviewRevision()
	loadNewReviewReaction()

	// course、teacher 如果新的没有，需要添加
	oldToNewTeacherMap = make(map[int64]po.TeacherPO)
	for _, teacher := range teacherMap {
		newTeacher := BuildNewTeacherFromOld(teacher)
		if _, ok := newTeacherMap[newTeacher.Code]; !ok {
			err := newDB.Model(&po.TeacherPO{}).Clauses(clause.OnConflict{DoNothing: true}).Create(&newTeacher).Error
			if err != nil {
				println("create teacher error", newTeacher.Code, newTeacher.Name, err.Error())
				continue
			}
			// println("create teacher ", newTeacher.Name, newTeacher.Code, newTeacher.Department)
			newTeacherMap[newTeacher.Code] = newTeacher
			oldToNewTeacherMap[teacher.ID] = newTeacher
		} else {
			oldToNewTeacherMap[teacher.ID] = newTeacherMap[newTeacher.Code]
		}
	}

	oldToNewCourseMap = make(map[int64]po.CoursePO)
	for _, course := range courseMap {
		newBaseCourse, newCourse := BuildNewCourseFromOld(course)
		if _, ok := newBaseCourseMap[newBaseCourse.Code]; !ok {
			err := newDB.Model(&po.BaseCoursePO{}).Clauses(clause.OnConflict{DoNothing: true}).Create(&newBaseCourse).Error
			if err != nil {
				println("create base course error", newBaseCourse.Code, newBaseCourse.Name, err.Error())
				continue
			}
			// println("created base course ", newBaseCourse.Code, newBaseCourse.Name)
			newBaseCourseMap[newBaseCourse.Code] = newBaseCourse
		}
		if _, ok := newCourseMap[makeCourseKey(newCourse)]; !ok {
			err := newDB.Model(&po.CoursePO{}).Clauses(clause.OnConflict{DoNothing: true}).Create(&newCourse).Error
			if err != nil {
				println("create course error", newCourse.Code, newCourse.Name, err.Error())
				continue
			}
			// println("created course ", newCourse.Code, newCourse.Name, newCourse.MainTeacherID, newCourse.MainTeacherName)
			newCourseMap[makeCourseKey(newCourse)] = newCourse
			oldToNewCourseMap[course.ID] = newCourse
		} else {
			oldToNewCourseMap[course.ID] = newCourseMap[makeCourseKey(newCourse)]
		}
	}

	// user 导入

	for _, user := range userMap {
		newUser := BuildNewUserFromOld(user)
		if _, ok := newUserMap[int64(newUser.ID)]; !ok {
			err := newDB.Model(&po.UserPO{}).Clauses(clause.OnConflict{DoNothing: true}).Create(&newUser).Error
			if err != nil {
				println("create user error", newUser.Username, err.Error())
				continue
			}
			// println("created user ", newUser.ID, newUser.Username)
			newUserMap[int64(newUser.ID)] = newUser
		}
	}

	// user point 导入

	for _, userPoint := range userPointMap {
		newUserPoint := BuildUserPointFromOld(userPoint)
		if _, ok := newUserPointMap[int64(newUserPoint.ID)]; !ok {
			err := newDB.Model(&po.UserPointDetailPO{}).Clauses(clause.OnConflict{DoNothing: true}).Create(&newUserPoint).Error
			if err != nil {
				println("create user point error", newUserPoint.UserID, err.Error())
				continue
			}
			// println("created user point ", newUserPoint.ID, newUserPoint.UserID, newUserPoint.Value)
			newUserPointMap[int64(newUserPoint.ID)] = newUserPoint
		}
	}

	// review 导入

	for _, review := range reviewMap {
		newReview := BuildNewReviewFormOld(review)
		if _, ok := newReviewMap[int64(newReview.ID)]; !ok {
			err := newDB.Model(&po.ReviewPO{}).Create(&newReview).Error
			if err != nil {
				println("create review error", newReview.CourseID, newReview.UserID, err.Error())
				continue
			}
			// println("created review ", newReview.ID, newReview.CourseID, newReview.UserID)
			newReviewMap[int64(newReview.ID)] = newReview
		}
	}

	// review revision 导入

	for _, reviewRevision := range reviewRevisionMap {
		newReviewRevision := BuildNewReviewRevisionFromOld(reviewRevision)
		if _, ok := newReviewRevisionMap[int64(newReviewRevision.ID)]; !ok {
			err := newDB.Model(&po.ReviewRevisionPO{}).Clauses(clause.OnConflict{DoNothing: true}).Create(&newReviewRevision).Error
			if err != nil {
				println("create review revision error", newReviewRevision.CourseID, newReviewRevision.UserID, err.Error())
				continue
			}
			// println("created review revision ", newReviewRevision.ID, newReviewRevision.CourseID, newReviewRevision.UserID, newReviewRevision.CreatedAt.Unix())
			newReviewRevisionMap[int64(newReviewRevision.ID)] = newReviewRevision
		}
	}

	// review reaction 导入

	for _, reviewReaction := range reviewReactionMap {
		newReviewReaction := BuildNewReviewReactionFromOld(reviewReaction)
		if _, ok := newReviewReactionMap[int64(newReviewReaction.ID)]; !ok {
			err := newDB.Model(&po.ReviewReactionPO{}).Clauses(clause.OnConflict{DoNothing: true}).Create(&newReviewReaction).Error
			if err != nil {
				println("create review reaction error", newReviewReaction.ReviewID, newReviewReaction.UserID, err.Error())
				continue
			}
			// println("created review reaction ", newReviewReaction.ID, newReviewReaction.CourseID, newReviewReaction.UserID, newReviewReaction.CreatedAt.Unix())
			newReviewReactionMap[int64(newReviewReaction.ID)] = newReviewReaction
		}
	}
}
