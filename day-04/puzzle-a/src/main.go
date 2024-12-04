package main

import (
	"bufio"
	"log"
	"os"
)

var (
	target     = []rune{'X', 'M', 'A', 'S'}                                                    //nolint:gochecknoglobals // Meant as a constant
	directions = [][]int{{1, 1}, {1, 0}, {1, -1}, {0, 1}, {0, -1}, {-1, 1}, {-1, 0}, {-1, -1}} //nolint:gochecknoglobals // Meant as a constant
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
	array := make([][]rune, 0)
	for scanner.Scan() {
		line := scanner.Text()
		runeSlice := []rune(line)
		array = append(array, runeSlice)
	}

	nFound := doSearch(array)

	log.Println(nFound)
}

func doSearch(array [][]rune) int {
	nFound := 0
	for i := range array {
		for j := range array[i] {
			for _, dir := range directions {
				if searchInDir(array, i, j, dir) {
					nFound++
				}
			}
		}
	}

	return nFound
}

func searchInDir(array [][]rune, i, j int, dir []int) bool {
	for step := range target {
		x := i + dir[0]*step
		y := j + dir[1]*step
		if x < 0 || x >= len(array) || y < 0 || y >= len(array[x]) {
			return false
		}
		if array[x][y] != target[step] {
			return false
		}
	}

	return true
}
