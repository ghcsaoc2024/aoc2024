package main

import (
	"bufio"
	"log"
	"os"
	"slices"

	"daysix/lib"
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

	// Read in the array
	array, initialCoords := readArray(scanner)

	dimensions := lib.Coord{Row: len(array), Col: len(array[0])}
	if dimensions.Row < 1 {
		log.Panic("no rows in input")
	}

	log.Printf("finished reading array (%d rows)", dimensions.Row)
	log.Printf("initial coordinates: %v", initialCoords)

	// Do the walkabouts
	nLoopifiers := 0
	for row := range dimensions.Row {
		log.Printf("row %d", row)
		for col := range dimensions.Col {
			currentCoords := lib.Coord{Row: row, Col: col}
			if currentCoords == initialCoords {
				continue
			}

			if array[currentCoords.Row][currentCoords.Col] == lib.Blocked {
				continue
			}

			array[currentCoords.Row][currentCoords.Col] = lib.Blocked
			if isLoopful(initialCoords, dimensions, array) {
				nLoopifiers++
			}
			array[currentCoords.Row][currentCoords.Col] = lib.Empty
		}
		log.Printf("found %d loopifiers", nLoopifiers)
	}

	log.Printf("found %d loopifiers", nLoopifiers)
}

func isLoopful(initialCoords, dimensions lib.Coord, array [][]lib.Cell) bool {
	initialDir := lib.Coord{Row: -1, Col: 0}
	currentCoords := initialCoords
	currentDir := initialDir
	visitationArray := make([][][]lib.Coord, len(array))
	for iRow, row := range array {
		visitationArray[iRow] = make([][]lib.Coord, len(row))
	}

	for {
		if slices.Contains(visitationArray[currentCoords.Row][currentCoords.Col], currentDir) {
			return true
		}
		visitationArray[currentCoords.Row][currentCoords.Col] = append(visitationArray[currentCoords.Row][currentCoords.Col], currentDir)

		provisionalNextCoords := lib.NextCoords(currentCoords, currentDir)
		if !lib.IsValidCoord(provisionalNextCoords, dimensions) {
			return false
		}

		switch array[provisionalNextCoords.Row][provisionalNextCoords.Col] {
		case lib.Empty:
			currentCoords = provisionalNextCoords
			continue
		case lib.Blocked:
			currentDir = lib.TurnRight(currentDir)
		}
	}
}

func readArray(scanner *bufio.Scanner) ([][]lib.Cell, lib.Coord) {
	array := make([][]lib.Cell, 0)
	initialCoords := lib.Coord{Row: -1, Col: -1}
	for scanner.Scan() {
		line := scanner.Text()
		row := make([]lib.Cell, 0)
		for _, char := range line {
			currentCoords := lib.Coord{Row: len(array), Col: len(row)}
			switch char {
			case '^':
				if initialCoords != (lib.Coord{Row: -1, Col: -1}) {
					log.Panicf("multiple starting points found: had already encountered %v, and now encountered %v", initialCoords, currentCoords) //nolint:revive // Toy code
				}
				initialCoords = currentCoords
				fallthrough
			case '.':
				row = append(row, lib.Empty)
			case '#':
				row = append(row, lib.Blocked)
			default:
				log.Panicf("unexpected character in input: %v (current coordinates: %v)", char, currentCoords) //nolint:revive // Toy code
			}
		}
		array = append(array, row)
	}
	return array, initialCoords
}
