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

		legal, failIdx := isLegal(values)
		if legal {
			nSafe++
			continue
		}

		repairedValues := make([]int, nValues-1)
		startCheck := max(0, failIdx-1)
		endCheck := min(nValues-1, failIdx+1)
		for idxToSkip := startCheck; idxToSkip <= endCheck; idxToSkip++ {
			copy(repairedValues[:idxToSkip], values[:idxToSkip])
			copy(repairedValues[idxToSkip:], values[idxToSkip+1:])
			legal, _ := isLegal(repairedValues)
			if legal {
				nSafe++
				break
			}
		}
	}

	log.Println(nSafe)
}

func isLegal(values []int) (bool, int) {
	ascending := false
	descending := false
	for idx := range len(values) - 1 {
		step := values[idx+1] - values[idx]
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
