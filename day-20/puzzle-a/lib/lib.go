package lib

import (
	"bufio"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-set/v3"
)

type Coord struct {
	Row int
	Col int
}

func (c Coord) Add(other Coord) Coord {
	return Coord{
		Row: c.Row + other.Row,
		Col: c.Col + other.Col,
	}
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

type DirKey int

const (
	UpDir DirKey = iota
	DownDir
	LeftDir
	RightDir
)

var Directions = map[DirKey]Coord{ //nolint:gochecknoglobals // Meant as a constant
	UpDir:    {Row: -1, Col: 0},
	DownDir:  {Row: 1, Col: 0},
	LeftDir:  {Row: 0, Col: -1},
	RightDir: {Row: 0, Col: 1},
}

type Cell int

const (
	Empty Cell = iota
	Wall
)

type Cost int64

type Maze struct {
	Board         [][]Cell
	Start         *Coord
	End           *Coord
	Pos           Coord
	Dimensions    Coord
	CheatCost     Cost
	HasCheated    bool
	BlockedCheats *set.Set[Coord]
}

type MoveFunc func(Coord) Coord

func DirKeyToMoveFunc(dirKey DirKey) MoveFunc {
	return func(coord Coord) Coord {
		return coord.Add(Directions[dirKey])
	}
}

type PrecondFunc func(Maze) bool

func DirKeyToPrecondFunc(dirKey DirKey) PrecondFunc {
	return func(maze Maze) bool {
		newPos := maze.Pos.Add(Directions[dirKey])
		if !newPos.IsValid(maze.Dimensions) {
			return false
		}

		if maze.Board[newPos.Row][newPos.Col] == Wall && maze.HasCheated {
			return false
		}

		return true
	}
}

type Move struct {
	Precondition PrecondFunc
	Func         MoveFunc
}

var Moves = []Move{ //nolint:gochecknoglobals // Meant as a constant
	{DirKeyToPrecondFunc(UpDir), DirKeyToMoveFunc(UpDir)},
	{DirKeyToPrecondFunc(DownDir), DirKeyToMoveFunc(DownDir)},
	{DirKeyToPrecondFunc(LeftDir), DirKeyToMoveFunc(LeftDir)},
	{DirKeyToPrecondFunc(RightDir), DirKeyToMoveFunc(RightDir)},
}

type State struct {
	Pos        Coord
	HasCheated bool
}

func ReadInput(scanner *bufio.Scanner) (*Maze, error) {
	maze := &Maze{}
	for scanner.Scan() {
		line := scanner.Text()
		if len(strings.TrimSpace(line)) < 1 {
			break
		}

		row := make([]Cell, len(line))
		for iCol, char := range line {
			coord := Coord{Row: len(maze.Board), Col: iCol}
			switch char {
			case 'S':
				if maze.Start != nil {
					return nil, errors.New("multiple starting points found")
				}
				maze.Start = &Coord{Row: coord.Row, Col: coord.Col}
				row[iCol] = Empty
			case '.':
				row[iCol] = Empty
			case '#':
				row[iCol] = Wall
			case 'E':
				if maze.End != nil {
					return nil, errors.New("multiple ending points found")
				}
				maze.End = &Coord{Row: coord.Row, Col: coord.Col}
				maze.End.Row = coord.Row
				maze.End.Col = coord.Col
				row[iCol] = Empty

			default:
				return nil, fmt.Errorf("unrecognized cell character: `%c`", char)
			}
		}

		maze.Board = append(maze.Board, row)
		maze.Dimensions = Coord{Row: len(maze.Board), Col: len(row)}
	}

	if maze.Start == nil {
		return nil, errors.New("no starting point found")
	}
	if maze.End == nil {
		return nil, errors.New("no ending point found")
	}
	maze.Pos = *maze.Start

	maze.HasCheated = false

	return maze, nil
}
