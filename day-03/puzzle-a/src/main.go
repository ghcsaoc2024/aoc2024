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

	// Create a regular expression that matches anything starting with mul
	pattern := `mul\(([1-9][0-9]*),([1-9][0-9]*)\)`
	re, err := regexp.Compile(pattern)
	if err != nil {
		log.Panicf("Failed to compile regex: %v", err)
	}

	// Find all matches in the data.
	matches := re.FindAllSubmatch(data, -1)
	sum := lo.Sum(lo.Map(matches, func(match [][]byte, _ int) int64 {
		operands := lo.Map(match[1:], func(item []byte, _ int) int64 {
			num, err := strconv.ParseInt(string(item), 10, 64)
			if err != nil {
				log.Panic(err)
			}
			return num
		})

		if len(operands) != 2 {
			log.Panic("Expected two operands")
		}

		return operands[0] * operands[1]
	}))

	log.Println(sum)
}
