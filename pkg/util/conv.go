package util

import "strconv"

func StringsToFloat64s(str []string) []float64 {
	res := make([]float64, 0, len(str))
	for _, s := range str {
		val, err := strconv.ParseFloat(s, 64)
		if err != nil {
			panic(err)
		}
		res = append(res, val)
	}
	return res
}
