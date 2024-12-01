package advent

import "golang.org/x/exp/constraints"

func MinInt[T constraints.Integer](x, y T) T {
	if x > y {
		return y
	}

	return x
}

func MaxInt[T constraints.Integer](x, y T) T {
	if x < y {
		return y
	}

	return x
}
