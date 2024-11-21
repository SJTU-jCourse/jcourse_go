package mock

import (
	"jcourse_go/model/po"
	"strconv"

	"gorm.io/gorm/clause"

	"syreclabs.com/go/faker"
)

func MockReviews(gen MockDBGenerator, n int, courses []po.CoursePO, users []po.UserPO) ([]po.ReviewPO, error) {
	reviews := make([]po.ReviewPO, n)
	type courseUserPair struct {
		courseID int64
		userID   int64
	}
	courseUserPairs, err := GenerateUniqueSet(n, func() courseUserPair {
		courseID := int64(courses[gen.Rand.Intn(len(courses))].ID)
		userID := int64(users[gen.Rand.Intn(len(users))].ID)
		return courseUserPair{userID: userID, courseID: courseID}
	})
	if err != nil {
		return nil, err
	}
	for i := 0; i < n; i++ {
		grade := gen.Rand.Intn(5) + 2020
		reviews[i] = po.ReviewPO{
			CourseID:    courseUserPairs[i].courseID,
			UserID:      courseUserPairs[i].userID,
			Comment:     faker.Lorem().Paragraph(gen.Rand.Intn(5) + 1),
			Rating:      int64(gen.Rand.Intn(20) + 80),
			Semester:    strconv.Itoa(grade) + "-" + strconv.Itoa(gen.Rand.Intn(3)+1),
			IsAnonymous: gen.Rand.Intn(2) == 0,
			Grade:       strconv.Itoa(grade),
		}
	}

	err = gen.db.Clauses(clause.OnConflict{
		DoNothing: true,
	}).CreateInBatches(reviews, gen.batchSize).Error
	return reviews, err
}
