package lib

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
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

func (c Coord) Mul(scalar int) Coord {
	return Coord{
		Row: c.Row * scalar,
		Col: c.Col * scalar,
	}
}

type Machine struct {
	ButtonA  Coord
	ButtonB  Coord
	PrizeLoc Coord
}

type ScalarMachinePair struct {
	A int
	B int
}

type Solution ScalarMachinePair

type Prices ScalarMachinePair

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
		machine.ButtonA.Row, err = strconv.Atoi(match[1])
		if err != nil {
			return nil, fmt.Errorf("internal error: failed to convert button A row: %w (input string: `%v`)", err, match[1])
		}
		machine.ButtonA.Col, err = strconv.Atoi(match[2])
		if err != nil {
			return nil, fmt.Errorf("internal error: failed to convert button A col: %w (input string: `%v`)", err, match[2])
		}
		machine.ButtonB.Row, err = strconv.Atoi(match[3])
		if err != nil {
			return nil, fmt.Errorf("internal error: failed to convert button B row: %w (input string: `%v`)", err, match[3])
		}
		machine.ButtonB.Col, err = strconv.Atoi(match[4])
		if err != nil {
			return nil, fmt.Errorf("internal error: failed to convert button B col: %w (input string: `%v`)", err, match[4])
		}
		machine.PrizeLoc.Row, err = strconv.Atoi(match[5])
		if err != nil {
			return nil, fmt.Errorf("internal error: failed to convert prize row: %w (input string: `%v`)", err, match[5])
		}
		machine.PrizeLoc.Col, err = strconv.Atoi(match[6])
		if err != nil {
			return nil, fmt.Errorf("internal error: failed to convert prize col: %w (input string: `%v`)", err, match[6])
		}
		machines = append(machines, machine)
	}

	return machines, nil
}
