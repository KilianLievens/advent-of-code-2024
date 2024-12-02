package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kilianlievens/advent-of-code-2024/advent"
)

func main() {
	exampleOneInput := advent.Read("./input/day_two/example_one.txt")
	fmt.Printf("Example: safety: %d\n", getSafetyNumber(exampleOneInput, false))         // 2
	fmt.Printf("Example: dampened safety: %d\n", getSafetyNumber(exampleOneInput, true)) // 4
	puzzleOneInput := advent.Read("./input/day_two/puzzle_one.txt")
	fmt.Printf("Puzzle: safety: %d\n", getSafetyNumber(puzzleOneInput, false))         // 306
	fmt.Printf("Puzzle: dampened safety: %d\n", getSafetyNumber(puzzleOneInput, true)) // 366
}

func parse(input []string) [][]int {
	reports := [][]int{}
	for _, line := range input {
		splitLine := strings.Split(line, " ")

		levels := []int{}
		for _, rawNum := range splitLine {
			num, _ := strconv.Atoi(rawNum)
			levels = append(levels, num)
		}

		reports = append(reports, levels)
	}

	return reports
}

func getDir(diff int) int {
	advent.Assert(diff != 0, "diff should not be 0")

	if diff > 0 {
		return 1
	}

	return -1
}

func pruneReport(report []int, toRemove int) []int {
	res := make([]int, 0)
	res = append(res, report[:toRemove]...)
	res = append(res, report[toRemove+1:]...)

	advent.Assert(len(res) == len(report)-1, "mockNewReport: length of new report should be one less than the original report")

	return res
}

func getReportSafety(report []int) int {
	advent.Assert(len(report) >= 2, "report should be at least two values")

	reportDir := 0
	for i, level := range report {
		if i == 0 {
			continue
		}

		diff := level - report[i-1]

		if advent.AbsInt(diff) < 1 || advent.AbsInt(diff) > 3 {
			return 0
		}

		dir := getDir(diff)

		if reportDir == 0 {
			reportDir = dir
		}

		if dir != reportDir {
			return 0
		}
	}

	return 1
}

func getSafetyNumber(input []string, dampen bool) int {
	reports := parse(input)

	safety := 0
	for _, report := range reports {
		reportSafety := 0

		reportSafety += getReportSafety(report)

		if dampen {
			for i := 0; i < len(report); i++ {
				newReport := pruneReport(report, i)
				reportSafety += getReportSafety(newReport)
			}
		}

		if reportSafety > 0 {
			safety++
		}
	}

	return safety
}
