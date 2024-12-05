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
	byPreceder := make(map[int][]int)
	byFollower := make(map[int][]int)
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

		byPreceder[values[0]] = append(byPreceder[values[0]], values[1])
		byFollower[values[1]] = append(byFollower[values[1]], values[0])
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

		if !isValid(values, byPreceder, byFollower) {
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

func isValid(values []int, byPreceder, byFollower map[int][]int) bool {
	nValues := len(values)
	posByValue := make(map[int]int)
	for pos, value := range values {
		posByValue[value] = pos
	}

	// Validate obligatory followers of each value
	for idx := 1; idx < nValues; idx++ {
		obligFollowers := byPreceder[values[idx]]
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

	// Validate obligatory preceders of each value
	for idx := nValues - 2; idx >= 0; idx-- {
		obligPreceders := byFollower[values[idx]]
		for _, preceder := range obligPreceders {
			pos, doesItOccur := posByValue[preceder]
			if !doesItOccur {
				continue
			}

			if pos > idx {
				return false
			}
		}
	}

	return true
}
