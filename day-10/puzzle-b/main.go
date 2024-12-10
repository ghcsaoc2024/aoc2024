package main

import (
	"bufio"
	"log"
	"os"

	"main/lib"

	"github.com/alexflint/go-arg"
)

type Args struct {
	InputFile string `arg:"positional,required" help:"input file"`
}

var directions = []lib.Coord{ //nolint:gochecknoglobals // Meant as a constant
	{Row: 1, Col: 0},
	{Row: 0, Col: 1},
	{Row: 0, Col: -1},
	{Row: -1, Col: 0},
}

func main() {
	var args Args
	arg.MustParse(&args)

	board := readInputFile(args)
	if len(board.Grid) < 1 {
		log.Panic("grid is empty")
	}

	dimensions := lib.Coord{Row: len(board.Grid), Col: len(board.Grid[0])}
	totalScore := 0
	for elevation := lib.TopElevation - 1; elevation >= lib.BottomElevation; elevation-- {
		for _, coord := range board.ByElevation[elevation] {
			trailCount := 0
			for _, dir := range directions {
				neighborCoord := coord.Add(lib.Coord{Row: dir.Row, Col: dir.Col})
				if !neighborCoord.IsValid(dimensions) {
					continue
				}

				neighborCell := &board.Grid[neighborCoord.Row][neighborCoord.Col]
				if neighborCell.Elevation != elevation+1 {
					continue
				}

				trailCount += neighborCell.TrailCount
			}

			cell := &board.Grid[coord.Row][coord.Col]
			cell.TrailCount = trailCount
			if elevation > lib.BottomElevation {
				continue
			}

			totalScore += trailCount
		}
	}

	log.Printf("grid dimensions: %v", dimensions)
	log.Printf("total score: %d", totalScore)
}

func readInputFile(args Args) lib.Board {
	file, err := os.Open(args.InputFile)
	if err != nil {
		log.Fatal(err) //nolint:revive // Toy code
	}
	defer func(file *os.File) {
		closeErr := file.Close()
		if closeErr != nil {
			log.Fatal(closeErr) //nolint:revive // Toy code
		}
	}(file)

	scanner := bufio.NewScanner(file)
	board := lib.ReadInput(scanner)

	return board
}
