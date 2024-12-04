package main

import (
	"bufio"
	"log"
	"os"
)

var (
	target = []rune{'M', 'A', 'S'}                       //nolint:gochecknoglobals // Meant as a constant
	diags  = [][]int{{1, 1}, {1, -1}, {-1, 1}, {-1, -1}} //nolint:gochecknoglobals // Meant as a constant
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
			for dirIdx1, dir1 := range diags {
				for dirIdx2 := dirIdx1 + 1; dirIdx2 < len(diags); dirIdx2++ {
					dir2 := diags[dirIdx2]
					if searchInTwoDirs(array, i, j, dir1, dir2) {
						nFound++
					}
				}
			}
		}
	}

	return nFound
}

func searchInTwoDirs(array [][]rune, i, j int, dir1, dir2 []int) bool {
	return searchInDir(array, i-dir1[0], j-dir1[1], dir1) && searchInDir(array, i-dir2[0], j-dir2[1], dir2)
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
