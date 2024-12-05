package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/samber/lo"
)

func main() {
	file, err := os.Open("../input/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	scanner := bufio.NewScanner(file)

	// Read in the precedence rules
	precedenceMap := make(map[int][]int)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "|")
		if len(fields) != 2 {
			break
		}

		values := lo.Map(fields, func(item string, _ int) int {
			num, err := strconv.Atoi(item)
			if err != nil {
				log.Panic(err)
			}
			return num
		})

		precedenceMap[values[0]] = append(precedenceMap[values[0]], values[1])
	}

	// Read in the page sets
	runningTotal := 0
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, ",")
		values := lo.Map(fields, func(item string, _ int) int {
			num, err := strconv.Atoi(item)
			if err != nil {
				log.Panic(err)
			}
			return num
		})

		if isValid(values, precedenceMap) {
			continue
		}

		nValues := len(values)
		if nValues%2 == 0 {
			continue
		}

		repairValues(&values, precedenceMap)
		if !isValid(values, precedenceMap) {
			log.Printf("could not repair %v", values)
			continue
		}

		middle := nValues / 2
		runningTotal += values[middle]
	}

	log.Println(runningTotal)
}

func isValid(values []int, precedenceMap map[int][]int) bool {
	posByValue := genPosByValueMap(values)

	// Validate obligatory followers of each value
	nValues := len(values)
	for idx := 1; idx < nValues; idx++ {
		obligFollowers := precedenceMap[values[idx]]
		for _, follower := range obligFollowers {
			pos, doesItOccur := posByValue[follower]
			if doesItOccur && pos < idx {
				return false
			}
		}
	}

	return true
}

func repairValues(values *[]int, precedenceMap map[int][]int) {
	posByValue := genPosByValueMap(*values)

	// Validate obligatory followers of each value
	nValues := len(*values)
	for idx := 1; idx < nValues; idx++ {
		currentVal := (*values)[idx]
		obligFollowers := precedenceMap[currentVal]
		for _, follower := range obligFollowers {
			pos, doesItOccur := posByValue[follower]
			if doesItOccur && pos < idx {
				(*values)[idx], (*values)[pos] = (*values)[pos], (*values)[idx]
				posByValue[follower], posByValue[currentVal] = idx, pos
				idx = pos - 1
				break
			}
		}
	}
}

func genPosByValueMap(values []int) map[int]int {
	return lo.FromEntries(lo.Map(values, func(value, pos int) lo.Entry[int, int] {
		return lo.Entry[int, int]{Key: value, Value: pos}
	}))
}
