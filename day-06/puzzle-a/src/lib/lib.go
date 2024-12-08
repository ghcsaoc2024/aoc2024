package lib

import (
	"bufio"
	"log"
)

type Coord struct {
	Row int
	Col int
}

type Cell int

const (
	Empty Cell = iota
	Visited
	Blocked
)

func (c Coord) MoveOne(dir Coord) Coord {
	return Coord{
		Row: c.Row + dir.Row,
		Col: c.Col + dir.Col,
	}
}

func (c Coord) IsValid(ddimensions Coord) bool {
	if c.Row < 0 || c.Row >= ddimensions.Row {
		return false
	}
	if c.Col < 0 || c.Col >= ddimensions.Col {
		return false
	}
	return true
}

func TurnRight(dir Coord) Coord {
	return Coord{
		Row: dir.Col,
		Col: -dir.Row,
	}
}

func ReadArray(scanner *bufio.Scanner) ([][]Cell, Coord) {
	array := make([][]Cell, 0)
	initialCoords := Coord{Row: -1, Col: -1}
	for scanner.Scan() {
		line := scanner.Text()
		row := make([]Cell, 0)
		for _, char := range line {
			currentCoords := Coord{Row: len(array), Col: len(row)}
			switch char {
			case '^':
				if initialCoords != (Coord{Row: -1, Col: -1}) {
					log.Panicf("multiple starting points found: had already encountered %v, and now encountered %v", initialCoords, currentCoords) //nolint:revive // Toy code
				}
				initialCoords = currentCoords
				row = append(row, Visited)
			case '.':
				row = append(row, Empty)
			case '#':
				row = append(row, Blocked)
			default:
				log.Panicf("unexpected character in input: %v (current coordinates: %v)", char, currentCoords) //nolint:revive // Toy code
			}
		}
		array = append(array, row)
	}
	return array, initialCoords
}
