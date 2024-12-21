package lib

import (
	"bufio"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/exp/constraints"
)

func Abs[T constraints.Signed](x T) T { //nolint:ireturn // false positive
	return max(x, -x)
}

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

func Distance(from, to Coord) int {
	return Abs(from.Row-to.Row) + Abs(from.Col-to.Col)
}

type dirKey int

const (
	UpDir dirKey = iota
	DownDir
	LeftDir
	RightDir
)

var directions = map[dirKey]Coord{ //nolint:gochecknoglobals // Meant as a constant
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
	Board      [][]Cell
	Start      *Coord
	End        *Coord
	Dimensions Coord
	Moves      []Move
}

type State struct {
	Pos Coord
}

type moveFunc func(State) (State, Cost)

func dirKeyToMoveFunc(dirKey dirKey) moveFunc {
	return func(oldState State) (State, Cost) {
		newState := oldState
		newState.Pos = oldState.Pos.Add(directions[dirKey])

		return newState, 1
	}
}

type precondFunc func(Maze, State) bool

func dirKeyToPrecondFunc(dirKey dirKey) precondFunc {
	return func(maze Maze, curState State) bool {
		newPos := curState.Pos.Add(directions[dirKey])
		if !newPos.IsValid(maze.Dimensions) {
			return false
		}

		newCell := maze.Board[newPos.Row][newPos.Col]

		return newCell != Wall
	}
}

type Move struct {
	Precondition precondFunc
	Func         moveFunc
}

var moves = []Move{ //nolint:gochecknoglobals // Meant as a constant
	{dirKeyToPrecondFunc(UpDir), dirKeyToMoveFunc(UpDir)},
	{dirKeyToPrecondFunc(DownDir), dirKeyToMoveFunc(DownDir)},
	{dirKeyToPrecondFunc(LeftDir), dirKeyToMoveFunc(LeftDir)},
	{dirKeyToPrecondFunc(RightDir), dirKeyToMoveFunc(RightDir)},
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

	maze.Moves = moves

	return maze, nil
}
