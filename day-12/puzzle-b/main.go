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

	cornerDict := make([][2]int, 4) //nolint:mnd // Four corners
	for iCorner, corner := range lib.Corners {
		items := corner.Slice()
		if len(items) != 2 {
			log.Panic("internal error: problem in lib.Corners (len(items) != 2)")
		}

		dir1, dir2 := items[0], items[1]
		iDir1 := lo.IndexOf(lib.Directions, dir1)
		iDir2 := lo.IndexOf(lib.Directions, dir2)
		if iDir1 < 0 || iDir2 < 0 {
			log.Panic("internal error: problem in lib.Corners (iDir1 < 0 || iDir2 < 0)")
		}

		cornerDict[iCorner] = [2]int{iDir1, iDir2}
	}

	totalCorners := 0
	totalArea := 0
	totalCost := 0
	for allCoords.Size() > 0 {
		next, _ := iter.Pull(allCoords.Items())
		startingCoord, ok := next()
		if !ok {
			log.Panic("internal error: allCoords should not be empty")
		}
		subArea, subCorners := expand(&board, dimensions, startingCoord, allCoords, cornerDict)
		totalCorners += subCorners
		totalArea += subArea
		totalCost += subArea * subCorners
	}

	log.Printf("total area: %d", totalArea)
	log.Printf("total corners: %d", totalCorners)
	log.Printf("total cost: %d", totalCost)
}

func expand(board *[][]lib.Cell, dimensions, startingCoord lib.Coord, allCoords *set.Set[lib.Coord], cornerDict [][2]int) (int, int) {
	if !allCoords.Contains(startingCoord) {
		return 0, 0
	}

	allCoords.Remove(startingCoord)

	cell := &(*board)[startingCoord.Row][startingCoord.Col]
	selfCorners := 0
	for _, dirs := range cornerDict {
		if cell.Boundaries[dirs[0]] && cell.Boundaries[dirs[1]] {
			selfCorners++
		}
	}

	convexCorners := 0
	for _, dirs := range cornerDict {
		if cell.Boundaries[dirs[0]] || cell.Boundaries[dirs[1]] {
			continue
		}

		diagNeighborCoords := startingCoord.Add(lib.Directions[dirs[0]]).Add(lib.Directions[dirs[1]])
		if !diagNeighborCoords.IsValid(dimensions) {
			continue
		}

		diagNeighbor := &(*board)[diagNeighborCoords.Row][diagNeighborCoords.Col]
		if diagNeighbor.Kind != cell.Kind {
			convexCorners++
		}
	}

	totalCorners := selfCorners + convexCorners
	totalArea := 1
	for iDir, dir := range lib.Directions {
		if cell.Boundaries[iDir] {
			continue
		}
		neighborCoord := startingCoord.Add(dir)
		if !neighborCoord.IsValid(dimensions) {
			continue
		}

		if !allCoords.Contains(neighborCoord) {
			continue
		}

		subArea, subCorners := expand(board, dimensions, neighborCoord, allCoords, cornerDict)
		totalCorners += subCorners
		totalArea += subArea
	}

	return totalArea, totalCorners
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
