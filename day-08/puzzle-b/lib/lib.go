package lib

import (
	"bufio"
)

type Coord struct {
	Row int
	Col int
}

func (c Coord) IsValid(dimensions Coord) bool {
	if c.Row < 0 || c.Row >= dimensions.Row {
		return false
	}
	if c.Col < 0 || c.Col >= dimensions.Col {
		return false
	}
	return true
}

func (c Coord) Add(other Coord) Coord {
	return Coord{
		Row: c.Row + other.Row,
		Col: c.Col + other.Col,
	}
}

func (c Coord) Subtract(other Coord) Coord {
	return Coord{
		Row: c.Row - other.Row,
		Col: c.Col - other.Col,
	}
}

func ReadArray(scanner *bufio.Scanner) (Coord, map[rune][]Coord) {
	nRows := 0
	nCols := -1
	antennae := make(map[rune][]Coord)
	for scanner.Scan() {
		line := scanner.Text()
		if nCols < 0 {
			nCols = len(line)
		}
		for iCol, char := range line {
			if char == '.' {
				continue
			}
			currentCoords := Coord{Row: nRows, Col: iCol}
			antennae[char] = append(antennae[char], currentCoords)
		}
		nRows++
	}

	return Coord{nRows, nCols}, antennae
}
