package mock

import (
	"math/rand"

	"syreclabs.com/go/faker"
)

// SetSeed 生成确定的随机数种子
func SetSeed(seed int64) rand.Source {
	faker.Seed(seed)
	return rand.NewSource(seed)
}

// ErrTooManyRecords 生成数量过多让一些唯一性约束可能无法满足
const ErrTooManyRecords = "too many records"
const MaxLoopLimit = 1000 // 防止无限循环, 1000次随机生成失败则报错
