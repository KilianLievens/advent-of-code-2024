package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/kilianlievens/advent-of-code-2024/advent"
)

func main() {
	exampleOneInput := advent.Read("./input/day_nine/example_one.txt")
	fmt.Printf("Example: defragment checksum: %d\n", defragment(exampleOneInput))
	fmt.Printf("Example: block defragment checksum: %d\n", blockDefragment(exampleOneInput))
	puzzleOneInput := advent.Read("./input/day_nine/puzzle_one.txt")
	fmt.Printf("Puzzle: defragment checksum: %d\n", defragment(puzzleOneInput))
	fmt.Printf("Puzzle: block defragment checksum: %d\n", blockDefragment(puzzleOneInput))
}

type space struct {
	id *uint64
}

func (b space) String() string {
	if b.id == nil {
		return "."
	}
	return fmt.Sprintf("%d", *b.id)
}

func parse(input []string) []space {
	advent.Assert(len(input) == 1, "Should only be one line")

	raw := strings.Split(input[0], "")

	disk := []space{}

	for i, v := range raw {
		size, _ := strconv.Atoi(v)

		for j := 0; j < size; j++ {
			if i%2 == 1 {
				disk = append(disk, space{id: nil})
				continue
			}
			id := uint64(i / 2)
			disk = append(disk, space{id: &id})
		}
	}

	return disk
}

func calcChecksum(disk []space) uint64 {
	var checksum uint64 = 0

	for i, v := range disk {
		if v.id == nil {
			continue
		}
		checksum += uint64(i) * *v.id
	}

	return checksum
}

func defragment(input []string) uint64 {
	disk := parse(input)

	leftPointer := 0
	rightPointer := len(disk) - 1

	for leftPointer < rightPointer {
		left := disk[leftPointer]

		if left.id == nil {
			right := disk[rightPointer]
			for right.id == nil {
				rightPointer--

				if leftPointer >= rightPointer {
					break
				}

				right = disk[rightPointer]
			}

			disk[leftPointer] = right
			disk[rightPointer] = left
		}

		leftPointer++
	}

	return calcChecksum(disk)
}

type blockType uint64

const (
	spaceType blockType = iota
	dataType
)

type block struct {
	size  uint64
	btype blockType
	id    *uint64
	moved bool
}

func (b block) String() string {
	if b.btype == spaceType {
		return strings.Repeat(".", int(b.size))
	}

	return strings.Repeat(fmt.Sprintf("%d", *b.id), int(b.size))
}

type disk []*block

func (d disk) toSpace() []space {
	s := []space{}

	for _, b := range d {
		for i := 0; i < int(b.size); i++ {
			if b.btype == spaceType {
				s = append(s, space{id: nil})
				continue
			}
			s = append(s, space{id: b.id})
		}
	}

	return s
}

func blockParse(input []string) disk {
	advent.Assert(len(input) == 1, "Should only be one line")

	raw := strings.Split(input[0], "")

	disk := []block{}

	for i, v := range raw {
		size, _ := strconv.Atoi(v)

		if i%2 == 1 {
			disk = append(disk, block{size: uint64(size), btype: spaceType})
			continue
		}
		id := uint64(i / 2)
		disk = append(disk, block{size: uint64(size), btype: dataType, id: &id})
	}

	merged := []*block{}
	for _, blk := range disk {
		if len(merged) > 0 && merged[len(merged)-1].btype == spaceType && blk.btype == spaceType {
			merged[len(merged)-1].size += blk.size
			continue
		}
		merged = append(merged, &blk)
	}

	return merged
}

func blockDefragment(input []string) uint64 {
	disk := blockParse(input)

	for r := len(disk) - 1; r >= 0; r-- {
		right := disk[r]

		if right.btype == spaceType {
			continue
		}

		if right.moved {
			continue
		}

		right.moved = true

		for l := 0; l < r; l++ {
			left := disk[l]

			if left.btype == dataType {
				continue
			}

			if left.size < right.size {
				continue
			}

			rest := left.size - right.size

			if rest == 0 {
				disk[l] = right
				disk[r] = left
				break
			}

			disk[l] = right
			disk[r] = &block{size: right.size, btype: spaceType}
			disk = slices.Insert(disk, l+1, &block{size: rest, btype: spaceType})
			break
		}
	}

	return calcChecksum(disk.toSpace())
}
