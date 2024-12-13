package main

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/kilianlievens/advent-of-code-2024/advent"
)

func main() {
	exampleOneInput := advent.Read("./input/day_thirteen/example_one.txt")
	fmt.Printf("Example: claw cost: %d\n", startMachines(exampleOneInput))
	fmt.Printf("Example: claw cost but also far: %d\n", startMachinesLA(exampleOneInput))
	puzzleOneInput := advent.Read("./input/day_thirteen/puzzle_one.txt")
	fmt.Printf("Puzzle: claw cost: %d\n", startMachines(puzzleOneInput))
	fmt.Printf("Puzzle: claw cost but also far: %d\n", startMachinesLA(puzzleOneInput))
}

type pos struct {
	x, y int
}

func (p pos) move(p2 pos, inverted bool) pos {
	if inverted {
		return pos{p.x - p2.x, p.y - p2.y}
	}

	return pos{p.x + p2.x, p.y + p2.y}
}

func (p pos) multiply(factor int) pos {
	return pos{p.x * factor, p.y * factor}
}

func (p pos) same(p2 pos) bool {
	return p.x == p2.x && p.y == p2.y
}

type button struct {
	movement pos
	cost     int
}

type machine struct {
	destination      pos
	aButton, bButton button
}

func (m machine) String() string {
	return fmt.Sprintf("(Dest: %v, A: %v, B: %v", m.destination, m.aButton, m.bButton)
}

func (m machine) pressButton(pos pos, button rune, inverted bool) (pos, int) {
	advent.Assert(button == 'A' || button == 'B', "Invalid button")

	if button == 'A' {
		return pos.move(m.aButton.movement, inverted), m.aButton.cost
	}

	return pos.move(m.bButton.movement, inverted), m.bButton.cost
}

func parse(input []string, aCost, bCost int, addToDest int) []machine {
	machines := []machine{}

	movementA := regexp.MustCompile(`X\+(\d+)`)
	movementB := regexp.MustCompile(`Y\+(\d+)`)
	prizeX := regexp.MustCompile(`X=(\d+)`)
	prizeY := regexp.MustCompile(`Y=(\d+)`)

	counter := 0
	curMachine := machine{}
	for _, line := range input {
		mod := counter % 3

		switch mod {
		case 0, 1:
			rawA := movementA.FindStringSubmatch(line)
			movA, _ := strconv.Atoi(rawA[1])
			rawB := movementB.FindStringSubmatch(line)
			movB, _ := strconv.Atoi(rawB[1])

			if mod == 0 {
				curMachine.aButton = button{pos{movA, movB}, aCost}
			} else {
				curMachine.bButton = button{pos{movA, movB}, bCost}
			}

		case 2:
			rawX := prizeX.FindStringSubmatch(line)
			rawY := prizeY.FindStringSubmatch(line)
			x, _ := strconv.Atoi(rawX[1])
			y, _ := strconv.Atoi(rawY[1])
			curMachine.destination = pos{x + addToDest, y + addToDest}

			machines = append(machines, curMachine)
			curMachine = machine{}
		}

		counter++
	}

	return machines
}

type point struct {
	pos  pos
	cost int
}

func playClaw(m machine) int {
	bLine := []point{}
	aLine := []point{}

	for i := 0; i <= 100; i++ {
		b := m.destination.move(m.bButton.movement.multiply(i), true)
		if b.x < 0 || b.y < 0 {
			break
		}

		point := point{b, m.bButton.cost * i}

		if b.x == 0 && b.y == 0 {
			return point.cost
		}

		bLine = append(bLine, point)
	}

	for i := 0; i <= 100; i++ {
		a := pos{x: 0, y: 0}.move(m.aButton.movement.multiply(i), false)
		if a.x > m.destination.x || a.y > m.destination.y {
			break
		}

		point := point{a, m.aButton.cost * i}

		if a.x == m.destination.x && a.y == m.destination.y {
			return point.cost
		}

		aLine = append(aLine, point)
	}

	for _, a := range aLine {
		for _, b := range bLine {
			if a.pos.same(b.pos) {
				return a.cost + b.cost
			}
		}
	}

	return 0
}

// ChatGippity because I can recognize Linear Algebra when I see it but I've forgotten how it works.
func playClawLA(m machine) int {
	// We want to solve:
	// A_x * a + B_x * b = X
	// A_y * a + B_y * b = Y
	//
	// Using Cramer's rule:
	// det = A_x*B_y - A_y*B_x
	// a = (X*B_y - Y*B_x) / det
	// b = (A_x*Y - A_y*X) / det

	A_x := m.aButton.movement.x
	A_y := m.aButton.movement.y
	B_x := m.bButton.movement.x
	B_y := m.bButton.movement.y
	X := m.destination.x
	Y := m.destination.y

	det := A_x*B_y - A_y*B_x
	if det == 0 {
		// No unique solution
		return 0
	}

	a_num := X*B_y - Y*B_x
	b_num := A_x*Y - A_y*X

	// Check divisibility
	if a_num%det != 0 || b_num%det != 0 {
		return 0
	}

	a := a_num / det
	b := b_num / det

	// Check non-negativity
	if a < 0 || b < 0 {
		return 0
	}

	return m.aButton.cost*a + m.bButton.cost*b
}

func startMachines(input []string) int {
	machines := parse(input, 3, 1, 0)

	sum := 0
	for _, m := range machines {
		sum += playClaw(m)
	}

	return sum
}

func startMachinesLA(input []string) int {
	machines := parse(input, 3, 1, 10000000000000)

	sum := 0
	for _, m := range machines {
		sum += playClawLA(m)
	}

	return sum
}
