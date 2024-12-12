package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kilianlievens/advent-of-code-2024/advent"
)

func main() {
	exampleOneInput := advent.Read("./input/day_eleven/example_one.txt")
	exampleTwoInput := advent.Read("./input/day_eleven/example_two.txt")
	fmt.Printf("Example 1, single blink: stone count: %d\n", countStoneRow(exampleOneInput, 1))
	fmt.Printf("Example 2, six blinks: stone count: %d\n", countStoneRow(exampleTwoInput, 6))
	fmt.Printf("Example 2, 22 blinks: stone count: %d\n", countStoneRow(exampleTwoInput, 22))
	fmt.Printf("Example 2, 25 blinks: stone count: %d\n", countStoneRow(exampleTwoInput, 25))
	puzzleOneInput := advent.Read("./input/day_eleven/puzzle_one.txt")
	fmt.Printf("Puzzle, 25 blinks: stone count: %d\n", countStoneRow(puzzleOneInput, 25))
	fmt.Printf("Puzzle, 75 blinks: stone count: %d\n", countStoneRow(puzzleOneInput, 75))
}

func parse(input []string) []stone {
	advent.Assert(len(input) == 1, "Input should be one line")

	row := []stone{}

	for _, line := range input {
		splitLine := strings.Split(line, " ")
		for _, num := range splitLine {
			num, _ := strconv.Atoi(num)
			row = append(row, stone(num))
		}
	}

	return row
}

type stone int

func (s stone) countDigits() int {
	str := fmt.Sprintf("%d", s)
	if s < 0 {
		return len(str) - 1
	}
	return len(str)
}

func (s stone) split() (stone, stone) {
	digitCount := s.countDigits()

	advent.Assert(digitCount > 0, "digit count must more than 0")
	advent.Assert(digitCount%2 == 0, "digit count must be even")

	str := fmt.Sprintf("%d", s)
	half := digitCount / 2
	first, _ := strconv.Atoi(str[:half])
	second, _ := strconv.Atoi(str[half:])

	return stone(first), stone(second)
}

func blink(row []stone) []stone {
	newRow := []stone{}
	for _, s := range row {
		if s == 0 {
			newRow = append(newRow, 1)
			continue
		}

		if s.countDigits()%2 == 0 {
			first, second := s.split()
			newRow = append(newRow, first, second)
			continue
		}

		newRow = append(newRow, s*2024)
	}

	return newRow
}

type predictState struct {
	num    int
	blinks uint
}

var predictCache map[predictState]uint = make(map[predictState]uint)

func countStonesDynamic(s stone, blinks uint) uint {
	cacheKey := predictState{int(s), blinks}

	if val, ok := predictCache[cacheKey]; ok {
		return val
	}

	res := blink([]stone{s})
	if blinks == 1 {
		predictCache[cacheKey] = uint(len(res))
		return uint(len(res))
	}

	var sum uint = 0
	for _, num := range res {
		sum += countStonesDynamic(num, blinks-1)
	}

	predictCache[cacheKey] = sum
	return sum
}

func countStoneRow(input []string, blinks uint) uint {
	row := parse(input)

	var sum uint = 0

	for _, s := range row {
		sum += countStonesDynamic(s, blinks)
	}

	return sum
}
