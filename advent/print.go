package advent

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

func PrintRuneMatrix(matrix [][]rune) {
	for _, line := range matrix {
		for _, r := range line {
			fmt.Printf("%s", string(r))
		}
		fmt.Printf("\n")
	}
}

func PrintStringMatrix(matrix [][]string) {
	for _, line := range matrix {
		for _, r := range line {
			fmt.Printf("%s", r)
		}
		fmt.Printf("\n")
	}
}

func PrintIntMatrix[T constraints.Integer](matrix [][]T) {
	for _, line := range matrix {
		for _, r := range line {
			fmt.Printf("%d", r)
		}
		fmt.Printf("\n")
	}
}
