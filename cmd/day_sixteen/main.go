package main

import (
	"container/heap"
	"fmt"

	"github.com/kilianlievens/advent-of-code-2024/advent"
)

func main() {
	exampleOneInput := advent.Read("./input/day_sixteen/example_one.txt")
	score, bestSeatCount := getScoreAndSeats(exampleOneInput)
	fmt.Printf("Example: Reindeer score: %d\n", score)
	fmt.Printf("Example: Best seat count: %d\n", bestSeatCount)
	puzzleOneInput := advent.Read("./input/day_sixteen/puzzle_one.txt")
	score, bestSeatCount = getScoreAndSeats(puzzleOneInput)
	fmt.Printf("Puzzle: Reindeer score: %d\n", score)
	fmt.Printf("Puzzle: Best seat count: %d\n", bestSeatCount)
}

func parse(input []string) (maze, pos, pos) {
	var deer, exit pos
	maze := maze{}

	for y, line := range input {
		row := []tileType{}
		for x, c := range line {
			if c == 'S' {
				deer = pos{x, y}
			}
			if c == 'E' {
				exit = pos{x, y}
			}
			row = append(row, newTile(c))
		}
		maze = append(maze, row)
	}

	advent.Assert(len(maze) == len(input) && len(maze[0]) == len(input[0]), "Invalid maze dimensions")

	return maze, deer, exit
}

type tileType int

func (t tileType) String() string {
	return tileName[t]
}

var tileName = map[tileType]string{
	wall:  "#",
	empty: ".",
}

func newTile(c rune) tileType {
	switch c {
	case '#':
		return wall
	case 'S', '.', 'E':
		return empty
	}
	panic("Unknown tile type")
}

const (
	wall tileType = iota
	empty
)

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

type movement struct {
	pos  pos
	dir  directionType
	cost int
}

func (p pos) moves(m maze, lastDirection directionType) []movement {
	potMoves := []movement{
		{pos: p.add(pos{y: -1}), dir: up, cost: up.getTurnCost(lastDirection) + 1},
		{pos: p.add(pos{y: 1}), dir: down, cost: down.getTurnCost(lastDirection) + 1},
		{pos: p.add(pos{x: -1}), dir: left, cost: left.getTurnCost(lastDirection) + 1},
		{pos: p.add(pos{x: 1}), dir: right, cost: right.getTurnCost(lastDirection) + 1},
	}

	prunedMoves := []movement{}
	for _, move := range potMoves {
		if m.inBounds(move.pos) && m.get(move.pos) != wall {
			prunedMoves = append(prunedMoves, move)
		}
	}

	return prunedMoves
}

type maze [][]tileType

func (w maze) print(deer, exit pos) {
	for y, line := range w {
		for x, r := range line {
			if deer.same(pos{y: y, x: x}) {
				fmt.Printf("S")
				continue
			}
			if exit.same(pos{y: y, x: x}) {
				fmt.Printf("E")
				continue
			}
			fmt.Printf("%s", r)
		}
		fmt.Printf("\n")
	}
}

func (w maze) printSeats(seats map[pos]bool) {
	for y, line := range w {
		for x, r := range line {
			if seats[pos{y: y, x: x}] {
				fmt.Printf("O")
				continue
			}

			fmt.Printf("%s", r)
		}
		fmt.Printf("\n")
	}
}

func (w maze) get(pos pos) tileType {
	return w[pos.y][pos.x]
}

func (w maze) inBounds(pos pos) bool {
	return pos.y >= 0 && pos.y < len(w) && pos.x >= 0 && pos.x < len(w[0])
}

type directionType int

const (
	up directionType = iota
	down
	left
	right
)

var directionName = map[directionType]string{
	up:    "up",
	down:  "down",
	left:  "left",
	right: "right",
}

func (d directionType) String() string {
	return directionName[d]
}

func (d directionType) inverse() directionType {
	switch d {
	case up:
		return down
	case down:
		return up
	case left:
		return right
	case right:
		return left
	default:
		panic("Unknown direction")
	}
}

func (d directionType) getTurnCost(lastDirection directionType) int {
	switch d {
	case lastDirection:
		return 0
	case d.inverse():
		return 2000
	default:
		return 1000
	}
}

type nodeKey struct {
	lastDirection directionType
	pos           pos
}

func (k nodeKey) String() string {
	return fmt.Sprintf("(Dir: %s, Pos: %s)", k.lastDirection, k.pos)
}

type node struct {
	pos           pos
	cost          int
	estimatedCost int
	lastDirection directionType
	path          []pos
}

type nodePrioQueue []node

func (q nodePrioQueue) Len() int {
	return len(q)
}

func (q nodePrioQueue) Less(i, j int) bool {
	return q[i].estimatedCost < q[j].estimatedCost
}

