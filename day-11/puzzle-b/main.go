package main

import (
	"bufio"
	"log"
	"math/big"
	"os"

	"main/lib"

	"github.com/alexflint/go-arg"
)

const BaseTen = 10

var MagicMultiplier = big.NewInt(2024) //nolint:gochecknoglobals,mnd // Meant as a constant.

type LaunchPoint struct {
	Str            string
	RemainingDepth int
}

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

	values := readInputFile(args)
	total := 0
	expansionCache := make(map[LaunchPoint]int)
	for _, value := range values {
		total += expandWithCache(value, args.NumSteps, &expansionCache)
	}

	log.Printf("total: %d", total)
}

func readInputFile(args Args) []*big.Int {
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
	values := lib.ReadInput(scanner)

	return values
}

func expandWithCache(value *big.Int, numSteps int, expansionCache *map[LaunchPoint]int) int {
	str := value.String()
	if total, ok := (*expansionCache)[LaunchPoint{Str: str, RemainingDepth: numSteps}]; ok {
		return total
	}

	total := expand(value, numSteps, expansionCache)
	(*expansionCache)[LaunchPoint{Str: str, RemainingDepth: numSteps}] = total

	return total
}

func expand(value *big.Int, numSteps int, expansionCache *map[LaunchPoint]int) int {
	str := value.String()
	if total, ok := (*expansionCache)[LaunchPoint{Str: str, RemainingDepth: numSteps}]; ok {
		return total
	}

	if numSteps < 1 {
		return 1
	}

	if value.Sign() == 0 {
		value.SetInt64(1)
		return expandWithCache(value, numSteps-1, expansionCache)
	}

	strLen := len(str)
	if strLen%2 == 0 {
		firstHalf := str[:strLen/2]
		secondHalf := str[strLen/2:]
		value.SetString(firstHalf, BaseTen)
		firstTotal := expandWithCache(value, numSteps-1, expansionCache)
		value.SetString(secondHalf, BaseTen)
		return firstTotal + expandWithCache(value, numSteps-1, expansionCache)
	}

	value.Mul(value, MagicMultiplier)

	return expandWithCache(value, numSteps-1, expansionCache)
}
