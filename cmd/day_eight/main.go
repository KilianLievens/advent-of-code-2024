package main

import (
	"fmt"

	"github.com/kilianlievens/advent-of-code-2024/advent"
)

func main() {
	exampleOneInput := advent.Read("./input/day_eight/example_one.txt")
	antiNodeCount, resonantAntiNodeCount := countUniqueAntiNode(exampleOneInput)
	fmt.Printf("Example: unique anti-node count: %d\n", antiNodeCount)
	fmt.Printf("Example: unique resonant anti-node count: %d\n", resonantAntiNodeCount)
	puzzleOneInput := advent.Read("./input/day_eight/puzzle_one.txt")
	antiNodeCount, resonantAntiNodeCount = countUniqueAntiNode(puzzleOneInput)
	fmt.Printf("Puzzle: unique anti-node count: %d\n", antiNodeCount)
	fmt.Printf("Puzzle: unique resonant anti-node count: %d\n", resonantAntiNodeCount)
}

type pos struct {
	x, y int
}

func (p *pos) diff(o pos) (int, int) {
	return p.y - o.y, p.x - o.x
}

func (p *pos) add(y, x int) pos {
	return pos{x: p.x + x, y: p.y + y}
}

func (p pos) inBounds(yMax, xMax int) bool {
	return p.x >= 0 && p.x < xMax && p.y >= 0 && p.y < yMax
}

func (p pos) String() string {
	return fmt.Sprintf("(y:%d, x:%d)", p.y, p.x)
}

type node struct {
	isAntenna bool
	value     string
}

func (n node) String() string {
	if n.isAntenna {
		return n.value
	}
	return "."
}

type plane [][]node

func (p *plane) Print() {
	for _, row := range *p {
		for _, n := range row {
			fmt.Print(n)
		}
		fmt.Println()
	}

}

func parse(input []string) (plane, map[string][]pos) {
	p := plane{}
	antennaGroup := map[string][]pos{}

	for y, row := range input {
		line := []node{}
		for x, c := range row {
			isAntenna := c != '.'
			if isAntenna {
				group := antennaGroup[string(c)]
				group = append(group, pos{x: x, y: y})
				antennaGroup[string(c)] = group
			}
			line = append(line, node{value: string(c), isAntenna: isAntenna})
		}
		p = append(p, line)
	}

	advent.Assert(len(p) == len(input), "parsed plane has wrong number of rows")
	advent.Assert(len(p[0]) == len(input[0]), "parsed plane has wrong number of columns")

	return p, antennaGroup
}

func getAntiNodes(p [2]pos, yMax, xMax int) []pos {
	yDiff, xDiff := p[0].diff(p[1])

	antiNodes := [2]pos{
		p[0].add(yDiff, xDiff),
		p[1].add(-yDiff, -xDiff),
	}

	filteredAntiNodes := []pos{}
	for _, antiNode := range antiNodes {
		if antiNode.inBounds(yMax, xMax) {
			filteredAntiNodes = append(filteredAntiNodes, antiNode)
		}
	}

	return filteredAntiNodes
}

func getResonantAntiNodes(p [2]pos, yMax, xMax int) []pos {
	yDiff, xDiff := p[0].diff(p[1])

	antiNodes := []pos{}

	i := 0
	for {
		breaker := true
		a := p[0].add(yDiff*i, xDiff*i)
		b := p[1].add(-yDiff*i, -xDiff*i)

		if a.inBounds(yMax, xMax) {
			antiNodes = append(antiNodes, a)
			breaker = false
		}

		if b.inBounds(yMax, xMax) {
			antiNodes = append(antiNodes, b)
			breaker = false
		}

		if breaker {
			break
		}

		i++
	}

	return antiNodes
}

func countUniqueAntiNode(input []string) (int, int) {
	_, aG := parse(input)

	uniqueAntiNodes := map[pos]bool{}
	uniqueResonantAntiNodes := map[pos]bool{}

	for _, group := range aG {
		for i := 0; i < len(group)-1; i++ {
			for j := i + 1; j < len(group); j++ {
				antiNodes := getAntiNodes([2]pos{group[i], group[j]}, len(input), len(input[0]))
				for _, antiNode := range antiNodes {
					uniqueAntiNodes[antiNode] = true
				}

				resonantAntiNodes := getResonantAntiNodes([2]pos{group[i], group[j]}, len(input), len(input[0]))
				for _, resAntiNode := range resonantAntiNodes {
					uniqueResonantAntiNodes[resAntiNode] = true
				}
			}
		}
	}

	return len(uniqueAntiNodes), len(uniqueResonantAntiNodes)
}
