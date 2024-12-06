package main

import (
	"bufio"
	"log"
	"os"

	"daysix/lib"

	"github.com/hashicorp/go-set/v3"
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
	array, initialCoords := lib.ReadArray(scanner)

	dimensions := lib.Coord{Row: len(array), Col: len(array[0])}
	if dimensions.Row < 1 {
		log.Panic("no rows in input")
	}

	if !initialCoords.IsValid(dimensions) {
		log.Panicf("initial coordinates %v are not valid (dimensions: %v)", initialCoords, dimensions)
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

func walkabout(initialCoords, dimensions lib.Coord, array [][]lib.Cell) {
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
}

func isLoopful(initialCoords, dimensions lib.Coord, array [][]lib.Cell) bool {
	initialDir := lib.Coord{Row: -1, Col: 0}
	current := lib.Visitation{
		Loc: initialCoords,
		Dir: initialDir,
	}
	turns := set.New[lib.Visitation](0)
	for {
		nextCoords := current.Loc.MoveOne(current.Dir)
		if !nextCoords.IsValid(dimensions) {
			return false
		}

		switch array[nextCoords.Row][nextCoords.Col] {
		case lib.Empty:
			fallthrough
		case lib.Visited:
			current.Loc = nextCoords
			continue
		case lib.Blocked:
			if turns.Contains(current) {
				return true
			}
			turns.Insert(current)
			current.Dir = lib.TurnRight(current.Dir)
		}
	}
}
