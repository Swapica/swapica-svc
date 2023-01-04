package mem

import (
	"golang.org/x/exp/constraints"
)

func contains[T constraints.Ordered](src []T, value T) bool {
	for _, v := range src {
		if v == value {
			return true
		}
	}

	return false
}
