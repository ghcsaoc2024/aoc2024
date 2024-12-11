package lib

import (
	"bufio"
	"math/big"
	"strings"
)

func ReadInput(scanner *bufio.Scanner) []*big.Int {
	values := make([]*big.Int, 0)
	for scanner.Scan() {
		line := scanner.Text()
		for _, str := range strings.Fields(line) {
			var value big.Int
			_, success := value.SetString(str, 10) //nolint:mnd // false positive
			if !success {
				continue
			}
			values = append(values, &value)
		}
	}

	return values
}