func (q nodePrioQueue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

func (q *nodePrioQueue) Push(x any) {
	*q = append(*q, x.(node))
}

func (q *nodePrioQueue) Pop() any {
	old := *q
	n := len(old)
	x := old[n-1]
	*q = old[:n-1]
	return x
}

func calcManhattanDistance(a, b pos) int {
	return advent.AbsInt(a.x-b.x) + advent.AbsInt(a.y-b.y)
}

func estimateTurnCost(curDir directionType, curPos, end pos) int {
	turnCost := 0
	if curPos.x < end.x {
		turnCost += right.getTurnCost(curDir)
	}
	if curPos.x > end.x {
		turnCost += left.getTurnCost(curDir)
	}
	if curPos.y < end.y {
		turnCost += down.getTurnCost(curDir)
	}
	if curPos.y > end.y {
		turnCost += up.getTurnCost(curDir)
	}

	return turnCost
}

func aStar(start, end pos, m maze) int {
	var prioQ = &nodePrioQueue{}
	var doneNodes = make(map[nodeKey]bool)

	firstNode := node{
		path:          []pos{start},
		pos:           start,
		lastDirection: right,
	}

	heap.Push(prioQ, firstNode)

	for prioQ.Len() > 0 {
		curNode := heap.Pop(prioQ).(node)

		if curNode.pos.same(end) {
			return curNode.cost
		}

		key := nodeKey{
			lastDirection: curNode.lastDirection,
			pos:           curNode.pos,
		}

		if _, done := doneNodes[key]; done {
			continue
		}

		doneNodes[key] = true

		moves := curNode.pos.moves(m, curNode.lastDirection)
		for _, move := range moves {
			cost := curNode.cost + move.cost

			path := make([]pos, len(curNode.path)+1)
			copy(path, curNode.path)

			path[len(curNode.path)] = move.pos

			nextNode := node{
				lastDirection: move.dir,
				path:          path,
				pos:           move.pos,
				cost:          cost,
				estimatedCost: cost + calcManhattanDistance(curNode.pos, move.pos) + estimateTurnCost(move.dir, move.pos, end),
			}
			nextNodeKey := nodeKey{
				lastDirection: nextNode.lastDirection,
				pos:           nextNode.pos,
			}

			if _, done := doneNodes[nextNodeKey]; !done {
				heap.Push(prioQ, nextNode)
			}
		}
	}

	panic("Could not find end.")
}

func aNotStarAtAll(start, end pos, m maze) int {
	var prioQ = &nodePrioQueue{}
	var doneNodes = make(map[nodeKey]int)
	bestCost := -1
	bestSpots := []pos{}
	altPaths := map[pos][]pos{}

	firstNode := node{
		path:          []pos{start},
		pos:           start,
		lastDirection: right,
	}

	heap.Push(prioQ, firstNode)

	for prioQ.Len() > 0 {
		curNode := heap.Pop(prioQ).(node)

		if bestCost >= 0 && curNode.cost > bestCost {
			break
		}

		key := nodeKey{
			lastDirection: curNode.lastDirection,
			pos:           curNode.pos,
		}

		if existingCost, ok := doneNodes[key]; ok {
			if curNode.cost <= existingCost {
				altPaths[curNode.pos] = append(altPaths[curNode.pos], curNode.path...)
			}
			continue
		}

		doneNodes[key] = curNode.cost

		if curNode.pos.same(end) {
			bestSpots = append(bestSpots, curNode.path...)

			if bestCost < 0 {
				bestCost = curNode.cost
				continue
			}
		}

		moves := curNode.pos.moves(m, curNode.lastDirection)
		for _, move := range moves {
			cost := curNode.cost + move.cost

			path := make([]pos, len(curNode.path)+1)
			copy(path, curNode.path)

			path[len(curNode.path)] = move.pos

			nextNode := node{
				lastDirection: move.dir,
				path:          path,
				pos:           move.pos,
				cost:          cost,
				estimatedCost: cost + calcManhattanDistance(curNode.pos, move.pos) + estimateTurnCost(move.dir, move.pos, end),
			}

			heap.Push(prioQ, nextNode)
		}
	}

	qualitySeats := map[pos]bool{}
	processedAltPaths := map[pos]bool{}
	for _, spot := range bestSpots {
		qualitySeats[spot] = true
		processAltPaths(spot, altPaths, qualitySeats, processedAltPaths)
	}

	// m.printSeats(qualitySeats)
	return len(qualitySeats)
}

func processAltPaths(p pos, altPaths map[pos][]pos, qualitySeats map[pos]bool, processedAltPaths map[pos]bool) {
	if processedAltPaths[p] {
		return
	}

	alts, ok := altPaths[p]
	if !ok {
		return
	}

	processedAltPaths[p] = true

	for _, alt := range alts {
		qualitySeats[alt] = true
		processAltPaths(alt, altPaths, qualitySeats, processedAltPaths)
	}
}

func getScoreAndSeats(input []string) (int, int) {
	maze, deer, exit := parse(input)
	// maze.print(deer, exit)
	return aStar(deer, exit, maze), aNotStarAtAll(deer, exit, maze)
}
