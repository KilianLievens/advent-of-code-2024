package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/kilianlievens/advent-of-code-2024/advent"
)

func main() {
	exampleOneInput := ReadWithNewline("./input/day_five/example_one.txt")
	sum, incorrectSum := getPageSum(exampleOneInput)
	fmt.Printf("Example: page : %d\n", sum)
	fmt.Printf("Example: incorrect sum: %d\n", incorrectSum)
	puzzleOneInput := ReadWithNewline("./input/day_five/puzzle_one.txt")
	sum, incorrectSum = getPageSum(puzzleOneInput)
	fmt.Printf("Puzzle: sum: %d\n", sum)
	fmt.Printf("Puzzle: incorrectSum sum: %d\n", incorrectSum)
}

func ReadWithNewline(fileName string) []string {
	body, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}

	var lines []string

	for _, line := range strings.Split(string(body), "\n") {
		lines = append(lines, line)
	}

	return lines
}

func parse(input []string) ([][2]int, [][]int) {
	orderingList := [][2]int{}
	updates := [][]int{}

	isOrdering := true
	for _, line := range input {
		if line == "" {
			if isOrdering {
				isOrdering = false
				continue
			}
			break
		}

		if isOrdering {
			parts := strings.Split(line, "|")
			advent.Assert(len(parts) == 2, "Invalid input")
			a, _ := strconv.Atoi(parts[0])
			b, _ := strconv.Atoi(parts[1])
			orderingList = append(orderingList, [2]int{a, b})

			continue
		}

		parts := strings.Split(line, ",")
		update := []int{}
		for _, part := range parts {
			num, _ := strconv.Atoi(part)
			update = append(update, num)
		}
		updates = append(updates, update)
	}

	return orderingList, updates
}

func checkOrdering(orderingList [][2]int, updates []int) bool {
	for i, current := range updates {
		for _, order := range orderingList {
			if order[0] == current {
				for k := i - 1; k >= 0; k-- {
					if updates[k] == order[1] {
						return false
					}
				}
			}

			if order[1] == current {
				for j := i + 1; j < len(updates); j++ {
					if updates[j] == order[0] {
						return false
					}
				}
			}
		}
	}
	return true
}

func createUpdate(input []int, next int, orderingList [][2]int) [][]int {
	possies := [][]int{}
	for i := 0; i <= len(input); i++ {
		c := slices.Clone(input)
		poss := slices.Insert(c, i, next)
		if checkOrdering(orderingList, poss) {
			possies = append(possies, poss)
		}
	}
	return possies
}

func recreateUpdate(update []int, orderingList [][2]int) []int {
	possies := [][]int{}
	for i, item := range update {
		if i == 0 {
			possies = append(possies, []int{item})
			continue
		}

		newPossies := [][]int{}
		for _, poss := range possies {
			newPossies = append(newPossies, createUpdate(poss, item, orderingList)...)
		}

		possies = newPossies
	}

	advent.Assert(len(possies) == 1, "More than one possibility")
	return possies[0]
}

func getPageSum(input []string) (int, int) {
	orderingList, updates := parse(input)
	sum := 0
	incorrectSum := 0

	incorrectlyOrdered := [][]int{}

	for _, update := range updates {
		orderOK := checkOrdering(orderingList, update)

		if orderOK {
			sum += update[int(math.Floor(float64(len(update))/float64(2)))]
			continue
		}

		incorrectlyOrdered = append(incorrectlyOrdered, update)
	}

	for _, update := range incorrectlyOrdered {
		fixed := recreateUpdate(update, orderingList)
		incorrectSum += fixed[int(math.Floor(float64(len(fixed))/float64(2)))]
	}

	return sum, incorrectSum
}
