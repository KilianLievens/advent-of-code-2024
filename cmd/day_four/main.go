package main

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/kilianlievens/advent-of-code-2024/advent"
)

func main() {
	exampleOneInput := advent.Read("./input/day_four/example_one.txt")
	fmt.Printf("Example: XMAS count: %d\n", countXMAS(exampleOneInput))
	fmt.Printf("Example: X-MAS count: %d\n", countCrossMAS(exampleOneInput))
	puzzleOneInput := advent.Read("./input/day_four/puzzle_one.txt")
	fmt.Printf("Puzzle: XMAS count: %d\n", countXMAS(puzzleOneInput))
	fmt.Printf("Puzzle: X-MAS count: %d\n", countCrossMAS(puzzleOneInput))
}

func parse(input []string) [][]string {
	matrix := [][]string{}

	for _, line := range input {
		splitLine := strings.Split(line, "")
		matrix = append(matrix, splitLine)
	}

	advent.Assert(len(matrix) > 1 && len(matrix[0]) > 1, "Matrix too small")

	return matrix
}

func getHorizontalLines(matrix [][]string) [][]string {
	lines := [][]string{}

	for _, line := range matrix {
		// Horizontal left to right
		lines = append(lines, line)

		// Horizontal right to left
		reverseLine := slices.Clone(line)
		slices.Reverse(reverseLine)
		lines = append(lines, reverseLine)
	}

	return lines
}

func diagonalizeNWSE(matrix [][]string) [][]string {
	lines := [][]string{}
	maxCount := len(matrix) + len(matrix[0]) - 1

	for i := 0; i < maxCount; i++ {
		line := []string{}
		for j := 0; j <= i; j++ {
			y := i - j
			x := 0 + j
			if y >= 0 && y < len(matrix) && x >= 0 && x < len(matrix[0]) {
				line = append(line, matrix[y][x])
			}
		}
		lines = append(lines, line)
	}

	return lines
}

func countXMAS(input []string) int {
	matrix := parse(input)
	lines := [][]string{}

	// Horizontal
	lines = append(lines, getHorizontalLines(matrix)...)

	// Diagonal SWNE
	NWSEMatrix := diagonalizeNWSE(matrix)
	lines = append(lines, getHorizontalLines(NWSEMatrix)...)

	rMatrix := advent.RotateRight2D(matrix)

	// Vertical
	lines = append(lines, getHorizontalLines(rMatrix)...)

	// Diagonal NWSE
	SWNEMatrix := diagonalizeNWSE(rMatrix)
	lines = append(lines, getHorizontalLines(SWNEMatrix)...)

	count := 0

	diagonalLines := len(matrix) + len(matrix[0]) - 1
	expectedLineCount := (len(matrix)+len(matrix[0]))*2 + diagonalLines*4
	advent.Assert(len(lines) == expectedLineCount, "Incorrect line count")

	re := regexp.MustCompile(`XMAS`)
	for _, line := range lines {
		matches := re.FindAllStringSubmatch(strings.Join(line, ""), -1)
		count += len(matches)
	}

	return count
}

func isCrossMas(matrix [][]string, y, x int) bool {
	advent.Assert(y >= 0, "y too low")
	advent.Assert(x >= 0, "x too low")
	advent.Assert(y+2 < len(matrix), "y too high")
	advent.Assert(x+2 < len(matrix[0]), "x too high")

	nwse := strings.Join([]string{matrix[y][x], matrix[y+1][x+1], matrix[y+2][x+2]}, "")
	swne := strings.Join([]string{matrix[y+2][x], matrix[y+1][x+1], matrix[y][x+2]}, "")

	if (nwse == "MAS" || nwse == "SAM") && (swne == "MAS" || swne == "SAM") {
		return true
	}

	return false

}

func countCrossMAS(input []string) int {
	matrix := parse(input)
	count := 0
	for y := 0; y < len(matrix)-2; y++ {
		for x := 0; x < len(matrix[0])-2; x++ {
			if isCrossMas(matrix, y, x) {
				count++
			}
		}
	}
	return count
}
