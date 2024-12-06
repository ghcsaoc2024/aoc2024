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
	array, initialCoords := lib.ReadArray(scanner)

	dimensions := lib.Coord{Row: len(array), Col: len(array[0])}
	if dimensions.Row < 1 {
		log.Panic("no rows in input")
	}

	log.Printf("finished reading array (%d rows)", dimensions.Row)
	log.Printf("initial coordinates: %v", initialCoords)

	// Do the walkabout
	nVisited := walkabout(initialCoords, dimensions, array)

	log.Printf("visited %d cells", nVisited)
}

func walkabout(initialCoords, dimensions lib.Coord, array [][]lib.Cell) int {
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
	return nVisited
}
