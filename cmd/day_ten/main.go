package main

import (
	"fmt"
	"strconv"

	"github.com/kilianlievens/advent-of-code-2024/advent"
)

func main() {
	exampleOneInput := advent.Read("./input/day_ten/example_one.txt")
	score, rating := countTrailheadScore(exampleOneInput)
	fmt.Printf("Example 1: Trailhead Score: %d\n", score)
	fmt.Printf("Example 1: Trailhead Rating: %d\n", rating)
	exampleTwoInput := advent.Read("./input/day_ten/example_two.txt")
	score, rating = countTrailheadScore(exampleTwoInput)
	fmt.Printf("Example 2: Trailhead Score: %d\n", score)
	fmt.Printf("Example 2: Trailhead Rating: %d\n", rating)
	puzzleOneInput := advent.Read("./input/day_ten/puzzle_one.txt")
	score, rating = countTrailheadScore(puzzleOneInput)
	fmt.Printf("Puzzle: Trailhead Score: %d\n", score)
	fmt.Printf("Puzzle: Trailhead Rating: %d\n", rating)
}

// ------------
type matrix [][]int

func (m matrix) get(pos pos) int {
	return m[pos.y][pos.x]
}

func (m matrix) inBounds(pos pos) bool {
	return pos.x >= 0 && pos.y >= 0 && pos.x < len(m[0]) && pos.y < len(m)
}

// ------------
type pos struct {
	x, y int
}

func (p pos) getNeighbours() []pos {
	return []pos{
		{p.x - 1, p.y},
		{p.x + 1, p.y},
		{p.x, p.y - 1},
		{p.x, p.y + 1},
	}
}

// ------------
func parse(input []string) matrix {
	matrix := matrix{}

	for _, line := range input {
		newLine := []int{}
		for _, char := range line {
			num, err := strconv.Atoi(string(char))
			advent.Assert(err == nil, "Could not convert character to int")
			newLine = append(newLine, num)
		}
		matrix = append(matrix, newLine)
	}

	advent.Assert(len(matrix) == len(input) && len(matrix[0]) == len(input[0]), "Matrix is not the same size as input")

	return matrix
}

func path(matrix matrix, pos pos, beenPos map[pos]bool) int {
	val := matrix.get(pos)

	if beenPos[pos] {
		return 0
	}

	beenPos[pos] = true

	if val == 9 {
		return 1
	}

	neighbours := pos.getNeighbours()
	count := 0

	for _, neighbour := range neighbours {
		if !matrix.inBounds(neighbour) {
			continue
		}
		neighVal := matrix.get(neighbour)
		if neighVal == val+1 {
			count += path(matrix, neighbour, beenPos)
		}
	}

	return count
}

func pathRating(matrix matrix, pos pos) int {
	val := matrix.get(pos)

	if val == 9 {
		return 1
	}

	neighbours := pos.getNeighbours()
	count := 0

	for _, neighbour := range neighbours {
		if !matrix.inBounds(neighbour) {
			continue
		}
		neighVal := matrix.get(neighbour)
		if neighVal == val+1 {
			count += pathRating(matrix, neighbour)
		}
	}

	return count
}

func countTrailheadScore(input []string) (int, int) {
	matrix := parse(input)

	sum := 0
	for y, row := range matrix {
		for x, item := range row {
			if item == 0 {
				sum += path(matrix, pos{x, y}, map[pos]bool{})
			}
		}
	}

	ratingSum := 0
	for y, row := range matrix {
		for x, item := range row {
			if item == 0 {
				ratingSum += pathRating(matrix, pos{x, y})
			}
		}
	}

	return sum, ratingSum
}
