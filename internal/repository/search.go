package repository

import (
	"strings"

	"jcourse_go/pkg/util"
)

// 空格分割的每个词都要匹配，词内分词做模糊匹配
func userQueryToTsQuery(query string) string {
	var sb strings.Builder
	words := strings.Fields(query)
	for i, word := range words {
		if i != 0 {
			sb.WriteString(" & ")
		}
		sb.WriteByte('(')
		segs := util.SegWord(word)
		for j, seg := range segs {
			if j != 0 {
				sb.WriteString(" | ")
			}
			sb.WriteString(seg)
		}
		sb.WriteByte(')')
	}
	return sb.String()
}
