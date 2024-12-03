package main

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/kilianlievens/advent-of-code-2024/advent"
)

func main() {
	exampleOneInput := advent.Read("./input/day_three/example_one.txt")
	fmt.Printf("Example one: mul: %d\n", getMul(exampleOneInput)) // 161
	exampleTwoInput := advent.Read("./input/day_three/example_two.txt")
	fmt.Printf("Example two: conditional mul: %d\n", getConditionalMul(exampleTwoInput)) // 48

	puzzleThreeInput := advent.Read("./input/day_three/puzzle_one.txt")
	fmt.Printf("Puzzle: mul: %d\n", getMul(puzzleThreeInput))                        // 179834255
	fmt.Printf("Puzzle: conditional mul: %d\n", getConditionalMul(puzzleThreeInput)) // 80570939
}

func getMul(input []string) int {
	re := regexp.MustCompile(`mul\((\d+),(\d+)\)`)
	mul := 0

	for _, line := range input {
		matches := re.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			a, _ := strconv.Atoi(match[1])
			b, _ := strconv.Atoi(match[2])
			mul += a * b
		}
	}
	return mul
}

func getConditionalMul(input []string) int {
	re := regexp.MustCompile(`(mul\((\d+),(\d+)\)|don't\(\)|do\(\))`)
	mul := 0
	multiplier := 1

	for _, line := range input {
		matches := re.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			if match[0] == "don't()" {
				multiplier = 0
				continue
			}
			if match[0] == "do()" {
				multiplier = 1
				continue
			}

			a, _ := strconv.Atoi(match[2])
			b, _ := strconv.Atoi(match[3])

			mul += a * b * multiplier
		}
	}

	return mul
}
