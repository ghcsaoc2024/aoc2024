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

type FileInfo struct {
	ID     int
	Length int
	Start  int
	End    int
}

func occupiedPred(val int) bool {
	return val != lib.FreeSpaceIndicator
}

func generateAntiPred(val int) func(int) bool {
	return func(v int) bool {
		return v != val
	}
}

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
	log.Printf("length of disk: %d", len(diskContents))

	freeSpacePtr := resetFreeSpacePtr(diskContents)
	_, lastFilledPtr, _ := lo.FindLastIndexOf(diskContents, occupiedPred)
	for freeSpacePtr != -1 && lastFilledPtr != -1 {
		fileInfo := getFileInfo(diskContents, lastFilledPtr)

		if freeSpacePtr >= lastFilledPtr {
			// The currently-pointed-to file does not fit in any free space to its left;
			// reset the free space pointer to the leftmost free space, and...
			freeSpacePtr = resetFreeSpacePtr(diskContents)
			// ...advance the last-filled pointer to the last filled space prior to the start
			// of the currently-pointed-to file.
			lastFilledPtr = advanceLastFilledPtr(diskContents, fileInfo)
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
		if freeSpaceLength >= fileInfo.Length {
			// It does! Move it over, overwriting its original position with empties.
			moveFile(&diskContents, fileInfo, freeSpacePtr)
			// Reset the free space pointer to the leftmost free space, and...
			freeSpacePtr = resetFreeSpacePtr(diskContents)
			// ...advance the last-filled pointer to the last filled space prior to the start
			// of the currently-pointed-to file. (Remember: we don't get to try a single file
			// more than once; so whatever lastFilledPtr has moved across, shall never be
			// revisited again!)
			lastFilledPtr = advanceLastFilledPtr(diskContents, fileInfo)
		} else {
			// It doesn't fit; move the free space pointer to the next free space, and try again.
			freeSpacePtr = nextFilled + lo.IndexOf(diskContents[nextFilled:], lib.FreeSpaceIndicator)
		}
	}

	checkSum := calcChecksum(diskContents)

	log.Printf("checksum: %s", checkSum.String())
}

func resetFreeSpacePtr(diskContents []int) int {
	return lo.IndexOf(diskContents, lib.FreeSpaceIndicator)
}

func advanceLastFilledPtr(diskContents []int, fileInfo FileInfo) int {
	_, lastFilledPtr, _ := lo.FindLastIndexOf(diskContents[:fileInfo.Start], occupiedPred)

	return lastFilledPtr
}

func getFileInfo(diskContents []int, lastFilledPtr int) FileInfo {
	currFileID := diskContents[lastFilledPtr]
	_, fileStart, _ := lo.FindLastIndexOf(diskContents[:lastFilledPtr], generateAntiPred(currFileID))
	fileStart++
	fileEnd := lastFilledPtr + 1
	fileLength := fileEnd - fileStart
	return FileInfo{ID: currFileID, Length: fileLength, Start: fileStart, End: fileEnd}
}

// Note that the values of freeSpacePtr (as well as any other pointers to disk
// locations currently being held) should be considered invalidated once this
// function has applied.
// Disk contents are modified in-place.
func moveFile(diskContents *[]int, fileInfo FileInfo, freeSpacePtr int) {
	emptyRun := lib.RunOf(fileInfo.Length, lib.FreeSpaceIndicator)
	copy((*diskContents)[freeSpacePtr:], (*diskContents)[fileInfo.Start:fileInfo.End])
	copy((*diskContents)[fileInfo.Start:], emptyRun)
}

func calcChecksum(diskContents []int) *big.Int {
	checkSum := big.NewInt(0)
	for i, val := range diskContents {
		if val == lib.FreeSpaceIndicator {
			continue
		}
		checkSum.Add(checkSum, big.NewInt(int64(i*val)))
	}
	return checkSum
}
