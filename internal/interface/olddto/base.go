package olddto

type BaseResponse struct {
	Message string `json:"message"`
}

type BasePaginateResponse[T any] struct {
	Page     int64 `json:"page"`
	PageSize int64 `json:"page_size"`
	Total    int64 `json:"total"`
	Data     []T   `json:"data"`
}
