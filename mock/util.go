package mock

import "github.com/pkg/errors"

func GenerateUniqueSet[T comparable](size int, generator func() T) ([]T, error) {
	set := make(map[T]bool)
	for len(set) < size {
		for i := 0; i < MaxLoopLimit+1; i++ {
			value := generator()
			if !set[value] {
				set[value] = true
				break
			}
			if i == MaxLoopLimit {
				return nil, errors.Errorf(ErrTooManyRecords)
			}
		}
	}
	list := make([]T, 0, size)
	for k := range set {
		list = append(list, k)
	}
	if len(list) != size {
		return nil, errors.Errorf("Unexpected len(list) != size")
	}
	return list, nil
}
