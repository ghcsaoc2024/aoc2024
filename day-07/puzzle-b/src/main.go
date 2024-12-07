package main

import (
	"bufio"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/hashicorp/go-set/v3"
	"github.com/samber/lo"
)

type Operator int

const (
	OpAdd Operator = iota
	OpMul
	OpConcat
	NumOfDiffOperators
)

func main() {
	var args struct {
		SweetSpot float64 `arg:"positional,required" help:"sweet spot for meet-in-the-\"middle\""`
	}
	arg.MustParse(&args)

	switch {
	case args.SweetSpot <= 0:
		fallthrough
	case args.SweetSpot >= 1:
		log.Fatalf("sweet spot must be larger than 0.0 and smaller than 1.0; got %f", args.SweetSpot)
	}

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
		if isSolvable(result, operands, args.SweetSpot) {
			runningTotal += result
		}
	}

	log.Printf("max attainable: %d", maxAttainable)
	log.Printf("running total: %d", runningTotal)
}

func isSolvable(desiredResult int64, operands []int64, sweetSpot float64) bool {
	nOperands := len(operands)
	nOps := nOperands - 1
	middleOpIdx := int(math.Round(float64(nOps) * sweetSpot))
	semiSolutions := set.New[int64](0)

	var ops []Operator
	var nCombinations int64

	// First half
	nFirstHalfOps := middleOpIdx
	ops = make([]Operator, nFirstHalfOps)
	nCombinations = int64(math.Pow(float64(NumOfDiffOperators), float64(nFirstHalfOps)))
	for iCombo := range nCombinations {
		combo := iCombo
		for iOp := range nFirstHalfOps {
			ops[iOp] = Operator(combo % int64(NumOfDiffOperators))
			combo /= int64(NumOfDiffOperators)
		}

		result := calcFwd(operands, ops, desiredResult)
		if result == -1 {
			continue
		}
		semiSolutions.Insert(result)
	}

	// Second half
	nSecondHalfOps := nOps - middleOpIdx
	ops = make([]Operator, nSecondHalfOps)
	nCombinations = int64(math.Pow(float64(NumOfDiffOperators), float64(nSecondHalfOps)))
	for iCombo := range nCombinations {
		combo := iCombo
		for iOp := range nSecondHalfOps {
			ops[iOp] = Operator(combo % int64(NumOfDiffOperators))
			combo /= int64(NumOfDiffOperators)
		}

		result := calcBack(operands, ops, desiredResult)
		if result == -1 {
			continue
		}

		if semiSolutions.Contains(result) {
			return true
		}
	}

	return false
}

func calcFwd(operands []int64, ops []Operator, desiredResult int64) int64 {
	nOperands := len(operands)
	if nOperands < 1 {
		return 0
	}

	result := operands[0]
	nOps := len(ops)
	var err error
	for idx := range nOps {
		if result > desiredResult {
			return -1
		}

		switch ops[idx] {
		case OpAdd:
			result += operands[idx+1]
		case OpMul:
			result *= operands[idx+1]
		case OpConcat:
			resultStr := strconv.FormatInt(result, 10)
			operandStr := strconv.FormatInt(operands[idx+1], 10)
			result, err = strconv.ParseInt(resultStr+operandStr, 10, 64)
			if err != nil {
				log.Panicf("internal error: could not convert `%v` to int64", operandStr) //nolint:revive // Toy code
			}
		case NumOfDiffOperators:
			log.Panicf("internal error: unknown operator %d", ops[idx]) //nolint:revive // Toy code
		}
	}

	return result
}

func calcBack(operands []int64, ops []Operator, desiredResult int64) int64 {
	nOperands := len(operands)
	if nOperands < 1 {
		return 0
	}

	result := desiredResult
	nOps := len(ops)
	var err error
	for idx := range nOps {
		if result < 1 {
			return -1
		}

		switch ops[idx] {
		case OpAdd:
			result -= operands[nOperands-idx-1]
		case OpMul:
			if result%operands[nOperands-idx-1] != 0 {
				return -1
			}
			result /= operands[nOperands-idx-1]
		case OpConcat:
			resultStr := strconv.FormatInt(result, 10)
			operandStr := strconv.FormatInt(operands[nOperands-idx-1], 10)
			if !strings.HasSuffix(resultStr, operandStr) {
				return -1
			}
			trimmed := strings.TrimSuffix(resultStr, operandStr)
			if len(trimmed) < 1 {
				return -1
			}
			result, err = strconv.ParseInt(trimmed, 10, 64)
			if err != nil {
				log.Panicf("internal error: could not convert `%v` to int64", trimmed) //nolint:revive // Toy code
			}
		case NumOfDiffOperators:
			log.Panicf("internal error: unknown operator %d", ops[idx]) //nolint:revive // Toy code
		}
	}

	return result
}
