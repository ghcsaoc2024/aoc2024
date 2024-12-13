package lib

import (
	"bufio"
	"fmt"
	"math/big"
	"regexp"
	"strings"
)

const BaseTen = 10

type Coord struct {
	Row *big.Int
	Col *big.Int
}

func (c Coord) Add(other Coord) Coord {
	return Coord{
		Row: big.NewInt(0).Add(c.Row, other.Row),
		Col: big.NewInt(0).Add(c.Col, other.Col),
	}
}

func (c Coord) Sub(other Coord) Coord {
	return Coord{
		Row: big.NewInt(0).Sub(c.Row, other.Row),
		Col: big.NewInt(0).Sub(c.Col, other.Col),
	}
}

func (c Coord) Mul(scalar *big.Int) Coord {
	return Coord{
		Row: big.NewInt(0).Mul(c.Row, scalar),
		Col: big.NewInt(0).Mul(c.Col, scalar),
	}
}

func GCD(a, b *big.Int) *big.Int {
	for b.Sign() != 0 {
		a, b = b, big.NewInt(0).Mod(a, b)
	}
	return big.NewInt(0).Set(a)
}

func LCM(a, b *big.Int) *big.Int {
	gcd := GCD(a, b)
	result := big.NewInt(0).Mul(a, b)
	result = result.Quo(result, gcd)

	return result
}

type Machine struct {
	ButtonA  Coord
	ButtonB  Coord
	PrizeLoc Coord
}

type Prices struct {
	A *big.Int
	B *big.Int
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
		var ok bool
		var value *big.Int
		value = big.NewInt(0)
		value, ok = value.SetString(match[1], BaseTen)
		if !ok {
			return nil, fmt.Errorf("internal error: failed to convert button A row: %w (input string: `%v`)", err, match[1])
		}
		machine.ButtonA.Row = value
		value = big.NewInt(0)
		value, ok = value.SetString(match[2], BaseTen)
		if !ok {
			return nil, fmt.Errorf("internal error: failed to convert button A col: %w (input string: `%v`)", err, match[2])
		}
		machine.ButtonA.Col = value
		value = big.NewInt(0)
		value, ok = value.SetString(match[3], BaseTen)
		if !ok {
			return nil, fmt.Errorf("internal error: failed to convert button B row: %w (input string: `%v`)", err, match[3])
		}
		machine.ButtonB.Row = value
		value = big.NewInt(0)
		value, ok = value.SetString(match[4], BaseTen)
		if !ok {
			return nil, fmt.Errorf("internal error: failed to convert button B col: %w (input string: `%v`)", err, match[4])
		}
		machine.ButtonB.Col = value
		value = big.NewInt(0)
		value, ok = value.SetString(match[5], BaseTen)
		if !ok {
			return nil, fmt.Errorf("internal error: failed to convert prize row: %w (input string: `%v`)", err, match[5])
		}
		machine.PrizeLoc.Row = value
		value = big.NewInt(0)
		value, ok = value.SetString(match[6], BaseTen)
		if !ok {
			return nil, fmt.Errorf("internal error: failed to convert prize col: %w (input string: `%v`)", err, match[6])
		}
		machine.PrizeLoc.Col = value
		machines = append(machines, machine)
	}

	return machines, nil
}
