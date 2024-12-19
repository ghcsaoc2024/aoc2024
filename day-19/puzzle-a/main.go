package main

import (
	"bufio"
	"log"
	"main/lib"
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

	inventory, patterns, err := readInputFile(args)
	if err != nil {
		log.Panic(err)
	}

	newInventory, err := pruneInventory(inventory)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("pruned inventory from %d to %d", len(inventory), len(newInventory))
	inventory = newInventory

	solutionsCache := make(map[string]bool)
	nSolvable := 0
	for iPattern, pattern := range patterns {
		solvable, err := solve(inventory, solutionsCache, pattern)
		if err != nil {
			log.Panic(err)
		}

		if solvable {
			nSolvable++
			log.Printf("pattern %d: solvable", iPattern)
		} else {
			log.Printf("pattern %d: not solvable", iPattern)
		}
	}

	log.Printf("number of solvable patterns: %d", nSolvable)
}

func pruneInventory(inventory []string) ([]string, error) {
	remainingInventory := make([]string, len(inventory)-1)
	for iTowel, towel := range inventory {
		copy(remainingInventory[:iTowel], inventory[:iTowel])
		copy(remainingInventory[iTowel:], inventory[iTowel+1:])
		solvable, err := solve(remainingInventory, nil, towel)
		if err != nil {
			return nil, err
		}

		if solvable {
			log.Printf("pruning %s", towel)
			return pruneInventory(remainingInventory)
		}
	}

	return inventory, nil
}

func solve(inventory []string, solutionsCache map[string]bool, pattern string) (bool, error) {
	if len(pattern) < 1 {
		return true, nil
	}

	if solutionsCache != nil {
		if solvable, ok := solutionsCache[pattern]; ok {
			return solvable, nil
		}
	}

	for _, towel := range inventory {
		if !strings.HasPrefix(pattern, towel) {
			continue
		}

		remainder := pattern[len(towel):]
		remainderSolvable, err := solve(inventory, solutionsCache, remainder)
		if err != nil {
			return false, err
		}

		if remainderSolvable {
			if solutionsCache != nil {
				solutionsCache[remainder] = true
			}

			return true, nil
		}
	}

	if solutionsCache != nil {
		solutionsCache[pattern] = false
	}

	return false, nil
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
