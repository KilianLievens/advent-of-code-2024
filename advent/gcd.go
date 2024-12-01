package advent

import "golang.org/x/exp/constraints"

func GCD[T constraints.Integer](a, b T) T {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}

	return a
}
