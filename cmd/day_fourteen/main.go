package main

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/kilianlievens/advent-of-code-2024/advent"
)

func main() {
	exampleOneInput := advent.Read("./input/day_fourteen/example_one.txt")
	fmt.Printf("Example: safety score: %d\n", getSafetyScore(exampleOneInput, 11, 7))
	puzzleOneInput := advent.Read("./input/day_fourteen/puzzle_one.txt")
	fmt.Printf("Puzzle: safety score: %d\n", getSafetyScore(puzzleOneInput, 101, 103))
	detectChristmas(puzzleOneInput, 101, 103)
	printTree(puzzleOneInput, 101, 103)
}

type pos struct {
	x, y int
}

type robot struct {
	pos    pos
	vector pos
	loop   int
}

func (r robot) String() string {
	return fmt.Sprintf("Pos: %v, Vector: %v, Loop: %d", r.pos, r.vector, r.loop)
}

func (r *robot) move(lenX, lenY int) {
	r.pos = pos{teleport(r.pos.x+r.vector.x, lenX), teleport(r.pos.y+r.vector.y, lenY)}
}

func teleport(z, length int) int {
	if z < 0 {
		return length + z
	}

	return z % length
}

func parse(input []string) []*robot {
	robots := []*robot{}

	re := regexp.MustCompile(`p=(\d+),(\d+) v=(-?\d+),(-?\d+)`)

	for _, line := range input {
		raw := re.FindStringSubmatch(line)

		advent.Assert(len(raw) == 5, "Raw not expected format")

		x, err1 := strconv.Atoi(raw[1])
		y, err2 := strconv.Atoi(raw[2])
		vx, err3 := strconv.Atoi(raw[3])
		vy, err4 := strconv.Atoi(raw[4])

		advent.Assert(err1 == nil && err2 == nil && err3 == nil && err4 == nil, "Error converting")

		robots = append(robots, &robot{pos{x, y}, pos{vx, vy}, 0})
	}

	return robots
}

func show(robots map[pos][]*robot, lenX, lenY int) {
	for y := 0; y < lenY; y++ {
		for x := 0; x < lenX; x++ {
			robots, ok := robots[pos{x, y}]
			if !ok {
				fmt.Print(".")
				continue
			}
			if len(robots) > 9 {
				fmt.Print("X")
				continue
			}
			fmt.Print(len(robots))
		}
		fmt.Println()
	}

}

type quadrant struct {
	start pos
	end   pos
}

func getSafetyScore(input []string, lenX, lenY int) int {
	robots := parse(input)

	for i := 0; i < 100; i++ {
		for _, r := range robots {
			r.move(lenX, lenY)
		}
	}

	robotMap := map[pos][]*robot{}
	for _, r := range robots {
		robotMap[r.pos] = append(robotMap[r.pos], r)
	}

	middleX := (lenX - 1) / 2
	middleY := (lenY - 1) / 2
	endX := lenX - 1
	endY := lenY - 1

	quadrants := []quadrant{
		{pos{0, 0}, pos{middleX - 1, middleY - 1}},
		{pos{middleX + 1, 0}, pos{endX, middleY - 1}},
		{pos{0, middleY + 1}, pos{middleX - 1, endY}},
		{pos{middleX + 1, middleY + 1}, pos{endX, endY}},
	}

	score := 1
	for _, q := range quadrants {
		count := 0
		for y := q.start.y; y <= q.end.y; y++ {
			for x := q.start.x; x <= q.end.x; x++ {
				if robots, ok := robotMap[pos{x, y}]; ok {
					count += len(robots)
				}
			}
		}
		score *= count
	}

	return score
}

func printTree(input []string, lenX, lenY int) {
	robots := parse(input)

	for i := 0; i < 10403; i++ {
		for _, r := range robots {
			r.move(lenX, lenY)
		}
		robotMap := map[pos][]*robot{}
		for _, r := range robots {
			robotMap[r.pos] = append(robotMap[r.pos], r)
		}

		middleX := (lenX - 1) / 2

		for y := 0; y < lenY; y++ {
			if robotMap[pos{middleX, y}] != nil && robotMap[pos{middleX - 1, y}] != nil && robotMap[pos{middleX + 1, y}] != nil && robotMap[pos{middleX - 2, y}] != nil && robotMap[pos{middleX + 2, y}] != nil {
				show(robotMap, lenX, lenY)
				fmt.Println(i + 1)
				fmt.Println()
				fmt.Println()
				fmt.Println()
			}
		}
	}
}

func detectChristmas(input []string, lenX, lenY int) {
	robots := parse(input)

	loops := []int{}

	for _, r := range robots {
		ogPos := r.pos
		for {
			r.move(lenX, lenY)
			r.loop++
			if r.pos == ogPos {
				loops = append(loops, r.loop)
				break
			}
		}
	}

	advent.Assert(len(loops) > 2, "Expected at least two loops")
	lcm := advent.LCM(loops[0], loops[1], loops[2:]...)
	fmt.Println("LCM:", lcm)
}
