package main

import (
	"bufio"
	"container/list"
	"log"
	"math/big"
	"os"

	"main/lib"

	"github.com/alexflint/go-arg"
)

const BaseTen = 10

var MagicMultiplier = big.NewInt(2024) //nolint:gochecknoglobals,mnd // Meant as a constant.

type Args struct {
	InputFile string `arg:"positional,required" help:"input file"`
	NumSteps  int    `arg:"-n"                  default:"25"      help:"number of steps to take"`
}

func main() {
	var args Args
	arg.MustParse(&args)

	if args.NumSteps < 1 {
		log.Fatalf("number of steps must be at least 1; got %d", args.NumSteps)
	}

	theList := readInputFile(args)

	for iStep := range args.NumSteps {
		log.Printf("step %d; current length of list: %d", iStep+1, theList.Len())
		for link := theList.Front(); link != nil; link = link.Next() {
			val, ok := link.Value.(*big.Int)
			if !ok {
				log.Panicf("expected value to be of type big.Int; got %T", link.Value)
			}

			if val.Sign() == 0 {
				val.SetInt64(1)
				continue
			}

			str := val.String()
			strLen := len(str)
			if strLen%2 == 0 {
				firstHalf := str[:strLen/2]
				secondHalf := str[strLen/2:]
				val.SetString(firstHalf, BaseTen)
				newVal := big.NewInt(0)
				newVal.SetString(secondHalf, BaseTen)
				theList.InsertAfter(newVal, link)
				link = link.Next()
				continue
			}

			val.Mul(val, MagicMultiplier)
		}
	}

	log.Printf("number of values in final list: %d", theList.Len())
}

func readInputFile(args Args) *list.List {
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
	theList := lib.ReadInput(scanner)

	return theList
}
