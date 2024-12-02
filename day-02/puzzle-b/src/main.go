package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/samber/lo"
	"golang.org/x/exp/constraints"
)

func Abs[T constraints.Signed](x T) T { //nolint:ireturn // false positive
	return max(x, -x)
}

func main() {
	file, err := os.Open("../input/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	scanner := bufio.NewScanner(file)
	nSafe := 0
	for scanner.Scan() {
		line := scanner.Text()

		fields := strings.Fields(line)
		values := lo.Map(fields, func(item string, _ int) int {
			num, err := strconv.Atoi(item)
			if err != nil {
				log.Panic(err)
			}
			return num
		})

		nValues := len(values)
		if nValues < 1 {
			continue
		}

		legal, failIdx := isLegal(values, -1)
		if legal {
			nSafe++
			continue
		}

		startCheck := max(0, failIdx-1)
		endCheck := min(nValues-1, failIdx+1)
		for idxToSkip := startCheck; idxToSkip <= endCheck; idxToSkip++ {
			legal, _ := isLegal(values, idxToSkip)
			if legal {
				nSafe++
				break
			}
		}
	}

	log.Println(nSafe)
}

func isLegal(values []int, skipIdx int) (bool, int) {
	ascending := false
	descending := false
	for idx := 0; idx < len(values)-1; idx++ {
		if idx == skipIdx {
			idx++
		}

		nextIdx := idx + 1
		if nextIdx == skipIdx {
			nextIdx++
		}

		if nextIdx >= len(values) {
			break
		}

		step := values[nextIdx] - values[idx]
		if step > 0 {
			ascending = true
		} else if step < 0 {
			descending = true
		}
		if ascending && descending {
			return false, idx
		}

		absStep := Abs(step)
		if (absStep < 1) || (absStep > 3) {
			return false, idx
		}
	}

	return true, -1
}
