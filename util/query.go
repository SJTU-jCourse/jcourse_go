package util

func CalcOffset(page int64, pageSize int64) int64 {
	return (page - 1) * pageSize
}
