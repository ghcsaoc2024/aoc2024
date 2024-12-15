package lib

import (
	"bufio"
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

func (c Coord) IsValid(dimensions Coord) bool {
	if c.Row < 0 || c.Row >= dimensions.Row {
		return false
	}
	if c.Col < 0 || c.Col >= dimensions.Col {
		return false
	}
	return true
}

var DirectionsByRune = map[rune]Coord{ //nolint:gochecknoglobals // Meant as a constant
	'<': {Row: 0, Col: -1},
	'>': {Row: 0, Col: 1},
	'^': {Row: -1, Col: 0},
	'v': {Row: 1, Col: 0},
}

type Cell int

const (
	Empty Cell = iota
	BoxL
	BoxR
	Wall
)

type Game struct {
	Board        [][]Cell
	Moves        []Coord
	Boxes        []Coord
	BoxesByCoord map[Coord]int
	Robot        Coord
}

func ReadInput(scanner *bufio.Scanner) (*Game, error) {
	game := &Game{}
	game.BoxesByCoord = make(map[Coord]int)
	for scanner.Scan() {
		line := scanner.Text()
		if len(strings.TrimSpace(line)) < 1 {
			break
		}

		row := make([]Cell, len(line)*2)
		for iCol, char := range line {
			realICol := iCol * 2
			coord := Coord{Row: len(game.Board), Col: realICol}
			switch char {
			case '@':
				game.Robot = coord
				fallthrough
			case '.':
				row[realICol] = Empty
				row[realICol+1] = Empty
			case '#':
				row[realICol] = Wall
				row[realICol+1] = Wall
			case 'O':
				row[realICol] = BoxL
				row[realICol+1] = BoxR
				game.BoxesByCoord[coord] = len(game.Boxes)
				game.Boxes = append(game.Boxes, coord)

			default:
				return nil, fmt.Errorf("unrecognized cell character: `%c`", char)
			}
		}

		game.Board = append(game.Board, row)
	}

	for scanner.Scan() {
		line := scanner.Text()
		for _, char := range line {
			dir, ok := DirectionsByRune[char]
			if !ok {
				return nil, fmt.Errorf("unrecognized direction character: `%c`", char)
			}
			game.Moves = append(game.Moves, dir)
		}
	}

	return game, nil
}
