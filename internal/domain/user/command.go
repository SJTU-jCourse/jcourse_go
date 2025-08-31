package user

import "jcourse_go/internal/domain/shared"

type UserFilterForQuery struct {
	shared.PaginationQuery
}

type UserPointDetailFilter struct {
	shared.PaginationQuery
	UserPointDetailID int64
	UserID            int64
	StartTime         int64
	EndTime           int64
}
