package main

import (
	"bufio"
	"log"
	"os"

	"main/lib"

	"github.com/hashicorp/go-set/v3"
)

func main() {
	file, err := os.Open("../input/input.txt")
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

	// Read in the array
	dimensions, antennae := lib.ReadArray(scanner)

	// Find the antinodes
	antinodes := set.New[lib.Coord](0)
	for _, locs := range antennae {
		for _, loc1 := range locs {
			for _, loc2 := range locs {
				if loc2 == loc1 {
					continue
				}
				diff := loc2.Subtract(loc1)
				projection := loc2.Add(diff)
				if projection.IsValid(dimensions) {
					antinodes.Insert(projection)
				}
			}
		}
	}

	log.Printf("number of antinodes found: %d", antinodes.Size())
}
