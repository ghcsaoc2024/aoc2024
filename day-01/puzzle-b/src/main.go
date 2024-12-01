package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
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

	var slice1, slice2 []int

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		numbers := strings.Fields(line)

		if len(numbers) >= 2 {
			num1, err1 := strconv.Atoi(numbers[0])
			num2, err2 := strconv.Atoi(numbers[1])

			if err1 == nil && err2 == nil {
				slice1 = append(slice1, num1)
				slice2 = append(slice2, num2)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Panic(err)
	}

	if len(slice1) != len(slice2) {
		log.Panic("Slices are not of the same length")
	}

	valCounts := make(map[int]int)
	for i := range slice2 {
		valCounts[slice2[i]]++
	}

	total := 0
	for i := range slice1 {
		nOccurrences := valCounts[slice1[i]]
		total += slice1[i] * nOccurrences
	}

	log.Println(total)
}
