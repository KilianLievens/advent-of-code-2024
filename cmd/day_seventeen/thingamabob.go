package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/kilianlievens/advent-of-code-2024/advent"
)

type registers struct {
	a, b, c int
}

func (rm registers) String() string {
	return fmt.Sprintf("A: %d, B: %d, C: %d", rm.a, rm.b, rm.c)
}

func (rm *registers) parse(r rune, v int) {
	switch r {
	case 'A':
		rm.a = v
	case 'B':
		rm.b = v
	case 'C':
		rm.c = v
	default:
		panic(fmt.Sprintf("Invalid register: %c", r))
	}
}

func (rm registers) comboOperand(operand int) int {
	switch operand {
	case 0, 1, 2, 3:
		return operand
	case 4:
		return int(rm.a)
	case 5:
		return int(rm.b)
	case 6:
		return int(rm.c)
	}

	panic(fmt.Sprintf("Invalid operand: %d", operand))
}

type output []string

func (o output) String() string {
	return strings.Join(o, ",")
}

type program []int

func (p program) execute(rm *registers) output {
	output := output{}
	ip := 0
	for ip < len(p) {
		instruction := instructions[p[ip]]
		advent.Assert(instruction != nil, "Invalid instruction")
		var iOut *string
		ip, iOut = instruction(rm, p[ip+1], ip)
		if iOut != nil {
			output = append(output, *iOut)
		}
	}

	return output
}

func (p program) String() string {
	strs := []string{}
	for _, i := range p {
		strs = append(strs, fmt.Sprintf("%d", i))
	}
	return strings.Join(strs, ",")
}

type instruction func(rm *registers, operand int, instructionPointer int) (nextIP int, output *string)
type instructionSet map[int]instruction

var instructions = instructionSet{
	0: adv,
	1: bxl,
	2: bst,
	3: jnz,
	4: bxc,
	5: out,
	6: bdv,
	7: cdv,
}

func adv(rm *registers, operand int, ip int) (int, *string) {
	numerator := rm.a
	demoninator := int(math.Pow(2, float64(rm.comboOperand(operand))))
	rm.a = int(numerator / demoninator)

	return ip + 2, nil
}

func bxl(rm *registers, operand int, ip int) (int, *string) {
	rm.b = rm.b ^ operand

	return ip + 2, nil
}

func bst(rm *registers, operand int, ip int) (int, *string) {
	rm.b = rm.comboOperand(operand) % 8

	return ip + 2, nil
}

func jnz(rm *registers, operand int, ip int) (int, *string) {
	if rm.a == 0 {
		return ip + 2, nil
	}

	return operand, nil
}

func bxc(rm *registers, _ int, ip int) (int, *string) {
	rm.b = rm.b ^ rm.c

	return ip + 2, nil
}

func out(rm *registers, operand int, ip int) (int, *string) {
	output := fmt.Sprintf("%d", rm.comboOperand(operand)%8)

	return ip + 2, &output
}

func bdv(rm *registers, operand int, ip int) (int, *string) {
	numerator := rm.a
	demoninator := int(math.Pow(2, float64(rm.comboOperand(operand))))
	rm.b = int(numerator / demoninator)

	return ip + 2, nil
}

func cdv(rm *registers, operand int, ip int) (int, *string) {
	numerator := rm.a
	demoninator := int(math.Pow(2, float64(rm.comboOperand(operand))))
	rm.c = int(numerator / demoninator)

	return ip + 2, nil
}
