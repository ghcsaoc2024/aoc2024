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

		if !isValid(values, precedenceMap) {
			continue
		}

		nValues := len(values)
		if nValues%2 == 0 {
			log.Printf("value list %v is not of odd length", values)
			continue
		}

		middle := nValues / 2
		runningTotal += values[middle]
	}

	log.Println(runningTotal)
}

func isValid(values []int, precedenceMap map[int][]int) bool {
	nValues := len(values)
	posByValue := make(map[int]int)
	for pos, value := range values {
		posByValue[value] = pos
	}

	// Validate obligatory followers of each value
	for idx := 1; idx < nValues; idx++ {
		obligFollowers := precedenceMap[values[idx]]
		for _, follower := range obligFollowers {
			pos, doesItOccur := posByValue[follower]
			if !doesItOccur {
				continue
			}

			if pos < idx {
				return false
			}
		}
	}

	return true
}
