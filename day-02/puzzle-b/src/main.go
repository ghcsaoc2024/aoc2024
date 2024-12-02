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

		legal := isLegal(values)
		if legal {
			nSafe++
			continue
		}

		repairedValues := make([]int, len(values)-1)
		for i := range values {
			copy(repairedValues[:i], values[:i])
			copy(repairedValues[i:], values[i+1:])
			if isLegal(repairedValues) {
				nSafe++
				break
			}
		}
	}

	log.Println(nSafe)
}

func isLegal(values []int) bool {
	ascending := false
	descending := false
	for i := range len(values) - 1 {
		step := values[i+1] - values[i]
		if step > 0 {
			ascending = true
		} else if step < 0 {
			descending = true
		}
		if ascending && descending {
			return false
		}

		absStep := Abs(step)
		if (absStep < 1) || (absStep > 3) {
			return false
		}
	}

	return true
}
