package main

import (
	"bufio"
	"log"
	"math/big"
	"os"

	"main/lib"

	"github.com/alexflint/go-arg"
	"github.com/samber/lo"
)

func main() {
	var args struct {
		InputFile string `arg:"positional,required" help:"input file"`
	}
	arg.MustParse(&args)

	file, err := os.Open(args.InputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		closeErr := file.Close()
		if closeErr != nil {
			log.Fatal(closeErr)
		}
	}(file)

	scanner := bufio.NewScanner(file)

	diskContents := lib.ReadInput(scanner)
	log.Printf("length of disk contents: %d", len(diskContents))

	freeSpacePtr := lo.IndexOf(diskContents, lib.FreeSpaceIndicator)
	usedPredicate := func(val int) bool {
		return val != lib.FreeSpaceIndicator
	}
	_, lastFilledPtr, _ := lo.FindLastIndexOf(diskContents, usedPredicate)
	for freeSpacePtr != -1 && lastFilledPtr != -1 {
		if freeSpacePtr >= lastFilledPtr {
			break
		}

		diskContents[freeSpacePtr], diskContents[lastFilledPtr] = diskContents[lastFilledPtr], diskContents[freeSpacePtr]
		freeSpacePtr += lo.IndexOf(diskContents[freeSpacePtr:], lib.FreeSpaceIndicator)
		_, lastFilledPtr, _ = lo.FindLastIndexOf(diskContents[:lastFilledPtr], usedPredicate)
	}

	checkSum := big.NewInt(0)
	for i, val := range diskContents {
		if val == lib.FreeSpaceIndicator {
			continue
		}
		checkSum.Add(checkSum, big.NewInt(int64(i*val)))
	}

	log.Printf("checksum: %s", checkSum.String())
}
