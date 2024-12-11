package lib

import (
	"bufio"
	"container/list"
	"math/big"
	"strings"
)

func ReadInput(scanner *bufio.Scanner) *list.List {
	theList := list.New()
	for scanner.Scan() {
		line := scanner.Text()
		for _, str := range strings.Fields(line) {
			var value big.Int
			_, success := value.SetString(str, 10) //nolint:mnd // false positive
			if !success {
				continue
			}
			theList.PushBack(&value)
		}
	}

	return theList
}
