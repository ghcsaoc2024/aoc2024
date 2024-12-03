package main

import (
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/samber/lo"
)

func main() {
	// Read the entire contents of "../input/input.txt" into a buffer.
	data, err := os.ReadFile("../input/input.txt")
	if err != nil {
		log.Panicf("Failed to read file: %v", err)
	}

	pattern := `(mul\(([1-9][0-9]*),([1-9][0-9]*)\))|(do\(\))|(don't\(\))`
	re, err := regexp.Compile(pattern)
	if err != nil {
		log.Panicf("Failed to compile regex: %v", err)
	}

	// Find all matches in the data.
	doIsOn := true
	runningSum := int64(0)
	matches := re.FindAllSubmatch(data, -1)
	for _, match := range matches {
		matchLen := len(match)
		switch {
		case match[1] != nil:
			prod, err := mul(match[2 : matchLen-2])
			if err != nil {
				log.Panic(err)
			}
			if doIsOn {
				runningSum += prod
			}
		case match[matchLen-2] != nil:
			doIsOn = true
		case match[matchLen-1] != nil:
			doIsOn = false
		default:
			log.Panicf("Internal error: unexpected kind of match `%v`", string(match[0]))
		}
	}

	log.Println(runningSum)
}

func mul(match [][]byte) (int64, error) {
	var firstError error
	operands := lo.Map(match, func(item []byte, _ int) int64 {
		num, err := strconv.ParseInt(string(item), 10, 64)
		if err != nil {
			if firstError == nil {
				firstError = err
			}
			return 1
		}
		return num
	})

	if firstError != nil {
		return 1, firstError
	}

	return lo.Reduce(operands, func(prodSoFar, operand int64, _ int) int64 {
		return prodSoFar * operand
	}, 1), nil
}
