package main

import (
	"bufio"
	"log"
	"os"
	"slices"

	"daysix/lib"

	"github.com/tiendc/go-deepcopy"
)

func main() {
	file, err := os.Open("../input/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		closeErr := file.Close()
		if closeErr != nil {
			log.Fatal(closeErr)
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

	// Do an initial walkabout to determine which coordinates are visited *without*
	// blocking any additional cells
	walkabout(initialCoords, dimensions, array)

	initialVisitationArray := make([][]lib.Cell, 0)
	err = deepcopy.Copy(&initialVisitationArray, array)
	if err != nil {
		log.Panic(err)
	}

	// Now, for each visited cell, check if blocking it would create a loop
	nLoopifiers := 0
	for row := range dimensions.Row {
		for col := range dimensions.Col {
			currentCoords := lib.Coord{Row: row, Col: col}
			if currentCoords == initialCoords {
				continue
			}

			if array[currentCoords.Row][currentCoords.Col] == lib.Blocked {
				continue
			}

			if initialVisitationArray[currentCoords.Row][currentCoords.Col] != lib.Visited {
				continue
			}

			array[currentCoords.Row][currentCoords.Col] = lib.Blocked
			if isLoopful(initialCoords, dimensions, array) {
				nLoopifiers++
			}
			array[currentCoords.Row][currentCoords.Col] = lib.Empty
		}
	}

	log.Printf("found %d loopifiers", nLoopifiers)
}

func walkabout(initialCoords, dimensions lib.Coord, array [][]lib.Cell) (lib.Coord, int) {
	initialDir := lib.Coord{Row: -1, Col: 0}
	currentCoords := initialCoords
	currentDir := initialDir
	nVisited := 1
	timesReset := 0
	for {
		if currentCoords == initialCoords && currentDir == initialDir {
			timesReset++
		}

		if timesReset > 1 {
			log.Panic("we're in a loop!") //nolint:revive // Toy code
		}

		nextCoords := currentCoords.MoveOne(currentDir)
		if !nextCoords.IsValid(dimensions) {
			break
		}

		switch array[nextCoords.Row][nextCoords.Col] {
		case lib.Empty:
			array[nextCoords.Row][nextCoords.Col] = lib.Visited
			nVisited++
			fallthrough
		case lib.Visited:
			currentCoords = nextCoords
			continue
		case lib.Blocked:
			currentDir = lib.TurnRight(currentDir)
		}
	}
	return currentCoords, nVisited
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

		nextCoords := currentCoords.MoveOne(currentDir)
		if !nextCoords.IsValid(dimensions) {
			return false
		}

		switch array[nextCoords.Row][nextCoords.Col] {
		case lib.Empty:
			fallthrough
		case lib.Visited:
			currentCoords = nextCoords
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
