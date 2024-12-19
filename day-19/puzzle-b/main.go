package main

import (
	"bufio"
	"log"
	"main/lib"
	"os"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/samber/lo"
)

type Args struct {
	InputFile string `arg:"positional,required" help:"input file"`
}

func main() {
	var args Args
	arg.MustParse(&args)

	inventory, patterns, err := readInputFile(args)
	if err != nil {
		log.Panic(err)
	}

	newInventory := pruneInventory(inventory)

	log.Printf("pruned inventory from %d to %d", len(inventory), len(newInventory))
	pruneventory, _ := lo.Difference(inventory, newInventory)
	inventory = newInventory

	solutionsCache := make(map[string]int64)
	nTotalSolutions := int64(0)
	for iPattern, pattern := range patterns {
		nSolutions := solve(inventory, pruneventory, solutionsCache, pattern)
		log.Printf("pattern %d: %d solutions", iPattern, nSolutions)
		nTotalSolutions += nSolutions
	}

	log.Printf("number of total solutions: %d", nTotalSolutions)
}

func pruneInventory(inventory []string) []string {
	remainingInventory := make([]string, len(inventory)-1)
	for iTowel, towel := range inventory {
		copy(remainingInventory[:iTowel], inventory[:iTowel])
		copy(remainingInventory[iTowel:], inventory[iTowel+1:])
		nSolutions := solve(remainingInventory, nil, nil, towel)
		if nSolutions > 0 {
			log.Printf("pruning %s", towel)
			return pruneInventory(remainingInventory)
		}
	}

	return inventory
}

func solve(inventory, pruneventory []string, solutionsCache map[string]int64, pattern string) int64 {
	if len(pattern) < 1 {
		return 1
	}

	if solutionsCache != nil {
		if nTotalSolutions, ok := solutionsCache[pattern]; ok {
			return nTotalSolutions
		}
	}

	nTotalSolutions := int64(0)
	for _, towel := range inventory {
		if !strings.HasPrefix(pattern, towel) {
			continue
		}

		remainder := pattern[len(towel):]
		nRemainderSolutions := solve(inventory, pruneventory, solutionsCache, remainder)

		if nRemainderSolutions >= 0 {
			if solutionsCache != nil {
				solutionsCache[remainder] = nRemainderSolutions
			}
			nTotalSolutions += nRemainderSolutions
		}
	}

	if pruneventory == nil {
		return nTotalSolutions
	}

	for _, towel := range pruneventory {
		if !strings.HasPrefix(pattern, towel) {
			continue
		}

		remainder := pattern[len(towel):]
		nRemainderSolutions, ok := solutionsCache[remainder]
		if !ok {
			log.Panic("internal error: no solution for pruned towel") //nolint:revive // Toy code
		}

		nTotalSolutions += nRemainderSolutions
	}

	return nTotalSolutions
}

func readInputFile(args Args) ([]string, []string, error) {
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
	inventory, patterns, err := lib.ReadInput(scanner)

	return inventory, patterns, err //nolint:wrapcheck // Toy code
}
