package lib

import (
	"bufio"

	"github.com/hashicorp/go-set/v3"
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

const TopElevation = 9
const BottomElevation = 0
const InvalidElevation = -1

type Cell struct {
	Elevation  int
	TrailCount int
}

type Board struct {
	Grid        [][]Cell
	ByElevation map[int][]Coord
}

func ReadInput(scanner *bufio.Scanner) Board {
	var board Board
	board.Grid = make([][]Cell, 0)
	board.ByElevation = make(map[int][]Coord)
	for scanner.Scan() {
		line := scanner.Text()
		iRow := len(board.Grid)
		board.Grid = append(board.Grid, make([]Cell, len(line)))
		for iCol, c := range line {
			elevation := int(c - '0')
			if elevation < BottomElevation || elevation > TopElevation {
				elevation = InvalidElevation
			}
			reachablePeaks := set.New[Coord](0)
			trailCount := 0
			if elevation == TopElevation {
				reachablePeaks.Insert(Coord{Row: iRow, Col: iCol})
				trailCount = 1
			}
			board.Grid[iRow][iCol] = Cell{Elevation: elevation, TrailCount: trailCount}
			board.ByElevation[elevation] = append(board.ByElevation[elevation], Coord{Row: iRow, Col: iCol})
		}
	}

	return board
}
