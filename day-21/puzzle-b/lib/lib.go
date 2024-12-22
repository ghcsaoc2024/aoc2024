package lib

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"

	"github.com/samber/lo"
)

const (
	InvalidNumPadKey int = iota
	NumPadKey1
	NumPadKey2
	NumPadKey3
	NumPadKey4
	NumPadKey5
	NumPadKey6
	NumPadKey7
	NumPadKey8
	NumPadKey9
	NumPadKey0
	NumPadKeyA
)

var NumPadKeyByRune = map[rune]int{ //nolint:gochecknoglobals // Meant as a constant
	'1': NumPadKey1,
	'2': NumPadKey2,
	'3': NumPadKey3,
	'4': NumPadKey4,
	'5': NumPadKey5,
	'6': NumPadKey6,
	'7': NumPadKey7,
	'8': NumPadKey8,
	'9': NumPadKey9,
	'A': NumPadKeyA,
	'0': NumPadKey0,
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

const (
	InvalidAction int = iota
	MoveUp
	MoveDown
	MoveLeft
	MoveRight
	Press
)

var Actions = map[int]Coord{ //nolint:gochecknoglobals // Meant as a constant
	MoveUp:    {Row: -1, Col: 0},
	MoveDown:  {Row: 1, Col: 0},
	MoveLeft:  {Row: 0, Col: -1},
	MoveRight: {Row: 0, Col: 1},
	Press:     {Row: 0, Col: 0},
}

var NumPadLayout = [][]int{ //nolint:gochecknoglobals // Meant as a constant
	{NumPadKey7, NumPadKey8, NumPadKey9},
	{NumPadKey4, NumPadKey5, NumPadKey6},
	{NumPadKey1, NumPadKey2, NumPadKey3},
	{InvalidNumPadKey, NumPadKey0, NumPadKeyA},
}

var ActionByRune = map[rune]int{ //nolint:gochecknoglobals // Meant as a constant
	'^': MoveUp,
	'v': MoveDown,
	'<': MoveLeft,
	'>': MoveRight,
	'A': Press,
}

var ActionPadLayout = [][]int{ //nolint:gochecknoglobals // Meant as a constant
	{InvalidAction, MoveUp, Press},
	{MoveLeft, MoveDown, MoveRight},
}

type NumPadCode []int

func ReadInput(scanner *bufio.Scanner) ([]NumPadCode, error) {
	pattern := `[0-9]+A`
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("internal error: failed to compile regex: %w", err)
	}

	var numPadCodes []NumPadCode
	iLine := 0
	for scanner.Scan() {
		line := scanner.Text()
		iLine++
		trimmedLine := strings.TrimSpace(line)
		match := re.FindString(trimmedLine)
		if len(match) == 0 {
			continue
		}

		numPadCode := lo.Map([]rune(match), func(item rune, _ int) int {
			return lo.ValueOr(NumPadKeyByRune, item, InvalidNumPadKey)
		})

		if lo.Contains(numPadCode, InvalidNumPadKey) {
			return nil, fmt.Errorf("invalid num pad code: %v (line %d)", numPadCode, iLine)
		}

		numPadCodes = append(numPadCodes, numPadCode)
	}

	return numPadCodes, nil
}
