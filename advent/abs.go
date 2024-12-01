package advent

import "golang.org/x/exp/constraints"

func AbsInt[T constraints.Integer](x T) T {
	if x < 0 {
		return -x
	}

	return x
}
