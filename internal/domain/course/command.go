package course

import "jcourse_go/internal/domain/shared"

type WriteReviewCommand struct {
	CourseID shared.IDType
	Comment  string
	Rating   int64
	Semester string
	Score    string
}

type UpdateReviewCommand struct {
	ReviewID shared.IDType
	Comment  string
	Rating   int64
	Semester string
	Score    string
}

type DeleteReviewCommand struct {
	ReviewID shared.IDType
}
