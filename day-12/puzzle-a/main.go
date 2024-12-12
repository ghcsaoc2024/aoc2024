package main

import (
	"bufio"
	"iter"
	"log"
	"os"

	"main/lib"

	"github.com/alexflint/go-arg"
	"github.com/hashicorp/go-set/v3"
	"github.com/samber/lo"
)

type Args struct {
	InputFile string `arg:"positional,required" help:"input file"`
}

func main() {
	var args Args
	arg.MustParse(&args)

	board := readInputFile(args)
	if len(board) < 1 {
		log.Panic("board is empty")
	}

	dimensions := lib.Coord{Row: len(board), Col: len(board[0])}

	antiDirs := make(map[int]int)
	for iDir, dir := range lib.Directions {
		antiDir := lo.IndexOf(lib.Directions, lib.Coord{Row: -dir.Row, Col: -dir.Col})
		if antiDir < 0 {
			log.Panic("invalid anti-direction")
		}
		antiDirs[iDir] = antiDir
	}

	for iRow := range dimensions.Row {
		for iCol := range dimensions.Col {
			cell := &board[iRow][iCol]
			for iDir, dir := range lib.Directions {
				neighborCoord := cell.Coord.Add(dir)
				if !neighborCoord.IsValid(dimensions) {
					continue
				}
				neighborCell := &board[neighborCoord.Row][neighborCoord.Col]
				if cell.Kind == neighborCell.Kind {
					cell.Boundaries[iDir] = false
					neighborCell.Boundaries[antiDirs[iDir]] = false
				}
			}
		}
	}

	allCoords := set.New[lib.Coord](dimensions.Row * dimensions.Col)
	for iRow := range dimensions.Row {
		for iCol := range dimensions.Col {
			allCoords.Insert(lib.Coord{Row: iRow, Col: iCol})
		}
	}

	totalFence := 0
	totalArea := 0
	totalCost := 0
	for allCoords.Size() > 0 {
		next, _ := iter.Pull(allCoords.Items())
		startingCoord, ok := next()
		if !ok {
			log.Panic("internal error: allCoords should not be empty")
		}
		subArea, subFence := expand(&board, dimensions, startingCoord, allCoords)
		totalFence += subFence
		totalArea += subArea
		totalCost += subArea * subFence
	}

	log.Printf("total area: %d", totalArea)
	log.Printf("total fence: %d", totalFence)
	log.Printf("total cost: %d", totalCost)
}

func expand(board *[][]lib.Cell, dimensions, startingCoord lib.Coord, allCoords *set.Set[lib.Coord]) (int, int) {
	if !allCoords.Contains(startingCoord) {
		return 0, 0
	}

	allCoords.Remove(startingCoord)

	totalFence := 0
	totalArea := 1
	for iDir, dir := range lib.Directions {
		cell := &(*board)[startingCoord.Row][startingCoord.Col]
		if cell.Boundaries[iDir] {
			totalFence++
			continue
		}
		neighborCoord := startingCoord.Add(dir)
		if !neighborCoord.IsValid(dimensions) {
			continue
		}

		if !allCoords.Contains(neighborCoord) {
			continue
		}

		subArea, subFence := expand(board, dimensions, neighborCoord, allCoords)
		totalFence += subFence
		totalArea += subArea
	}

	return totalArea, totalFence
}

func readInputFile(args Args) [][]lib.Cell {
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
	values := lib.ReadInput(scanner)

	return values
}
