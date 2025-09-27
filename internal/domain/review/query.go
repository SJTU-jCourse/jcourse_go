package review

import "jcourse_go/internal/domain/shared"

type ReviewQuery struct {
	shared.PaginationQuery
	Semester string `json:"semester" form:"semester"`
	Rating   int64  `json:"rating" form:"rating"`
}
