package main

import (
	"context"
	"fmt"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"jcourse_go/internal/dal"
	"jcourse_go/internal/model/converter"
	po2 "jcourse_go/internal/model/po"
	"jcourse_go/internal/model/types"
	"jcourse_go/internal/repository"
	"jcourse_go/internal/service"
	"jcourse_go/pkg/util"
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

	oldToNewCourseMap  map[int64]po2.CoursePO  // old id -> new model
	oldToNewTeacherMap map[int64]po2.TeacherPO // old id -> new model

	newBaseCourseMap     map[string]po2.BaseCoursePO     // code -> new
	newCourseMap         map[string]po2.CoursePO         // course key -> new
	newTeacherMap        map[string]po2.TeacherPO        // teacher code -> new
	newUserMap           map[int64]po2.UserPO            // uid -> new
	newReviewMap         map[int64]po2.ReviewPO          // uid -> new
	newReviewRevisionMap map[int64]po2.ReviewRevisionPO  // uid -> new
	newUserPointMap      map[int64]po2.UserPointDetailPO // uid -> new
	newReviewReactionMap map[int64]po2.ReviewReactionPO  // uid -> new
	newRatingMap         map[string]po2.RatingPO         // uid & course_id -> new rating
)

func makeCourseKey(course po2.CoursePO) string {
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

func makeRatingKey(userID int64, courseID int64) string {
	return fmt.Sprintf("%d:%d", userID, courseID)
}

func loadNewBaseCourse() {
	baseCourses := make([]po2.BaseCoursePO, 0)
	newDB.Model(&po2.BaseCoursePO{}).Find(&baseCourses)
	newBaseCourseMap = make(map[string]po2.BaseCoursePO)
	for _, baseCourse := range baseCourses {
		newBaseCourseMap[baseCourse.Code] = baseCourse
	}
}

func loadNewCourse() {
	courses := make([]po2.CoursePO, 0)
	newDB.Model(&po2.CoursePO{}).Find(&courses)
	newCourseMap = make(map[string]po2.CoursePO)
	for _, course := range courses {
		newCourseMap[makeCourseKey(course)] = course
	}
}

func loadNewTeacher() {
	teachers := make([]po2.TeacherPO, 0)
	newDB.Model(&po2.TeacherPO{}).Find(&teachers)
	newTeacherMap = make(map[string]po2.TeacherPO)
	for _, teacher := range teachers {
		newTeacherMap[teacher.Code] = teacher
	}
}

func loadNewUser() {
	users := make([]po2.UserPO, 0)
	newDB.Model(&po2.UserPO{}).Find(&users)
	newUserMap = make(map[int64]po2.UserPO)
	for _, user := range users {
		newUserMap[int64(user.ID)] = user
	}
}

func loadNewUserPoint() {
	userPoints := make([]po2.UserPointDetailPO, 0)
	newDB.Model(&po2.UserPointDetailPO{}).Find(&userPoints)
	newUserPointMap = make(map[int64]po2.UserPointDetailPO)
	for _, userPoint := range userPoints {
		newUserPointMap[int64(userPoint.ID)] = userPoint
	}
}

func loadNewReview() {
	reviews := make([]po2.ReviewPO, 0)
	newDB.Model(&po2.ReviewPO{}).Find(&reviews)
	newReviewMap = make(map[int64]po2.ReviewPO)
	for _, review := range reviews {
		newReviewMap[int64(review.ID)] = review
	}
}

func loadNewReviewRevision() {
	reviewRevisions := make([]po2.ReviewRevisionPO, 0)
	newDB.Model(&po2.ReviewRevisionPO{}).Find(&reviewRevisions)
	newReviewRevisionMap = make(map[int64]po2.ReviewRevisionPO)
	for _, reviewRevision := range reviewRevisions {
		newReviewRevisionMap[int64(reviewRevision.ID)] = reviewRevision
	}
}

func loadNewReviewReaction() {
	reviewReactions := make([]po2.ReviewReactionPO, 0)
	newDB.Model(&po2.ReviewReactionPO{}).Find(&reviewReactions)
	newReviewReactionMap = make(map[int64]po2.ReviewReactionPO)
	for _, reviewReaction := range reviewReactions {
		newReviewReactionMap[int64(reviewReaction.ID)] = reviewReaction
	}
}

func loadNewRating() {
	ratings := make([]po2.RatingPO, 0)
	newDB.Model(&po2.RatingPO{}).Where("related_type = ?", "course").Find(&ratings)
	newRatingMap = make(map[string]po2.RatingPO)
	for _, rating := range ratings {
		newRatingMap[makeRatingKey(rating.UserID, rating.RelatedID)] = rating
	}
}

func BuildNewTeacherFromOld(teacher TeacherV1) po2.TeacherPO {
	return po2.TeacherPO{
		Code:       teacher.Tid,
		Name:       teacher.Name,
		Title:      teacher.Title,
		Pinyin:     teacher.Pinyin,
		PinyinAbbr: teacher.AbbrPinyin,
		Department: departmentMap[teacher.DepartmentID],
	}
}

func BuildNewCourseFromOld(course CourseV1) (po2.BaseCoursePO, po2.CoursePO) {
	newBaseCourse := po2.BaseCoursePO{
		Code:   course.Code,
		Name:   course.Name,
		Credit: course.Credit,
	}

	teacher := oldToNewTeacherMap[course.MainTeacherID]
	newCourse := po2.CoursePO{
		Code:            course.Code,
		Name:            course.Name,
		Credit:          course.Credit,
		MainTeacherID:   int64(teacher.ID),
		MainTeacherName: teacher.Name,
		Department:      departmentMap[course.DepartmentID],
	}
	return newBaseCourse, newCourse
}

func BuildNewUserFromOld(user UserV1) po2.UserPO {
	newUser := po2.UserPO{
		ID:        user.ID,
		CreatedAt: user.DateJoined,
		Username:  user.Username,
		Password:  user.Password,
		Email:     user.Username,
	}

	profile := userProfileMap[user.ID]
	newUser.LastSeenAt = profile.LastSeenAt
	newUser.Type = profile.UserType
	return newUser
}

func BuildNewReviewFormOld(review ReviewV1) po2.ReviewPO {
	newReview := po2.ReviewPO{
		ID:          review.ID,
		CreatedAt:   review.CreatedAt,
		UpdatedAt:   review.ModifiedAt,
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

func BuildNewReviewRevisionFromOld(revision ReviewRevisionV1) po2.ReviewRevisionPO {
	newCourse := oldToNewCourseMap[revision.CourseID]

	return po2.ReviewRevisionPO{
		ID:        revision.ID,
		CreatedAt: revision.CreatedAt,
		CourseID:  newCourse.ID,
		ReviewID:  revision.ReviewID,
		UserID:    revision.UserID,
		Comment:   revision.Comment,
		Rating:    revision.Rating,
		Grade:     revision.Score,
		Semester:  semesterMap[revision.SemesterID],
	}
}

func BuildUserPointFromOld(point UserPointV1) po2.UserPointDetailPO {
	return po2.UserPointDetailPO{
		ID:          point.ID,
		CreatedAt:   point.Time,
		UserID:      point.UserID,
		Description: point.Description,
		Value:       point.Value,
	}
}

func BuildNewReviewReactionFromOld(reaction ReviewReactionV1) po2.ReviewReactionPO {
	reactionMapping := map[int64]string{
		1:  "+1",
		-1: "-1",
	}
	return po2.ReviewReactionPO{
		ID:       reaction.ID,
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
	repository.SetDefault(newDB)

	println("loading old")
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

	println("loading new")
	loadNewBaseCourse()
	loadNewCourse()
	loadNewTeacher()
	loadNewUser()
	loadNewUserPoint()
	loadNewReview()
	loadNewReviewRevision()
	loadNewReviewReaction()

	println("start import")
	println("importing teacher")
	// course、teacher 如果新的没有，需要添加
	oldToNewTeacherMap = make(map[int64]po2.TeacherPO)
	for _, teacher := range teacherMap {
		newTeacher := BuildNewTeacherFromOld(teacher)
		if _, ok := newTeacherMap[newTeacher.Code]; !ok {
			err := newDB.Model(&po2.TeacherPO{}).Clauses(clause.OnConflict{DoNothing: true}).Create(&newTeacher).Error
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

	newDB.Exec("SELECT setval('teachers_id_seq', (SELECT MAX(id) FROM teachers));")

	println("importing course")
	oldToNewCourseMap = make(map[int64]po2.CoursePO)
	for _, course := range courseMap {
		newBaseCourse, newCourse := BuildNewCourseFromOld(course)
		if _, ok := newBaseCourseMap[newBaseCourse.Code]; !ok {
			err := newDB.Model(&po2.BaseCoursePO{}).Clauses(clause.OnConflict{DoNothing: true}).Create(&newBaseCourse).Error
			if err != nil {
				println("create base course error", newBaseCourse.Code, newBaseCourse.Name, err.Error())
				continue
			}
			// println("created base course ", newBaseCourse.Code, newBaseCourse.Name)
			newBaseCourseMap[newBaseCourse.Code] = newBaseCourse
		}
		if _, ok := newCourseMap[makeCourseKey(newCourse)]; !ok {
			err := newDB.Model(&po2.CoursePO{}).Clauses(clause.OnConflict{DoNothing: true}).Create(&newCourse).Error
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

	newDB.Exec("SELECT setval('courses_id_seq', (SELECT MAX(id) FROM courses));")

	// user 导入
	println("importing user")
	for _, user := range userMap {
		newUser := BuildNewUserFromOld(user)
		if _, ok := newUserMap[int64(newUser.ID)]; !ok {
			err := newDB.Model(&po2.UserPO{}).Clauses(clause.OnConflict{DoNothing: true}).Create(&newUser).Error
			if err != nil {
				println("create user error", newUser.Username, err.Error())
				continue
			}
			// println("created user ", newUser.ID, newUser.Username)
			newUserMap[int64(newUser.ID)] = newUser
		}
	}

	newDB.Exec("SELECT setval('users_id_seq', (SELECT MAX(id) FROM users));")

	// user point 导入
	println("importing user point")
	for _, userPoint := range userPointMap {
		newUserPoint := BuildUserPointFromOld(userPoint)
		if _, ok := newUserPointMap[int64(newUserPoint.ID)]; !ok {
			err := newDB.Model(&po2.UserPointDetailPO{}).Clauses(clause.OnConflict{DoNothing: true}).Create(&newUserPoint).Error
			if err != nil {
				println("create user point error", newUserPoint.UserID, err.Error())
				continue
			}
			// println("created user point ", newUserPoint.ID, newUserPoint.UserID, newUserPoint.Value)
			newUserPointMap[int64(newUserPoint.ID)] = newUserPoint
		}
	}

	newDB.Exec("SELECT setval('user_points_id_seq', (SELECT MAX(id) FROM user_points));")

	// review 导入
	println("importing review")
	for _, review := range reviewMap {
		newReview := BuildNewReviewFormOld(review)
		if _, ok := newReviewMap[int64(newReview.ID)]; !ok {
			err := newDB.Model(&po2.ReviewPO{}).Create(&newReview).Error
			if err != nil {
				println("create review error", newReview.CourseID, newReview.UserID, err.Error())
				continue
			}
			// println("created review ", newReview.ID, newReview.CourseID, newReview.UserID)
			newReviewMap[int64(newReview.ID)] = newReview
		}
	}

	newDB.Exec("SELECT setval('reviews_id_seq', (SELECT MAX(id) FROM reviews));")

	// review revision 导入
	println("importing review revision")
	for _, reviewRevision := range reviewRevisionMap {
		newReviewRevision := BuildNewReviewRevisionFromOld(reviewRevision)
		if _, ok := newReviewRevisionMap[int64(newReviewRevision.ID)]; !ok {
			err := newDB.Model(&po2.ReviewRevisionPO{}).Clauses(clause.OnConflict{DoNothing: true}).Create(&newReviewRevision).Error
			if err != nil {
				println("create review revision error", newReviewRevision.CourseID, newReviewRevision.UserID, err.Error())
				continue
			}
			// println("created review revision ", newReviewRevision.ID, newReviewRevision.CourseID, newReviewRevision.UserID, newReviewRevision.CreatedAt.Unix())
			newReviewRevisionMap[int64(newReviewRevision.ID)] = newReviewRevision
		}
	}

	newDB.Exec("SELECT setval('review_revisions_id_seq', (SELECT MAX(id) FROM review_revisions));")

	// review reaction 导入
	println("importing review reaction")
	for _, reviewReaction := range reviewReactionMap {
		newReviewReaction := BuildNewReviewReactionFromOld(reviewReaction)
		if _, ok := newReviewReactionMap[int64(newReviewReaction.ID)]; !ok {
			err := newDB.Model(&po2.ReviewReactionPO{}).Clauses(clause.OnConflict{DoNothing: true}).Create(&newReviewReaction).Error
			if err != nil {
				println("create review reaction error", newReviewReaction.ReviewID, newReviewReaction.UserID, err.Error())
				continue
			}
			// println("created review reaction ", newReviewReaction.ID, newReviewReaction.CourseID, newReviewReaction.UserID, newReviewReaction.CreatedAt.Unix())
			newReviewReactionMap[int64(newReviewReaction.ID)] = newReviewReaction
		}
	}

	newDB.Exec("SELECT setval('review_reactions_id_seq', (SELECT MAX(id) FROM review_reactions));")

	loadNewRating()
	for _, review := range newReviewMap {
		ratingKey := makeRatingKey(review.UserID, review.CourseID)
		if _, ok := newRatingMap[ratingKey]; !ok {
			rating := converter.BuildRatingFromReview(review)
			err := newDB.Model(&po2.RatingPO{}).Clauses(clause.OnConflict{DoNothing: true}).Create(&rating).Error
			if err != nil {
				println("create rating error", rating.RelatedID, rating.UserID, err.Error())
				continue
			}
			newRatingMap[ratingKey] = rating
		}
	}

	needToSyncCourseID := make(map[int64]struct{})
	for _, rating := range newRatingMap {
		needToSyncCourseID[rating.RelatedID] = struct{}{}
	}

	for courseID, _ := range needToSyncCourseID {
		err := service.SyncRating(context.Background(), types.RelatedTypeCourse, courseID)
		if err != nil {
			println("sync rating error for course id = ", courseID, err.Error())
			continue
		}
	}
}
