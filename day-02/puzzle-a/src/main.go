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

		ascending := false
		descending := false
		legal := true
		for i := range len(values) - 1 {
			step := values[i+1] - values[i]
			if step > 0 {
				ascending = true
			} else if step < 0 {
				descending = true
			}
			if ascending && descending {
				legal = false
				break
			}

			absStep := Abs(step)
			if (absStep < 1) || (absStep > 3) {
				legal = false
				break
			}
		}

		if legal {
			nSafe++
		}
	}

	log.Println(nSafe)
}
