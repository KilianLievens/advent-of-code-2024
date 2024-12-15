package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/kilianlievens/advent-of-code-2024/advent"
)

func main() {
	exampleTwoInput := read("./input/day_fifteen/example_two.txt")
	fmt.Printf("Example 2: Boxes : %d\n", findBoxes(exampleTwoInput))
	exampleOneInput := read("./input/day_fifteen/example_one.txt")
	fmt.Printf("Example 1: Boxes : %d\n", findBoxes(exampleOneInput))
	fmt.Printf("Example 1: Alt Boxes : %d\n", findAltBoxes(exampleOneInput))
	puzzleOneInput := read("./input/day_fifteen/puzzle_one.txt")
	fmt.Printf("Puzzle: Boxes : %d\n", findBoxes(puzzleOneInput))
	fmt.Printf("Puzzle: Alt Boxes : %d\n", findAltBoxes(puzzleOneInput))
}

func read(fileName string) []string {
	body, err := os.ReadFile(fileName)
	advent.Assert(err == nil, "unable to read file: %v")

	var lines []string
	for _, line := range strings.Split(string(body), "\n") {
		lines = append(lines, line)
	}

	return lines
}

type tileType int

func (t tileType) String() string {
	return tileName[t]
}

const (
	wall tileType = iota
	box
	empty
	leftBox
	rightBox
)

var tileName = map[tileType]string{
	wall:     "#",
	box:      "O",
	empty:    ".",
	leftBox:  "[",
	rightBox: "]",
}

func newTile(c rune) tileType {
	switch c {
	case '#':
		return wall
	case 'O':
		return box
	case '@':
		return empty
	case '.':
		return empty
	}
	panic("Unknown tile type")
}

func newAltTile(c rune) []tileType {
	switch c {
	case '#':
		return []tileType{wall, wall}
	case 'O':
		return []tileType{leftBox, rightBox}
	case '@':
		return []tileType{empty, empty}
	case '.':
		return []tileType{empty, empty}
	}
	panic("Unknown tile type")
}

type warehouse [][]tileType

func (w warehouse) print(robotPos pos) {
	for y, line := range w {
		for x, r := range line {
			if robotPos.same(pos{y: y, x: x}) {
				fmt.Printf("@")
				continue
			}
			fmt.Printf("%s", r)
		}
		fmt.Printf("\n")
	}
}

func (w warehouse) get(pos pos) *tileType {
	return &w[pos.y][pos.x]
}

func (w warehouse) inBounds(pos pos) bool {
	return pos.y >= 0 && pos.y < len(w) && pos.x >= 0 && pos.x < len(w[0])
}

func (w warehouse) checkMove(pos pos, dir pos) bool {
	nextPos := pos.add(dir)

	if !w.inBounds(nextPos) {
		advent.Assert(false, "Should not be possible to go OB")
	}

	tile := *w.get(nextPos)

	switch tile {
	case wall:
		return false
	case box, leftBox, rightBox:
		return w.checkMove(nextPos, dir)
	case empty:
		return true
	default:
		panic("Unknown tile type")
	}
}

func (w warehouse) checkMoveDouble(p pos, dir pos) bool {
	nextPos := p.add(dir)

	if !w.inBounds(nextPos) {
		advent.Assert(false, "Should not be possible to go OB")
	}

	tile := *w.get(nextPos)

	switch tile {
	case wall:
		return false
	case empty:
		return true
	case box:
		panic("Should not be possible to move into a box in double move")
	case leftBox:
		return w.checkMoveDouble(nextPos, dir) && w.checkMoveDouble(nextPos.add(pos{x: 1}), dir)
	case rightBox:
		return w.checkMoveDouble(nextPos, dir) && w.checkMoveDouble(nextPos.add(pos{x: -1}), dir)
	default:
		panic("Unknown tile type")
	}
}

func (w warehouse) move(pos pos, dir pos) pos {
	nextPos := pos.add(dir)

	if !w.inBounds(nextPos) {
		advent.Assert(false, "Should not be possible to go OB")
	}

	nextTile := w.get(nextPos)

	advent.Assert(*nextTile != wall, "Should not be possible to move into a wall")

	if *nextTile == empty {
		return nextPos
	}

	if *nextTile == box || *nextTile == leftBox || *nextTile == rightBox {
		w.move(nextPos, dir)
	}

	nextNextTile := w.get(nextPos.add(dir))

	advent.Assert(*nextNextTile != wall, "Should not be possible to move a box into a wall")

	*nextNextTile = *nextTile
	*nextTile = empty

	return nextPos
}

