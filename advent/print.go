package advent

import (
	"fmt"
)

func PrintMatrix[T fmt.Stringer](matrix [][]T) {
	for _, line := range matrix {
		for _, r := range line {
			fmt.Printf("%s", r)
		}
		fmt.Printf("\n")
	}
}
