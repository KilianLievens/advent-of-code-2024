package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kilianlievens/advent-of-code-2024/advent"
)

func main() {
	exampleOneInput := advent.Read("./input/day_seven/example_one.txt")
	count, concatCount := calibrate(exampleOneInput)
	fmt.Printf("Example: calibration result: %d\n", count)
	fmt.Printf("Example: calibration result with concatenation: %d\n", concatCount)
	puzzleOneInput := advent.Read("./input/day_seven/puzzle_one.txt")
	count, concatCount = calibrate(puzzleOneInput)
	fmt.Printf("Puzzle: calibration result: %d\n", count)
	fmt.Printf("Puzzle: calibration result with concatenation: %d\n", concatCount)
}

type equation struct {
	result     int
	components []int
}

func (e *equation) check(useConcatenation bool) int {
	advent.Assert(len(e.components) > 0, "Equation has no components")
	if e.components[0] == e.result && len(e.components) == 1 {
		return 1
	}

	if e.components[0] > e.result {
		return 0
	}

	counter := 0
	intermediateResults := []int{e.components[0]}

	for i := 1; i < len(e.components); i++ {
		newIntermediateResults := []int{}
		for _, intermediateResult := range intermediateResults {
			multiplied := intermediateResult * e.components[i]
			added := intermediateResult + e.components[i]
			concatenated, _ := strconv.Atoi(fmt.Sprintf("%d%d", intermediateResult, e.components[i]))

			toCheck := []int{multiplied, added}
			if useConcatenation {
				toCheck = append(toCheck, concatenated)
			}

			for _, newResult := range toCheck {
				if newResult == e.result && i == len(e.components)-1 {
					counter++
					continue
				}
				if newResult <= e.result {
					newIntermediateResults = append(newIntermediateResults, newResult)
				}
			}
		}

		intermediateResults = newIntermediateResults
	}

	return counter
}

func parse(input []string) []equation {
	equations := make([]equation, len(input))

	for i, line := range input {
		splitLine := strings.Split(line, ": ")
		advent.Assert(len(splitLine) == 2, "Invalid input")

		result, _ := strconv.Atoi(splitLine[0])

		componentsSplit := strings.Split(splitLine[1], " ")
		components := make([]int, len(componentsSplit))
		for j, component := range componentsSplit {
			c, _ := strconv.Atoi(component)
			components[j] = c
		}

		equations[i] = equation{result, components}
	}

	advent.Assert(len(equations) == len(input), "Not all equations parsed")

	return equations
}

func calibrate(input []string) (int, int) {
	equations := parse(input)

	sum := 0
	concatSum := 0

	for _, e := range equations {
		count := e.check(false)
		countConcat := e.check(true)
		if count > 0 {
			sum += e.result
		}
		if countConcat > 0 {
			concatSum += e.result
		}
	}

	return sum, concatSum
}
