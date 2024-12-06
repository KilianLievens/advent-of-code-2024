package advent

import (
	"fmt"
)

type Stringer interface {
	String() string
}

func PrintMatrix[T Stringer](matrix [][]T) {
	for _, line := range matrix {
		for _, r := range line {
			fmt.Printf("%s", r)
		}
		fmt.Printf("\n")
	}
}
