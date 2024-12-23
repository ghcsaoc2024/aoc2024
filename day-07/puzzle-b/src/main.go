// This is the one from Advent of Code 2024 with the simple demonstration
// of implementing a WorkerPool pattern in Go.

package main

import (
	"bufio"
	"log"
	"math"
	"math/big"
	"os"
	"regexp"
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

const MaxNumWorkers = 65536

type WorkerTask struct {
	result    big.Int
	operands  []big.Int
	sweetSpot float64
}

type Args struct {
	InputFile  string  `arg:"positional,required" help:"input file"`
	SweetSpot  float64 `arg:"positional,required" help:"sweet spot for meet-in-the-\"middle\""`
	NumWorkers int     `arg:"-n"                  default:"1"                                  help:"number of workers to use"`
}

func main() {
	var args Args
	arg.MustParse(&args)

	if args.SweetSpot <= 0 || args.SweetSpot >= 1 {
		log.Fatalf("sweet spot must be larger than 0.0 and smaller than 1.0; got %f", args.SweetSpot)
	}

	if args.NumWorkers < 1 || args.NumWorkers > MaxNumWorkers {
		log.Fatalf("number of workers must be at least 1 and no more than %d; got %d", MaxNumWorkers, args.NumWorkers)
	}

	file, err := os.Open(args.InputFile)
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
	maxAttainable := big.NewInt(0)
	runningTotal := big.NewInt(0)

	// Task and result channels
	taskChan := make(chan WorkerTask, args.NumWorkers)
	resultsChan := make(chan big.Int)

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

			var result big.Int
			_, success := result.SetString(matches[1], 10) //nolint:mnd // false positive
			if !success {
				log.Panicf("internal error: could not convert `%v` to big.Int", matches[1])
			}

			operandStrings := strings.Fields(matches[2])
			operands := lo.Map(operandStrings, func(s string, _ int) big.Int {
				var operand big.Int
				_, success := operand.SetString(s, 10) //nolint:mnd // false positive
				if !success {
					log.Panicf("internal error: could not convert `%v` to big.Int", s)
				}
				return operand
			})

			maxAttainable.Add(maxAttainable, &result)

			// Send task to the worker pool
			taskChan <- WorkerTask{result: result, operands: operands, sweetSpot: args.SweetSpot}
		}
		close(taskChan)
	}()

	go func() {
		wg.Wait()
		close(resultsChan)
		// Contrary to appearances, this is in fact hermetically sealed:
		// because resultsChan has no buffer, and individual workers don't call
		// `wg.Done()` until *after* their call to `resultsChan <- task.result`
		// (which will block unless and until something reads the `task.result`
		// from the channel), it follows that there is no way `resultsChan` will
		// be closed by this deferred func *before* the result has the chance to
		// have been read by the loop below.
	}()

	for result := range resultsChan {
		runningTotal.Add(runningTotal, &result)
	}

	maxAttainableStr := maxAttainable.String()
	log.Printf("max attainable:  %s", maxAttainableStr)
	log.Printf("running total:   %*v", len(maxAttainableStr), runningTotal)
}

func worker(taskChan <-chan WorkerTask, resultsChan chan<- big.Int, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range taskChan {
		if isSolvable(task.result, task.operands, task.sweetSpot) {
			resultsChan <- task.result
		}
	}
}

func isSolvable(desiredResult big.Int, operands []big.Int, sweetSpot float64) bool {
	nOperands := len(operands)
	nOperators := nOperands - 1
	middleOpIdx := int(math.Round(float64(nOperators) * sweetSpot))
	semiSolutions := set.New[string](0)

	var operators []Operator
	var nCombinations int64

	// First half
	nFirstHalfOperators := middleOpIdx
	operators = make([]Operator, nFirstHalfOperators)
	nCombinations = int64(math.Pow(float64(NumOfDiffOperators), float64(nFirstHalfOperators)))
	for iCombo := range nCombinations {
		combo := iCombo
		for iOperator := range nFirstHalfOperators {
			operators[iOperator] = Operator(combo % int64(NumOfDiffOperators))
			combo /= int64(NumOfDiffOperators)
		}

		result := calcFwd(operands, operators, desiredResult)
		if result == nil {
			continue
		}
		semiSolutions.Insert(result.String())
	}

	// Second half
	nSecondHalfOperators := nOperators - middleOpIdx
	operators = make([]Operator, nSecondHalfOperators)
	nCombinations = int64(math.Pow(float64(NumOfDiffOperators), float64(nSecondHalfOperators)))
	for iCombo := range nCombinations {
		combo := iCombo
		for iOperator := range nSecondHalfOperators {
			operators[iOperator] = Operator(combo % int64(NumOfDiffOperators))
			combo /= int64(NumOfDiffOperators)
		}

		result := calcBack(operands, operators, desiredResult)
		if result == nil {
			continue
		}

		if semiSolutions.Contains(result.String()) {
			return true
		}
	}

	return false
}

func calcFwd(operands []big.Int, operators []Operator, desiredResult big.Int) *big.Int {
	nOperands := len(operands)
	if nOperands < 1 {
		return nil
	}

	var result big.Int
	result.Set(&operands[0])
	nOperators := len(operators)
	for iOperator := range nOperators {
		if result.Cmp(&desiredResult) > 0 {
			return nil
		}

		operand := operands[iOperator+1]
		switch operators[iOperator] {
		case OpAdd:
			result.Add(&result, &operand)
		case OpMul:
			result.Mul(&result, &operand)
		case OpConcat:
			resultStr := result.String()
			operandStr := operand.String()
			_, success := result.SetString(resultStr+operandStr, 10) //nolint:mnd // false positive
			if !success {
				log.Panicf("internal error: could not convert `%v` to big.Int", operandStr) //nolint:revive // Toy code
			}
		case NumOfDiffOperators:
			log.Panicf("internal error: unknown operator %d", operators[iOperator]) //nolint:revive // Toy code
		}
	}

	return &result
}

func calcBack(operands []big.Int, operators []Operator, desiredResult big.Int) *big.Int {
	nOperands := len(operands)
	if nOperands < 1 {
		return nil
	}

	var result big.Int
	result.Set(&desiredResult)
	nOperators := len(operators)
	for iOperator := range nOperators {
		if result.Sign() < 1 {
			return nil
		}

		operand := operands[nOperands-iOperator-1]
		switch operators[iOperator] {
		case OpAdd:
			result.Sub(&result, &operand)
		case OpMul:
			var remainder big.Int
			result.QuoRem(&result, &operand, &remainder)
			if remainder.Sign() != 0 {
				return nil
			}
		case OpConcat:
			resultStr := result.String()
			operandStr := operand.String()
			if !strings.HasSuffix(resultStr, operandStr) {
				return nil
			}
			trimmed := strings.TrimSuffix(resultStr, operandStr)
			if len(trimmed) < 1 {
				return nil
			}
			_, success := result.SetString(trimmed, 10) //nolint:mnd // false positive
			if !success {
				log.Panicf("internal error: could not convert `%v` to big.Int", trimmed) //nolint:revive // Toy code
			}
		case NumOfDiffOperators:
			log.Panicf("internal error: unknown operator %d", operators[iOperator]) //nolint:revive // Toy code
		}
	}

	return &result
}
