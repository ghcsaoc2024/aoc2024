//nolint:mnd // Magic numbers all over the place
package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"main/lib"
	"math"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/hashicorp/go-set/v3"
	"github.com/samber/lo"
)

type Args struct {
	InputFile string `arg:"positional,required" help:"input file"`
}

func main() {
	var args Args
	arg.MustParse(&args)

	computer, program, err := readInputFile(args)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("computer: %v", *computer)
	log.Printf("program: %v", *program)

	solution, err := solve(*computer, *program)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("solution: %v", solution)
}

func solve(computer lib.Computer, program lib.Program) (lib.Register, error) {
	progLen := len(program)
	output := program
	outputLen := progLen

	valuesToTest := set.New[lib.Register](1)
	valuesToTest.Insert(0)
	var previousGeneration *set.Set[lib.Register]
	for outputIdx := outputLen - 1; outputIdx >= 0; outputIdx-- {
		if valuesToTest.Empty() {
			return -1, errors.New("ran out of possible states")
		}

		previousGeneration, valuesToTest = valuesToTest, set.New[lib.Register](0)
		for value := range previousGeneration.Items() {
			for variation := range 8 {
				valueToCheck := 8*value + lib.Register(variation)
				computer.A = valueToCheck
				newOutput, err := runProgram(computer, output)
				if err != nil {
					return -1, err
				}

				if (*newOutput)[0] == output[outputIdx] {
					valuesToTest.Insert(valueToCheck)
				}
			}
		}
	}

	return lo.Min(valuesToTest.Slice()), nil
}

func runProgram(computer lib.Computer, program lib.Program) (*lib.Program, error) {
	output := lib.Program{}
	instructionPointer := 0
	for instructionPointer < len(program) {
		opCode := program[instructionPointer]
		operand := program[instructionPointer+1]
		switch opCode {
		case lib.OpADiv, lib.OpBDiv, lib.OpCDiv:
			numerator := computer.A
			cOperand, err := comboOperand(computer, operand)
			if err != nil {
				return nil, err
			}
			result := lib.Register(float64(numerator) / math.Pow(2, float64(cOperand)))

			switch opCode {
			case lib.OpADiv:
				computer.A = result
			case lib.OpBDiv:
				computer.B = result
			default:
				computer.C = result
			}
		case lib.OpXOR:
			computer.B = lib.Register(operand) ^ computer.B
		case lib.OpStore:
			cOperand, err := comboOperand(computer, operand)
			if err != nil {
				return nil, err
			}
			computer.B = cOperand % 8
		case lib.OpJumpNonZero:
			if computer.A != 0 {
				instructionPointer = int(operand)
				continue
			}
		case lib.OpRegisterXOR:
			computer.B ^= computer.C
		case lib.OpOutput:
			cOperand, err := comboOperand(computer, operand)
			if err != nil {
				return nil, err
			}
			output = append(output, lib.OpVal(cOperand%8)) //nolint:gosec // Truncation intentional
		}
		instructionPointer += 2
	}

	return &output, nil
}

func comboOperand(computer lib.Computer, operand lib.OpVal) (lib.Register, error) {
	switch operand {
	case 0, 1, 2, 3:
		return lib.Register(operand), nil
	case 4:
		return computer.A, nil
	case 5:
		return computer.B, nil
	case 6:
		return computer.C, nil
	default:
		return lib.Register(0), fmt.Errorf("invalid operand: %v", operand)
	}
}

func readInputFile(args Args) (*lib.Computer, *lib.Program, error) {
	file, err := os.Open(args.InputFile)
	if err != nil {
		log.Fatal(err) //nolint:revive // Toy code
	}
	defer func(file *os.File) {
		closeErr := file.Close()
		if closeErr != nil {
			log.Fatal(closeErr) //nolint:revive // Toy code
		}
	}(file)

	scanner := bufio.NewScanner(file)
	computer, program, err := lib.ReadInput(scanner)

	return computer, program, err //nolint:wrapcheck // Toy code
}
