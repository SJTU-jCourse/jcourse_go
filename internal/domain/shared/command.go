package shared

type PaginationFilterForQuery struct {
	Page      int64  `json:"page" form:"page"`
	PageSize  int64  `json:"page_size" form:"page_size"`
	Search    string `json:"search" form:"search"`
	Order     string `json:"order" form:"order"`
	Ascending bool   `json:"ascending" form:"ascending"`
}
