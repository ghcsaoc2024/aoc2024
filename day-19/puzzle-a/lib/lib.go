package lib

import (
	"bufio"
	"errors"
	"strings"
)

func ReadInput(scanner *bufio.Scanner) ([]string, []string, error) {
	var inventory []string
	for scanner.Scan() && inventory == nil {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if len(trimmed) < 1 {
			continue
		}

		towels := strings.Split(line, ", ")
		inventory = append(inventory, towels...)
	}

	if inventory == nil {
		return nil, nil, errors.New("no inventory found")
	}

	var patterns []string
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if len(trimmed) < 1 {
			continue
		}

		patterns = append(patterns, trimmed)
	}

	return inventory, patterns, nil
}
