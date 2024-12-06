package main

import (
	"bufio"
	"log"
	"os"

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

	// Do the walkabout
	currentCoords, nVisited := walkabout(initialCoords, dimensions, array)

	log.Printf("current coordinates: %v", currentCoords)
	log.Printf("visited %d cells", nVisited)
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

		nextCoords := lib.NextCoords(currentCoords, currentDir)
		if !lib.IsValidCoord(nextCoords, dimensions) {
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
				row = append(row, lib.Visited)
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
