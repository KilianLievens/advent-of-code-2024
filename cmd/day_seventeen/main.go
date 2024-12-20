package main

import (
	"fmt"
	"math"
)

func main() {
	exampleOneInput := read("./input/day_seventeen/example_one.txt")
	exampleTwoInput := read("./input/day_seventeen/example_two.txt")

	fmt.Println("Example: Run")
	runThingamabob(exampleOneInput)
	fmt.Println("Example: Quine")
	findThingamabobQuine(exampleTwoInput)

	puzzleOneInput := read("./input/day_seventeen/puzzle_one.txt")
	fmt.Println("Puzzle One")
	runThingamabob(puzzleOneInput)
	// fmt.Println("Puzzle: Quine")
	// findThingamabobQuine(puzzleOneInput)
}

// TODO: main function and file naming
func runThingamabob(input []string) {
	registerMap, program := parse(input)
	output := program.execute(&registerMap)
	fmt.Println(output)
}

func findThingamabobQuine(input []string) {
	_, program := parse(input)
	programString := program.String()
	for i := 0; i < math.MaxInt32; i++ {
		registerMap := registers{a: i, b: 0, c: 0}
		output := program.execute(&registerMap)

		if output.String() == programString {
			fmt.Println(i)
			break
		}
	}
}