func (w warehouse) moveDouble(p pos, dir pos) pos {
	nextPos := p.add(dir)

	if !w.inBounds(nextPos) {
		advent.Assert(false, "Should not be possible to go OB")
	}

	nextTile := w.get(nextPos)

	advent.Assert(*nextTile != wall, "Should not be possible to move into a wall")

	if *nextTile == empty {
		return nextPos
	}

	advent.Assert(*nextTile != box, "Should not be possible to move into a box in double move")

	var altNextTile *tileType
	var altNextNextTile *tileType
	switch *nextTile {
	case leftBox:
		altNextTile = w.get(nextPos.add(pos{x: 1}))
		altNextNextTile = w.get(nextPos.add(pos{x: 1}).add(dir))
		w.moveDouble(nextPos, dir)
		w.moveDouble(nextPos.add(pos{x: 1}), dir)
	case rightBox:
		altNextTile = w.get(nextPos.add(pos{x: -1}))
		altNextNextTile = w.get(nextPos.add(pos{x: -1}).add(dir))
		w.moveDouble(nextPos, dir)
		w.moveDouble(nextPos.add(pos{x: -1}), dir)
	default:
		panic("Box type not found")
	}

	nextNextTile := w.get(nextPos.add(dir))

	advent.Assert(*altNextTile != wall && *nextNextTile != wall && *altNextNextTile != wall, "Should not be possible to move a box into a wall")

	*nextNextTile = *nextTile
	*altNextNextTile = *altNextTile
	*nextTile = empty
	*altNextTile = empty

	return nextPos
}

type move int

func (m move) String() string {
	return moveName[m]
}

const (
	up move = iota
	down
	left
	right
)

var moveName = map[move]string{
	up:    "^",
	down:  "v",
	left:  "<",
	right: ">",
}

var moveDir = map[move]pos{
	up:    {y: -1},
	down:  {y: 1},
	left:  {x: -1},
	right: {x: 1},
}

func newMove(c rune) move {
	switch c {
	case '^':
		return up
	case 'v':
		return down
	case '<':
		return left
	case '>':
		return right
	}
	panic("Unknown move")
}

type pos struct {
	x, y int
}

func (p pos) String() string {
	return fmt.Sprintf("(x: %d,y: %d)", p.x, p.y)
}

func (p pos) add(dir pos) pos {
	return pos{x: p.x + dir.x, y: p.y + dir.y}
}

func (p pos) same(other pos) bool {
	return p.x == other.x && p.y == other.y
}

func parse(input []string) (warehouse, []move, pos) {
	parsingWarehouse := true
	warehouse := warehouse{}
	moves := []move{}
	robotPos := pos{}

	for y, line := range input {
		if line == "" {
			parsingWarehouse = false
			continue
		}

		if parsingWarehouse {
			row := []tileType{}
			for x, c := range line {
				row = append(row, newTile(c))
				if c == '@' {
					robotPos = pos{y: y, x: x}
				}
			}
			warehouse = append(warehouse, row)
			continue
		}

		for _, c := range line {
			moves = append(moves, newMove(c))
		}
	}

	return warehouse, moves, robotPos
}

func parseAlt(input []string) (warehouse, []move, pos) {
	parsingWarehouse := true
	warehouse := warehouse{}
	moves := []move{}
	robotPos := pos{}

	for y, line := range input {
		if line == "" {
			parsingWarehouse = false
			continue
		}

		if parsingWarehouse {
			row := []tileType{}
			for x, c := range line {
				row = append(row, newAltTile(c)...)
				if c == '@' {
					robotPos = pos{y: y, x: x * 2}
				}
			}
			warehouse = append(warehouse, row)
			continue
		}

		for _, c := range line {
			moves = append(moves, newMove(c))
		}
	}

	return warehouse, moves, robotPos
}

func findBoxes(input []string) int {
	warehouse, moves, robotPos := parse(input)

	// warehouse.print(robotPos)
	for _, move := range moves {
		if warehouse.checkMove(robotPos, moveDir[move]) {
			robotPos = warehouse.move(robotPos, moveDir[move])
		}
		// fmt.Println("move", move)
		// warehouse.print(robotPos)
	}

	sum := 0
	for y, row := range warehouse {
		for x, tile := range row {
			if tile == box {
				sum += 100*y + x
			}
		}
	}

	return sum
}

func findAltBoxes(input []string) int {
	warehouse, moves, robotPos := parseAlt(input)

	// warehouse.print(robotPos)
	for _, move := range moves {
		switch move {
		case left, right:
			if warehouse.checkMove(robotPos, moveDir[move]) {
				robotPos = warehouse.move(robotPos, moveDir[move])
			}
		case up, down:
			if warehouse.checkMoveDouble(robotPos, moveDir[move]) {
				robotPos = warehouse.moveDouble(robotPos, moveDir[move])
			}
		}
		// fmt.Println("move", move)
		// warehouse.print(robotPos)
	}

	sum := 0
	for y, row := range warehouse {
		for x, tile := range row {
			if tile == leftBox {
				sum += 100*y + x
			}
		}
	}

	return sum
}
