package main

import (
	"bufio"
	"log"
	"os"

	"main/lib"

	"github.com/alexflint/go-arg"
)

type Args struct {
	InputFile   string `arg:"positional,required" help:"input file"`
	NumMaxSteps int    `arg:"-n, --max-steps"     default:"100"     help:"maximum number of steps"`
}

func main() {
	var args Args
	arg.MustParse(&args)

	machines, err := readInputFile(args)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("%d machines read", len(machines))

	totalCost := 0
	prices := lib.Prices{A: 3, B: 1}
	for iMachine, machine := range machines {
		var bestPrice int
		alreadySolved := false
		for numA := range args.NumMaxSteps {
			numB := (machine.PrizeLoc.Row - numA*machine.ButtonA.Row) / machine.ButtonB.Row
			if numB >= args.NumMaxSteps {
				continue
			}

			price := numA*prices.A + numB*prices.B
			result := machine.ButtonA.Mul(numA)
			result = result.Add(machine.ButtonB.Mul(numB))
			if result == machine.PrizeLoc {
				if !alreadySolved {
					alreadySolved = true
					bestPrice = price
					continue
				}

				bestPrice = min(bestPrice, price)
			}
		}
		if !alreadySolved {
			log.Printf("machine %d: not solvable", iMachine)
			continue
		}

		totalCost += bestPrice
		log.Printf("machine %d: solvable; best price: %d", iMachine, bestPrice)
	}

	log.Printf("total cost: %d", totalCost)
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
