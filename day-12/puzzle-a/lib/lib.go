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

var Directions = []Coord{ //nolint:gochecknoglobals // Meant as a constant
	{Row: 1, Col: 0},
	{Row: 0, Col: 1},
	{Row: 0, Col: -1},
	{Row: -1, Col: 0},
}

type Cell struct {
	Kind       rune
	Coord      Coord
	Boundaries [4]bool
}

func ReadInput(scanner *bufio.Scanner) [][]Cell {
	board := make([][]Cell, 0)
	for scanner.Scan() {
		line := scanner.Text()
		iRow := len(board)
		row := make([]Cell, 0)
		for iCol, c := range line {
			coord := Coord{Row: iRow, Col: iCol}
			row = append(row, Cell{Kind: c, Coord: coord, Boundaries: [4]bool{true, true, true, true}})
		}
		board = append(board, row)
	}

	return board
}
