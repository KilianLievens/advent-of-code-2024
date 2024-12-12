package main

import (
	"fmt"
	"strings"

	"github.com/kilianlievens/advent-of-code-2024/advent"
)

func main() {
	exampleOneInput := advent.Read("./input/day_twelve/example_one.txt")
	price, discountPrice := getPrice(exampleOneInput)
	fmt.Printf("Example 1: price: %d\n", price)
	fmt.Printf("Example 1: discount price: %d\n", discountPrice)
	exampleTwoInput := advent.Read("./input/day_twelve/example_two.txt")
	price, _ = getPrice(exampleTwoInput)
	fmt.Printf("Example 2: price: %d\n", price)
	exampleThreeInput := advent.Read("./input/day_twelve/example_three.txt")
	price, _ = getPrice(exampleThreeInput)
	fmt.Printf("Example 3: price: %d\n", price)
	exampleFourInput := advent.Read("./input/day_twelve/example_four.txt")
	_, discountPrice = getPrice(exampleFourInput)
	fmt.Printf("Example 4: discount price: %d\n", discountPrice)
	exampleFiveInput := advent.Read("./input/day_twelve/example_five.txt")
	_, discountPrice = getPrice(exampleFiveInput)
	fmt.Printf("Example 5: discount price: %d\n", discountPrice)
	puzzleOneInput := advent.Read("./input/day_twelve/puzzle_one.txt")
	price, discountPrice = getPrice(puzzleOneInput)
	fmt.Printf("Puzzle: price: %d\n", price)
	fmt.Printf("Puzzle: discount price: %d\n", discountPrice)
}

func parse(input []string) [][]string {
	matrix := [][]string{}

	for _, line := range input {
		splitLine := strings.Split(line, "")
		matrix = append(matrix, splitLine)
	}

	advent.Assert(len(matrix) == len(input), "Matrix height does not match input height")
	advent.Assert(len(matrix[0]) == len(input[0]), "Matrix width does not match input width")

	return matrix
}

type matrix [][]string

func (m matrix) get(pos pos) string {
	return m[pos.y][pos.x]
}

type pos struct {
	x, y uint
}

func (p pos) neighbours() []side {
	return []side{
		{pos: pos{p.x, p.y - 1}, alignment: vertical, dir: up},
		{pos: pos{p.x, p.y + 1}, alignment: vertical, dir: down},
		{pos: pos{p.x - 1, p.y}, alignment: horizontal, dir: left},
		{pos: pos{p.x + 1, p.y}, alignment: horizontal, dir: right},
	}
}

func (p pos) inBounds(obY, obX uint) bool {
	return p.y >= 0 && p.y < obY && p.x >= 0 && p.x < obX
}

type alignmentType int

const (
	horizontal alignmentType = iota
	vertical
)

type directionType int

const (
	up directionType = iota
	down
	left
	right
)

type side struct {
	pos       pos
	alignment alignmentType
	dir       directionType
}

type region struct {
	plots     *[]pos
	perimeter *[]side
	plant     string
}

func (r region) size() uint {
	return uint(len(*r.plots))
}

func floodFillRegion(matrix matrix, pos pos, beenPos map[pos]bool, r *region) {
	advent.Assert(r.plots != nil, "Region plots should be initialized")
	advent.Assert(r.perimeter != nil, "Perimeter should be initialized")
	advent.Assert(r.plant != "", "Region plant should be initialized")

	beenPos[pos] = true
	*r.plots = append(*r.plots, pos)

	obY := uint(len(matrix))
	obX := uint(len(matrix[0]))

	neighs := pos.neighbours()
	for _, n := range neighs {
		if !n.pos.inBounds(obY, obX) {
			*r.perimeter = append(*r.perimeter, n)
			continue
		}

		if matrix.get(n.pos) != r.plant {
			*r.perimeter = append(*r.perimeter, n)
			continue
		}

		if !beenPos[n.pos] {
			floodFillRegion(matrix, n.pos, beenPos, r)
		}
	}
}

type edge struct {
	val uint
	dir directionType
}

func combineHorizontalSides(sides *[]side) uint {
	ys := map[uint][]side{}
	for _, s := range *sides {
		if s.alignment != horizontal {
			ys[s.pos.y] = append(ys[s.pos.y], s)
		}
	}

	sideCount := uint(0)
	for _, xs := range ys {
		xsMap := map[edge]bool{}

		for _, x := range xs {
			xsMap[edge{x.pos.x, x.dir}] = true
		}

		for x := range xsMap {
			_, ok := xsMap[edge{x.val + 1, x.dir}]
			if !ok {
				sideCount++
			}
		}
	}

	return sideCount
}

func combineVerticalSides(sides *[]side) uint {
	xs := map[uint][]side{}
	for _, s := range *sides {
		if s.alignment != vertical {
			xs[s.pos.x] = append(xs[s.pos.x], s)
		}
	}

	sideCount := uint(0)
	for _, ys := range xs {
		ysMap := map[edge]bool{}

		for _, y := range ys {
			ysMap[edge{y.pos.y, y.dir}] = true
		}

		for y := range ysMap {
			_, ok := ysMap[edge{y.val + 1, y.dir}]
			if !ok {
				sideCount++
			}
		}
	}

	return sideCount
}

func getPrice(input []string) (uint, uint) {
	matrix := parse(input)

	beenPos := map[pos]bool{}

	regions := []*region{}
	for y, row := range matrix {
		for x, plant := range row {
			if beenPos[pos{uint(x), uint(y)}] {
				continue
			}

			r := &region{plant: plant, plots: &[]pos{}, perimeter: &[]side{}}
			floodFillRegion(matrix, pos{uint(x), uint(y)}, beenPos, r)
			regions = append(regions, r)
		}
	}

	price := uint(0)
	discountPrice := uint(0)
	for _, r := range regions {
		price += r.size() * uint(len(*r.perimeter))

		sides := combineHorizontalSides(r.perimeter) + combineVerticalSides(r.perimeter)

		discountPrice += r.size() * sides
	}

	return price, discountPrice
}
