package main

import (
	"bufio"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/samber/lo"
)

type Operator int

const (
	OpAdd Operator = iota
	OpMul
	NumOfDiffOperators
)

func main() {
	file, err := os.Open("../input/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		closeErr := file.Close()
		if closeErr != nil {
			log.Fatal(closeErr)
		}
	}(file)

	pattern := `^([1-9][0-9]*):((\s+([1-9][0-9]*))+)\s*$`
	re, err := regexp.Compile(pattern)
	if err != nil {
		log.Panicf("Failed to compile regex: %v", err)
	}

	scanner := bufio.NewScanner(file)
	iLine := 0
	maxAttainable := int64(0)
	runningTotal := int64(0)
	for scanner.Scan() {
		line := scanner.Text()
		iLine++
		matches := re.FindStringSubmatch(line)
		if len(matches) == 0 {
			log.Printf("could not match line %d (`%v`)", iLine, line)
		}

		result, err := strconv.ParseInt(matches[1], 10, 64)
		if err != nil {
			log.Panicf("internal error: could not convert `%v` to int64", matches[1])
		}

		operandStrings := strings.Fields(matches[2])
		operands := lo.Map(operandStrings, func(s string, _ int) int64 {
			operand, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				log.Panicf("internal error: could not convert `%v` to int64", s)
			}
			return operand
		})

		maxAttainable += result
		if solvable(result, operands) {
			runningTotal += result
		}
	}

	log.Printf("max attainable: %d", maxAttainable)
	log.Printf("running total: %d", runningTotal)
}

func solvable(result int64, operands []int64) bool {
	nOperands := len(operands)
	nOps := nOperands - 1
	ops := make([]Operator, nOps)
	nCombinations := math.Pow(float64(NumOfDiffOperators), float64(nOps))
	for iCombo := range int64(nCombinations) {
		combo := iCombo
		for iOp := range nOps {
			ops[iOp] = Operator(combo % int64(NumOfDiffOperators))
			combo /= int64(NumOfDiffOperators)
		}

		if result == calc(operands, ops) {
			return true
		}
	}

	return false
}

func calc(operands []int64, ops []Operator) int64 {
	nOperands := len(operands)
	if nOperands < 1 {
		return 0
	}

	result := operands[0]
	for idx := 1; idx < nOperands; idx++ {
		switch ops[idx-1] {
		case OpAdd:
			result += operands[idx]
		case OpMul:
			result *= operands[idx]
		case NumOfDiffOperators:
			log.Panicf("internal error: unknown operator %d", ops[idx-1]) //nolint:revive // Toy code
		}
	}

	return result
}
