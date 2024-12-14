package lib

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Coord struct {
	X int64
	Y int64
}

func (c Coord) Add(other Coord) Coord {
	return Coord{
		X: c.X + other.X,
		Y: c.Y + other.Y,
	}
}

func (c Coord) Mul(scalar int64) Coord {
	return Coord{
		X: c.X * scalar,
		Y: c.Y * scalar,
	}
}

func (c Coord) Div(scalar int64) Coord {
	return Coord{
		X: c.X / scalar,
		Y: c.Y / scalar,
	}
}

func (c Coord) DivOther(other Coord) Coord {
	return Coord{
		X: c.X / other.X,
		Y: c.Y / other.Y,
	}
}

func (c Coord) ModOther(other Coord) Coord {
	return Coord{
		X: c.X % other.X,
		Y: c.Y % other.Y,
	}
}

type Robot struct {
	Pos Coord
	Vel Coord
}

func ReadInput(scanner *bufio.Scanner) ([]Robot, error) {
	robots := make([]Robot, 0)
	stringBuffer := strings.Builder{}
	for scanner.Scan() {
		line := scanner.Text()
		stringBuffer.WriteString(line)
		stringBuffer.WriteString("\n")
	}

	pattern := `p=([0-9]*),([0-9]*)\s+v=(-?[0-9]*),(-?[0-9]*)`
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("internal error: failed to compile regex: %w", err)
	}

	allMatches := re.FindAllStringSubmatch(stringBuffer.String(), -1)
	for _, match := range allMatches {
		robot := Robot{}
		if len(match) != 5 { //nolint:mnd // Regexp-dependent, check is for internal error.
			return nil, fmt.Errorf("internal error: unexpected number of matches: %v", len(match))
		}
		var err error
		robot.Pos.X, err = strconv.ParseInt(match[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("internal error: failed to convert button A row: %w (input string: `%v`)", err, match[1])
		}
		robot.Pos.Y, err = strconv.ParseInt(match[2], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("internal error: failed to convert button A col: %w (input string: `%v`)", err, match[2])
		}
		robot.Vel.X, err = strconv.ParseInt(match[3], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("internal error: failed to convert button B row: %w (input string: `%v`)", err, match[3])
		}
		robot.Vel.Y, err = strconv.ParseInt(match[4], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("internal error: failed to convert button B col: %w (input string: `%v`)", err, match[4])
		}
		robots = append(robots, robot)
	}

	return robots, nil
}
