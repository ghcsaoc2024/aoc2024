package main

import (
	"bufio"
	"log"
	"main/lib"
	"os"

	"github.com/alexflint/go-arg"
)

type Args struct {
	InputFile string `arg:"positional,required" help:"input file"`
}

func main() {
	var args Args
	arg.MustParse(&args)

	machines, err := readInputFile(args)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("%d machines read", len(machines))

	totalPrice := int64(0)
	bump := lib.Coord{Row: 10000000000000, Col: 10000000000000} //nolint:mnd // From puzzle
	for iMachine, machine := range machines {
		machine.PrizeLoc = machine.PrizeLoc.Add(bump)

		lcmA := lib.LCM(machine.ButtonA.Row, machine.ButtonA.Col)
		rowMulA := lcmA / machine.ButtonA.Row
		colMulA := lcmA / machine.ButtonA.Col

		newMulRowB := rowMulA * machine.ButtonB.Row
		newMulColB := colMulA * machine.ButtonB.Col
		newPrizeLoc := lib.Coord{Row: machine.PrizeLoc.Row * rowMulA, Col: machine.PrizeLoc.Col * colMulA}
		mulDiffB := newMulColB - newMulRowB
		prizeLocDiffB := newPrizeLoc.Col - newPrizeLoc.Row
		priceA := int64(-1)
		if prizeLocDiffB%mulDiffB == 0 {
			numB := prizeLocDiffB / mulDiffB
			locB := machine.ButtonB.Mul(numB)
			diff := machine.PrizeLoc.Sub(locB)
			if diff.Row/machine.ButtonA.Row != diff.Col/machine.ButtonA.Col {
				log.Panic("internal error: the math isn't mathing (the quotients are unequal)")
			}
			numA := diff.Row / machine.ButtonA.Row
			priceA = numA*3 + numB*1
		}

		lcmB := lib.LCM(machine.ButtonB.Row, machine.ButtonB.Col)
		rowMulB := lcmB / machine.ButtonB.Row
		colMulB := lcmB / machine.ButtonB.Col

		newMulRowA := rowMulB * machine.ButtonA.Row
		newMulColA := colMulB * machine.ButtonA.Col
		newPrizeLoc = lib.Coord{Row: machine.PrizeLoc.Row * rowMulB, Col: machine.PrizeLoc.Col * colMulB}
		mulDiffA := newMulColA - newMulRowA
		prizeLocDiffA := newPrizeLoc.Col - newPrizeLoc.Row
		priceB := int64(-1)
		if prizeLocDiffA%mulDiffA == 0 {
			numA := prizeLocDiffA / mulDiffA
			locA := machine.ButtonA.Mul(numA)
			diff := machine.PrizeLoc.Sub(locA)
			if diff.Row/machine.ButtonB.Row != diff.Col/machine.ButtonB.Col {
				log.Panic("internal error: the math isn't mathing (the quotients are unequal)")
			}
			numB := diff.Row / machine.ButtonB.Row
			priceB = numA*3 + numB*1
		}

		if priceA != priceB {
			log.Panicf("internal error: the math isn't mathing (priceA: %d, priceB: %d)", priceA, priceB)
		}

		var price int64
		switch {
		case priceA < 0 && priceB < 0:

			log.Printf("machine %d: not solvable", iMachine)
			continue
		case priceA < 0:
			price = priceB
		case priceB < 0:
			price = priceA
		default:
			price = min(priceA, priceB)
		}
		log.Printf("machine %d: solvable; price: %d", iMachine, price)
		totalPrice += price
	}

	log.Printf("total price: %d", totalPrice)
}

func readInputFile(args Args) ([]lib.Machine, error) {
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
	values, err := lib.ReadInput(scanner)

	return values, err //nolint:wrapcheck // Toy code
}
