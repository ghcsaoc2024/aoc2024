//nolint:mnd // Magic numbers all over the place
package main

import (
	"bufio"
	"fmt"
	"log"
	"main/lib"
	"math"
	"os"
	"strings"

	"github.com/alexflint/go-arg"
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

	output, err := runProgram(*computer, *program)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("output: %v", output)
}

func runProgram(computer lib.Computer, program lib.Program) (string, error) {
	stringBuffer := strings.Builder{}
	instructionPointer := 0
	for instructionPointer < len(program) {
		opCode := program[instructionPointer]
		operand := program[instructionPointer+1]
		switch opCode {
		case lib.OpADiv, lib.OpBDiv, lib.OpCDiv:
			numerator := computer.A
			cOperand, err := comboOperand(computer, operand)
			if err != nil {
				return "", err
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
				return "", err
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
				return "", err
			}
			stringBuffer.WriteString(fmt.Sprintf("%d,", cOperand%8))
		}
		instructionPointer += 2
	}

	return stringBuffer.String(), nil
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
