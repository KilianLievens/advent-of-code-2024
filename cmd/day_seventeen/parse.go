package main

import (
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/kilianlievens/advent-of-code-2024/advent"
)

func read(fileName string) []string {
	body, err := os.ReadFile(fileName)
	advent.Assert(err == nil, "unable to read file: %v")

	var lines []string
	for _, line := range strings.Split(string(body), "\n") {
		lines = append(lines, line)
	}

	return lines
}

func parse(input []string) (registers, program) {
	readingRegisters := true
	registers := registers{}
	program := make(program, 0)

	regMatch := regexp.MustCompile(`Register (.{1}): (\d+)`)
	progMatch := regexp.MustCompile(`Program: (.+)`)

	for _, line := range input {
		if line == "" {
			readingRegisters = false
			continue
		}

		if readingRegisters {
			matches := regMatch.FindStringSubmatch(line)
			advent.Assert(len(matches) == 3, "Invalid register line")
			advent.Assert(len(matches[1]) == 1, "Invalid register name")
			reg := rune(matches[1][0])
			val, err := strconv.Atoi(matches[2])
			advent.Assert(err == nil, "Invalid register value")
			registers.parse(reg, val)
			continue
		}

		matches := progMatch.FindStringSubmatch(line)
		advent.Assert(len(matches) == 2, "Invalid program line")
		for _, instr := range strings.Split(matches[1], ",") {
			val, err := strconv.Atoi(instr)
			advent.Assert(err == nil, "Invalid program instruction")
			program = append(program, val)
		}
	}

	return registers, program
}
