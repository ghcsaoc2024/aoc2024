package lib

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Register int64

type Computer struct {
	A Register
	B Register
	C Register
}

type OpVal int8

const (
	OpADiv OpVal = iota
	OpXOR
	OpStore
	OpJumpNonZero
	OpRegisterXOR
	OpOutput
	OpBDiv
	OpCDiv
)

type Program []OpVal

func ReadInput(scanner *bufio.Scanner) (*Computer, *Program, error) {
	stringBuffer := strings.Builder{}
	for scanner.Scan() {
		line := scanner.Text()
		stringBuffer.WriteString(line)
		stringBuffer.WriteString("\n")
	}

	pattern := `Register A: (0|[1-9][0-9]*)
Register B: (0|[1-9][0-9]*)
Register C: (0|[1-9][0-9]*)
\n*
Program: ([0-7](,[0-7])*)
`
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, nil, fmt.Errorf("internal error: failed to compile regex: %w", err)
	}

	match := re.FindStringSubmatch(stringBuffer.String())
	computer := Computer{}
	program := Program{}
	if len(match) != 6 { //nolint:mnd // Regexp-dependent, check is for internal error.
		return nil, nil, fmt.Errorf("internal error: unexpected number of matches: %v", len(match))
	}
	var regVal int
	regVal, err = strconv.Atoi(match[1])
	computer.A = Register(regVal)
	if err != nil {
		return nil, nil, fmt.Errorf("internal error: failed to convert button A row: %w (input string: `%v`)", err, match[1])
	}
	regVal, err = strconv.Atoi(match[2])
	computer.B = Register(regVal)
	if err != nil {
		return nil, nil, fmt.Errorf("internal error: failed to convert button A col: %w (input string: `%v`)", err, match[2])
	}
	regVal, err = strconv.Atoi(match[3])
	computer.C = Register(regVal)
	if err != nil {
		return nil, nil, fmt.Errorf("internal error: failed to convert button B row: %w (input string: `%v`)", err, match[3])
	}
	programString := match[4]
	opCodeStrings := strings.Split(programString, ",")
	for _, opCodeString := range opCodeStrings {
		opCode, err := strconv.Atoi(opCodeString)
		if err != nil {
			return nil, nil, fmt.Errorf("internal error: failed to convert op code: %w (input string: `%v`)", err, opCodeString)
		}
		if opCode < 0 || opCode > 7 {
			return nil, nil, fmt.Errorf("internal error: invalid op code: %v (input string: `%v`)", opCode, opCodeString)
		}
		program = append(program, OpVal(opCode))
	}

	return &computer, &program, nil
}
