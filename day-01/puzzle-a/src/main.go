package main

import (
	"bufio"
	"log"
	"os"
	"sort"
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

	sort.Ints(slice1)
	sort.Ints(slice2)

	distance := 0
	for i := range slice1 {
		diff := slice1[i] - slice2[i]
		if diff > 0 {
			distance += diff
		} else {
			distance -= diff
		}
	}

	log.Println(distance)
}
