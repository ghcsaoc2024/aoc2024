package main

import (
	"bufio"
	"log"
	"main/lib"
	"math/big"
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

	totalPrice := big.NewInt(0)
	bumpSum, ok := big.NewInt(0).SetString("10000000000000", lib.BaseTen)
	if !ok {
		log.Panic("internal error: could not set big.Int to bumpSum string")
	}
	bump := lib.Coord{
		Row: big.NewInt(0).Set(bumpSum),
		Col: big.NewInt(0).Set(bumpSum),
	}
	for iMachine, machine := range machines {
		machine.PrizeLoc = machine.PrizeLoc.Add(bump)

		lcmA := lib.LCM(machine.ButtonA.Row, machine.ButtonA.Col)
		rowMulA := big.NewInt(0).Quo(lcmA, machine.ButtonA.Row)
		colMulA := big.NewInt(0).Quo(lcmA, machine.ButtonA.Col)

		newMulRowB := big.NewInt(0).Mul(rowMulA, machine.ButtonB.Row)
		newMulColB := big.NewInt(0).Mul(colMulA, machine.ButtonB.Col)

		newPrizeLoc := lib.Coord{
			Row: big.NewInt(0).Mul(machine.PrizeLoc.Row, rowMulA),
			Col: big.NewInt(0).Mul(machine.PrizeLoc.Col, colMulA),
		}
		mulDiffB := big.NewInt(0).Sub(newMulColB, newMulRowB)
		prizeLocDiffB := big.NewInt(0).Sub(newPrizeLoc.Col, newPrizeLoc.Row)
		var priceA *big.Int

		if big.NewInt(0).Mod(prizeLocDiffB, mulDiffB).Sign() == 0 {
			numB := big.NewInt(0).Quo(prizeLocDiffB, mulDiffB)
			locB := machine.ButtonB.Mul(numB)
			diff := machine.PrizeLoc.Sub(locB)
			if big.NewInt(0).Quo(diff.Row, machine.ButtonA.Row).Cmp(big.NewInt(0).Quo(diff.Col, machine.ButtonA.Col)) != 0 {
				log.Panic("internal error: the math isn't mathing (the quotients are unequal)")
			}
			numA := big.NewInt(0).Quo(diff.Row, machine.ButtonA.Row)
			priceA = big.NewInt(0).Add(big.NewInt(0).Mul(numA, big.NewInt(3)), numB)
		}

		lcmB := lib.LCM(machine.ButtonB.Row, machine.ButtonB.Col)
		rowMulB := big.NewInt(0).Quo(lcmB, machine.ButtonB.Row)
		colMulB := big.NewInt(0).Quo(lcmB, machine.ButtonB.Col)

		newMulRowA := big.NewInt(0).Mul(rowMulB, machine.ButtonA.Row)
		newMulColA := big.NewInt(0).Mul(colMulB, machine.ButtonA.Col)

		newPrizeLoc = lib.Coord{
			Row: big.NewInt(0).Mul(machine.PrizeLoc.Row, rowMulB),
			Col: big.NewInt(0).Mul(machine.PrizeLoc.Col, colMulB),
		}
		mulDiffA := big.NewInt(0).Sub(newMulColA, newMulRowA)
		prizeLocDiffA := big.NewInt(0).Sub(newPrizeLoc.Col, newPrizeLoc.Row)
		var priceB *big.Int

		if big.NewInt(0).Mod(prizeLocDiffA, mulDiffA).Sign() == 0 {
			numA := big.NewInt(0).Quo(prizeLocDiffA, mulDiffA)
			locA := machine.ButtonA.Mul(numA)
			diff := machine.PrizeLoc.Sub(locA)
			if big.NewInt(0).Quo(diff.Row, machine.ButtonB.Row).Cmp(big.NewInt(0).Quo(diff.Col, machine.ButtonB.Col)) != 0 {
				log.Panic("internal error: the math isn't mathing (the quotients are unequal)")
			}
			numB := big.NewInt(0).Quo(diff.Row, machine.ButtonB.Row)
			priceB = big.NewInt(0).Add(big.NewInt(0).Mul(numA, big.NewInt(3)), numB)
		}

		if (priceA != priceB) && (priceA == nil || priceB == nil || priceA.Cmp(priceB) != 0) {
			log.Panicf("internal error: the math isn't mathing (priceA: %d, priceB: %d)", priceA, priceB)
		}

		if priceA == nil {
			log.Printf("machine %d: not solvable", iMachine)
			continue
		}
		log.Printf("machine %d: solvable; price: %d", iMachine, priceA)
		totalPrice = totalPrice.Add(totalPrice, priceA)
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
