package advent

import "golang.org/x/exp/constraints"

func LCM[T constraints.Integer](a, b T, nums ...T) T {
	result := a * b / GCD(a, b)

	for i := 0; i < len(nums); i++ {
		result = LCM(result, nums[i])
	}

	return result
}
