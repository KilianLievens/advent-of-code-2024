package main

import (
	"fmt"
	"slices"

	"github.com/kilianlievens/advent-of-code-2024/advent"
)

func main() {
	exampleOneInput := advent.Read("./input/day_six/example_one.txt")
	count, loops := followGuard(exampleOneInput)
	fmt.Printf("Example: pos count: %d\n", count)
	fmt.Printf("Example: loops count: %d\n", loops)
	puzzleOneInput := advent.Read("./input/day_six/puzzle_one.txt")
	count, loops = followGuard(puzzleOneInput)
	fmt.Printf("Puzzle: pos count: %d\n", count)
	fmt.Printf("Puzzle: loops count: %d\n", loops)
}

type tile int

const (
	open tile = iota
	blocked
)

var tileName = map[tile]string{
	open:    ".",
	blocked: "#",
}

func (t tile) String() string {
	return tileName[t]
}

type plane [][]tile

func (p *plane) Print() {
	for _, row := range *p {
		for _, r := range row {
			fmt.Printf("%s", r)
		}
		fmt.Printf("\n")
	}
}

func (p *plane) PrintBeen(been map[pos]bool) {
	for y, row := range *p {
		for x, r := range row {
			if been[pos{x: x, y: y}] {
				fmt.Printf("X")
				continue
			}
			fmt.Printf("%s", r)
		}
		fmt.Printf("\n")
	}
}

func (p *plane) Clone() plane {
	clone := make([][]tile, len(*p))
	for i, row := range *p {
		clone[i] = slices.Clone(row)
	}
	return clone
}

type direction int

const (
	up direction = iota
	right
	down
	left
)

var directionName = map[direction]string{
	up:    "^",
	right: ">",
	down:  "v",
	left:  "<",
}

func (d direction) String() string {
	return directionName[d]
}

type guard struct {
	pos       *pos
	direction direction
}

func (g *guard) Clone() guard {
	return guard{
		pos:       &pos{x: g.pos.x, y: g.pos.y},
		direction: g.direction,
	}
}

func (g *guard) Move(plane plane) bool {
	switch g.direction {
	case up:
		if g.pos.y == 0 {
			return true
		}

		proposedY := g.pos.y - 1
		if plane[proposedY][g.pos.x] == blocked {
			g.direction = right
			return g.Move(plane)
		}

		g.pos.y--
		return false

	case right:
		if g.pos.x == len(plane[0])-1 {
			return true
		}

		proposedX := g.pos.x + 1
		if plane[g.pos.y][proposedX] == blocked {
			g.direction = down
			return g.Move(plane)
		}

		g.pos.x++
		return false

	case down:
		if g.pos.y == len(plane)-1 {
			return true
		}

		proposedY := g.pos.y + 1
		if plane[proposedY][g.pos.x] == blocked {
			g.direction = left
			return g.Move(plane)
		}

		g.pos.y++
		return false

	case left:
		if g.pos.x == 0 {
			return true
		}

		proposedX := g.pos.x - 1
		if plane[g.pos.y][proposedX] == blocked {
			g.direction = up
			return g.Move(plane)
		}

		g.pos.x--
		return false
	}
	panic("unreachable")
}

func (g *guard) BeenState() beenState {
	return beenState{pos: *g.pos, direction: g.direction}
}

type pos struct {
	x, y int
}

func parse(input []string) (plane, *guard) {
	plane := plane{}
	var g *guard
	for y, inputLine := range input {
		parsedLine := []tile{}
		for x, char := range inputLine {
			if char == '.' {
				parsedLine = append(parsedLine, open)
				continue
			}

			if char == '#' {
				parsedLine = append(parsedLine, blocked)
				continue
			}

			if char == '^' {
				parsedLine = append(parsedLine, open)
				g = &guard{pos: &pos{x: x, y: y}, direction: up}
				continue
			}

			panic("Invalid character")
		}
		plane = append(plane, parsedLine)
	}

	advent.Assert(len(plane) == len(input), "Invalid input: too few Y")
	advent.Assert(len(plane[0]) == len(input[0]), "Invalid input: too few X")
	advent.Assert(g != nil, "Guard not found")

	return plane, g
}

type beenState struct {
	pos       pos
	direction direction
}

func followGuard(input []string) (int, int) {
	plane, og := parse(input)

	guard := og.Clone()
	beenPositions := map[pos]bool{}

	done := false
	for !done {
		beenPositions[*guard.pos] = true
		done = guard.Move(plane)
	}

	loops := 0

	for pos := range beenPositions {
		clone := plane.Clone()
		clone[pos.y][pos.x] = blocked
		beenStates := map[beenState]bool{}

		guard := og.Clone()
		done := false

		for !done {
			beenState := guard.BeenState()
			if beenStates[beenState] {
				loops++
				break
			}
			beenStates[beenState] = true
			done = guard.Move(clone)
		}
	}

	return len(beenPositions), loops
}
