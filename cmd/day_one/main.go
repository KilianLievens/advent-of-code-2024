package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/kilianlievens/advent-of-code-2024/advent"
)

func main() {
	exampleOneInput := advent.Read("./input/day_one/example_one.txt")
	fmt.Printf("Example: distance: %d\n", getDistance(exampleOneInput))     // 11
	fmt.Printf("Example: similarity: %d\n", getSimilarity(exampleOneInput)) // 31
	puzzleOneInput := advent.Read("./input/day_one/puzzle_one.txt")
	fmt.Printf("Puzzle: distance: %d\n", getDistance(puzzleOneInput))     // 1970720
	fmt.Printf("Puzzle: similarity: %d\n", getSimilarity(puzzleOneInput)) // 17191599
}

func parse(input []string) ([]int, []int) {
	listOne := []int{}
	listTwo := []int{}
	for _, line := range input {
		splitLine := strings.Split(line, "   ")

		advent.Assert(len(splitLine) == 2, "Line should contain exactly two numbers")

		numberOne, _ := strconv.Atoi(splitLine[0])
		numberTwo, _ := strconv.Atoi(splitLine[1])

		listOne = append(listOne, numberOne)
		listTwo = append(listTwo, numberTwo)
	}

	advent.Assert(len(listOne) == len(listTwo), "Lists are not the same length")

	return listOne, listTwo
}

func getDistance(input []string) int {
	listOne, listTwo := parse(input)

	slices.Sort(listOne)
	slices.Sort(listTwo)

	distance := 0
	for i, number := range listOne {
		distance += advent.AbsInt(number - listTwo[i])
	}

	return distance
}

func getSimilarity(input []string) int {
	listOne, listTwo := parse(input)

	counts := make(map[int]int)
	for _, number := range listTwo {
		counts[number]++
	}

	similarity := 0
	for _, number := range listOne {
		similarity += number * counts[number]
	}

	return similarity
}
