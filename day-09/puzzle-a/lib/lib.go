package lib

import (
	"bufio"

	"github.com/samber/lo"
)

const FreeSpaceIndicator = -1

func ReadInput(scanner *bufio.Scanner) []int {
	diskContents := make([]int, 0)
	for scanner.Scan() {
		line := scanner.Text()
		nextValIsFreeSpace := false
		fileID := 0
		for _, c := range line {
			val := int(c - '0')
			if nextValIsFreeSpace {
				diskContents = append(diskContents, runOf(val, FreeSpaceIndicator)...)
			} else {
				diskContents = append(diskContents, runOf(val, fileID)...)
				fileID++
			}
			nextValIsFreeSpace = !nextValIsFreeSpace
		}
	}

	return diskContents
}

func runOf[T any](runLength int, val T) []T {
	return lo.Map(lo.Range(runLength), func(_, _ int) T {
		return val
	})
}
