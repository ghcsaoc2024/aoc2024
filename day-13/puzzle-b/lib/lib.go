package lib

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Coord struct {
	Row int64
	Col int64
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

func (c Coord) Mul(scalar int64) Coord {
	return Coord{
		Row: c.Row * scalar,
		Col: c.Col * scalar,
	}
}

func (c Coord) Div(scalar int64) Coord {
	return Coord{
		Row: c.Row / scalar,
		Col: c.Col / scalar,
	}
}

func (c Coord) Mod(scalar int64) Coord {
	return Coord{
		Row: c.Row % scalar,
		Col: c.Col % scalar,
	}
}

func (c Coord) MulOther(other Coord) Coord {
	return Coord{
		Row: c.Row * other.Row,
		Col: c.Col * other.Col,
	}
}

func (c Coord) DivOther(other Coord) Coord {
	return Coord{
		Row: c.Row / other.Row,
		Col: c.Col / other.Col,
	}
}

func (c Coord) ModOther(other Coord) Coord {
	return Coord{
		Row: c.Row % other.Row,
		Col: c.Col % other.Col,
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

func GCD(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func LCM(a, b int64) int64 {
	return a * b / GCD(a, b)
}

type Machine struct {
	ButtonA  Coord
	ButtonB  Coord
	PrizeLoc Coord
}

type Prices struct {
	A int64
	B int64
}

func ReadInput(scanner *bufio.Scanner) ([]Machine, error) {
	machines := make([]Machine, 0)
	stringBuffer := strings.Builder{}
	for scanner.Scan() {
		line := scanner.Text()
		stringBuffer.WriteString(line)
		stringBuffer.WriteString("\n")
	}

	pattern := `Button A: X\+([1-9][0-9]*), Y\+([1-9][0-9]*)
Button B: X\+([1-9][0-9]*), Y\+([1-9][0-9]*)
Prize: X=([1-9][0-9]*), Y=([1-9][0-9]*)`
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("internal error: failed to compile regex: %w", err)
	}

	allMatches := re.FindAllStringSubmatch(stringBuffer.String(), -1)
	for _, match := range allMatches {
		machine := Machine{}
		if len(match) != 7 { //nolint:mnd // Regexp-dependent, check is for internal error.
			return nil, fmt.Errorf("internal error: unexpected number of matches: %v", len(match))
		}
		var err error
		machine.ButtonA.Row, err = strconv.ParseInt(match[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("internal error: failed to convert button A row: %w (input string: `%v`)", err, match[1])
		}
		machine.ButtonA.Col, err = strconv.ParseInt(match[2], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("internal error: failed to convert button A col: %w (input string: `%v`)", err, match[2])
		}
		machine.ButtonB.Row, err = strconv.ParseInt(match[3], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("internal error: failed to convert button B row: %w (input string: `%v`)", err, match[3])
		}
		machine.ButtonB.Col, err = strconv.ParseInt(match[4], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("internal error: failed to convert button B col: %w (input string: `%v`)", err, match[4])
		}
		machine.PrizeLoc.Row, err = strconv.ParseInt(match[5], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("internal error: failed to convert prize row: %w (input string: `%v`)", err, match[5])
		}
		machine.PrizeLoc.Col, err = strconv.ParseInt(match[6], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("internal error: failed to convert prize col: %w (input string: `%v`)", err, match[6])
		}
		machines = append(machines, machine)
	}

	return machines, nil
}
