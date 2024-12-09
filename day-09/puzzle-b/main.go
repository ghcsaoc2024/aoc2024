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

	occupiedPred := func(val int) bool {
		return val != lib.FreeSpaceIndicator
	}
	nonValuePred := func(val int) func(int) bool {
		return func(v int) bool {
			return v != val
		}
	}
	freeSpacePtr := lo.IndexOf(diskContents, lib.FreeSpaceIndicator)
	_, lastFilledPtr, _ := lo.FindLastIndexOf(diskContents, occupiedPred)
	for freeSpacePtr != -1 && lastFilledPtr != -1 {
		currFileID := diskContents[lastFilledPtr]
		_, fileStart, _ := lo.FindLastIndexOf(diskContents[:lastFilledPtr], nonValuePred(currFileID))
		fileStart++
		fileEnd := lastFilledPtr + 1
		fileLength := fileEnd - fileStart
		if freeSpacePtr >= lastFilledPtr {
			freeSpacePtr = lo.IndexOf(diskContents, lib.FreeSpaceIndicator)
			if freeSpacePtr == -1 {
				break
			}
			_, lastFilledPtr, _ = lo.FindLastIndexOf(diskContents[:fileStart], occupiedPred)
			continue
		}

		_, nextFilled, _ := lo.FindIndexOf(diskContents[freeSpacePtr:], occupiedPred)
		if nextFilled == -1 {
			nextFilled = len(diskContents)
		} else {
			nextFilled += freeSpacePtr
		}

		freeSpaceLength := nextFilled - freeSpacePtr

		if freeSpaceLength >= fileLength {
			emptyRun := lib.RunOf(fileLength, lib.FreeSpaceIndicator)
			copy(diskContents[freeSpacePtr:], diskContents[fileStart:fileEnd])
			copy(diskContents[fileStart:], emptyRun)
			freeSpacePtr = lo.IndexOf(diskContents, lib.FreeSpaceIndicator)
			_, lastFilledPtr, _ = lo.FindLastIndexOf(diskContents[:fileStart], occupiedPred)
		} else {
			freeSpacePtr = nextFilled + lo.IndexOf(diskContents[nextFilled:], lib.FreeSpaceIndicator)
		}
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
