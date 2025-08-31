package user

import "jcourse_go/internal/domain/shared"

type UserFilterForQuery struct {
	shared.PaginationFilterForQuery
}

type UserPointDetailFilter struct {
	shared.PaginationFilterForQuery
	UserPointDetailID int64
	UserID            int64
	StartTime         int64
	EndTime           int64
}
