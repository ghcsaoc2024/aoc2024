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
			// The currently-pointed-to file does not fit in any free space to its left;
			// reset the free space pointer to the leftmost free space, and...
			freeSpacePtr = lo.IndexOf(diskContents, lib.FreeSpaceIndicator)
			// ...advance the last-filled pointer to the last filled space prior to the start
			// of the currently-pointed-to file.
			_, lastFilledPtr, _ = lo.FindLastIndexOf(diskContents[:fileStart], occupiedPred)
			continue
		}

		// Calculate how much free space we're looking at
		_, nextFilled, _ := lo.FindIndexOf(diskContents[freeSpacePtr:], occupiedPred)
		if nextFilled == -1 {
			nextFilled = len(diskContents)
		} else {
			nextFilled += freeSpacePtr
		}
		freeSpaceLength := nextFilled - freeSpacePtr

		// Does the currently-pointed-to file fit in the free space?
		if freeSpaceLength >= fileLength {
			// It does! Move it over, overwriting its original position with empties.
			emptyRun := lib.RunOf(fileLength, lib.FreeSpaceIndicator)
			copy(diskContents[freeSpacePtr:], diskContents[fileStart:fileEnd])
			copy(diskContents[fileStart:], emptyRun)
			// Reset the free space pointer to the leftmost free space, and...
			freeSpacePtr = lo.IndexOf(diskContents, lib.FreeSpaceIndicator)
			// ...advance the last-filled pointer to the last filled space prior to the start
			// of the currently-pointed-to file. (Remember: we don't get to try a single file
			// more than once; so whatever lastFilledPtr has moved across, shall never be
			// revisited again!)
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
