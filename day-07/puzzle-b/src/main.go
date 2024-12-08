package main

import (
	"bufio"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"

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

type WorkerTask struct {
	result    int64
	operands  []int64
	sweetSpot float64
}

func worker(tasks <-chan WorkerTask, resultsChan chan<- int64, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range tasks {
		if isSolvable(task.result, task.operands, task.sweetSpot) {
			resultsChan <- task.result
		}
	}
}

func main() {
	var args struct {
		SweetSpot  float64 `arg:"positional,required" help:"sweet spot for meet-in-the-\"middle\""`
		NumWorkers int     `arg:"-n,env:NUM_WORKERS"  default:"1"                                  help:"number of workers to use"`
	}
	arg.MustParse(&args)

	if args.SweetSpot <= 0 || args.SweetSpot >= 1 {
		log.Fatalf("sweet spot must be larger than 0.0 and smaller than 1.0; got %f", args.SweetSpot)
	}

	file, err := os.Open("../input/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		closeErr := file.Close()
		if closeErr != nil {
			log.Fatal(closeErr)
		}
	}()

	pattern := `^([1-9][0-9]*):((\s+([1-9][0-9]*))+)\s*$`
	re, err := regexp.Compile(pattern)
	if err != nil {
		log.Panicf("internal error: failed to compile regex: `%v`", err)
	}

	scanner := bufio.NewScanner(file)
	iLine := 0
	maxAttainable := int64(0)
	runningTotal := int64(0)

	// Task and result channels
	taskChan := make(chan WorkerTask, args.NumWorkers)
	resultsChan := make(chan int64)

	// Setup worker pool
	var wg sync.WaitGroup
	wg.Add(args.NumWorkers)

	for range args.NumWorkers {
		go worker(taskChan, resultsChan, &wg)
	}

	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			iLine++
			matches := re.FindStringSubmatch(line)
			if len(matches) == 0 {
				log.Printf("could not match line %d (`%v`)", iLine, line)
				continue
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

			// Send task to the worker pool
			taskChan <- WorkerTask{result: result, operands: operands, sweetSpot: args.SweetSpot}
		}
		close(taskChan)
	}()

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	for result := range resultsChan {
		runningTotal += result
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

		operand := operands[idx+1]
		switch ops[idx] {
		case OpAdd:
			result += operand
		case OpMul:
			result *= operand
		case OpConcat:
			resultStr := strconv.FormatInt(result, 10)
			operandStr := strconv.FormatInt(operand, 10)
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

		operand := operands[nOperands-idx-1]
		switch ops[idx] {
		case OpAdd:
			result -= operand
		case OpMul:
			if result%operand != 0 {
				return -1
			}
			result /= operand
		case OpConcat:
			resultStr := strconv.FormatInt(result, 10)
			operandStr := strconv.FormatInt(operand, 10)
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
