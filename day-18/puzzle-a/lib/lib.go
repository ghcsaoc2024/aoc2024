package lib

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"github.com/samber/lo"
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

func (c Coord) Mul(scalar int) Coord {
	return Coord{
		Row: c.Row * scalar,
		Col: c.Col * scalar,
	}
}

func (c Coord) IsValid(dims Coord) bool {
	if c.Row < 0 || c.Row >= dims.Row {
		return false
	}
	if c.Col < 0 || c.Col >= dims.Col {
		return false
	}
	return true
}

var Directions = []Coord{ //nolint:gochecknoglobals // Meant as a constant
	{Row: 0, Col: 1},
	{Row: 1, Col: 0},
	{Row: 0, Col: -1},
	{Row: -1, Col: 0},
}

type Board [][]bool

type Game struct {
	Dims       Coord
	BlockSched map[int]Coord
	StartPos   Coord
	EndPos     Coord
	StateCache map[int]Board
}

func ReadInput(scanner *bufio.Scanner, dims Coord) (*Game, error) {
	game := Game{}
	game.BlockSched = make(map[int]Coord)
	game.StateCache = make(map[int]Board)
	game.StateCache[0] = Board{}
	lineCounter := 0
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if len(trimmed) < 1 {
			continue
		}

		lineCounter++
		fields := strings.Split(trimmed, ",")
		if len(fields) != 2 {
			return nil, fmt.Errorf("invalid input on line %d: `%s`", lineCounter, trimmed)
		}

		blockLoc := Coord{}
		var err error
		values := lo.Map(fields, func(str string, _ int) int {
			value, innerErr := strconv.Atoi(str)
			if innerErr != nil && err == nil {
				err = innerErr
			}
			return value
		})
		if err != nil {
			return nil, err
		}
		blockLoc.Row, blockLoc.Col = values[0], values[1]
		if !blockLoc.IsValid(dims) {
			return nil, fmt.Errorf("invalid block location on line %d: %v", lineCounter, blockLoc)
		}
		game.BlockSched[lineCounter] = blockLoc
	}

	game.Dims = dims

	return &game, nil
}
