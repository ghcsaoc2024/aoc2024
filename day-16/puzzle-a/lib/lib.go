package lib

import (
	"bufio"
	"errors"
	"fmt"
	"strings"
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

func (c Coord) Sub(other Coord) Coord {
	return Coord{
		Row: c.Row - other.Row,
		Col: c.Col - other.Col,
	}
}

func (c Coord) TurnRight() Coord {
	dir := &c
	return Coord{
		Row: dir.Col,
		Col: -dir.Row,
	}
}

func (c Coord) TurnLeft() Coord {
	dir := &c
	return Coord{
		Row: -dir.Col,
		Col: dir.Row,
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

var StartDirection = Coord{Row: 0, Col: 1} //nolint:gochecknoglobals // Meant as a constant

type Cell int

const (
	Empty Cell = iota
	Wall
)

type Cursor struct {
	Coord
	Dir Coord
}

type Cost int64

type Maze struct {
	Board      [][]Cell
	Start      *Coord
	End        *Coord
	Cursor     Cursor
	Dimensions Coord
	Cost       Cost
	Solved     bool
}

func Forward(m Maze) Maze {
	m.Cursor.Coord = m.Cursor.Coord.Add(m.Cursor.Dir)
	m.Cost += 1
	return m
}

func TurnRight(m Maze) Maze {
	m.Cursor.Dir = m.Cursor.Dir.TurnRight()
	m.Cost += 1000
	return m
}

func TurnLeft(m Maze) Maze {
	m.Cursor.Dir = m.Cursor.Dir.TurnLeft()
	m.Cost += 1000
	return m
}

type MoveFunc func(Maze) Maze

type Move struct {
	Precondition func(Maze) bool
	Func         MoveFunc
}

func alwaysTruePrecondition(_ Maze) bool {
	return true
}

func forwardPrecondition(maze Maze) bool {
	nextCoord := maze.Cursor.Coord.Add(maze.Cursor.Dir)
	if !nextCoord.IsValid(maze.Dimensions) {
		return false
	}

	if maze.Board[nextCoord.Row][nextCoord.Col] == Wall {
		return false
	}

	return true
}

var Moves = []Move{ //nolint:gochecknoglobals // Meant as a constant
	{forwardPrecondition, Forward},
	{alwaysTruePrecondition, TurnRight},
	{alwaysTruePrecondition, TurnLeft},
}

func ReadInput(scanner *bufio.Scanner) (*Maze, error) {
	maze := &Maze{}
	for scanner.Scan() {
		line := scanner.Text()
		if len(strings.TrimSpace(line)) < 1 {
			break
		}

		row := make([]Cell, len(line)*2)
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

	maze.Cursor = Cursor{*maze.Start, StartDirection}
	maze.Cost = 0
	maze.Solved = false

	return maze, nil
}
