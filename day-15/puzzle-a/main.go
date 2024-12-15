package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"main/lib"

	"github.com/alexflint/go-arg"
	"github.com/samber/lo"
)

type Args struct {
	InputFile string `arg:"positional,required" help:"input file"`
}

func main() {
	var args Args
	arg.MustParse(&args)

	game, err := readInputFile(args)
	if err != nil {
		log.Panic(err)
	}

	if len(game.Board) < 1 {
		log.Panic("board is empty")
	}

	dimensions := lib.Coord{Row: len(game.Board), Col: len(game.Board[0])}
	log.Printf("game board dimensions: %v", dimensions)
	log.Printf("number of boxes: %d", len(game.Boxes))
	log.Printf("initial robot position: %v", game.Robot)
	log.Printf("number of moves: %d", len(game.Moves))

	for _, move := range game.Moves {
		nextCoords := game.Robot.Add(move)
		if !nextCoords.IsValid(dimensions) {
			continue
		}

		nextCell := &game.Board[nextCoords.Row][nextCoords.Col]
		if *nextCell == lib.Wall {
			continue
		}

		if *nextCell == lib.Empty {
			game.Robot = nextCoords
			continue
		}

		pushDest, ok := evalPush(*game, dimensions, move)
		if !ok {
			continue
		}

		err := execPush(game, dimensions, move, *pushDest)
		if err != nil {
			log.Panic(err)
		}
	}

	log.Printf("final robot position: %v", game.Robot)

	totalScore := lo.Sum(lo.Map(game.Boxes, func(box lib.Coord, _ int) int {
		return 100*box.Row + box.Col
	}))
	log.Printf("total score: %d", totalScore)
}

func execPush(game *lib.Game, dimensions lib.Coord, move lib.Coord, pushDest lib.Coord) error {
	for {
		curCell := &game.Board[pushDest.Row][pushDest.Col]
		nextCoords := pushDest.Sub(move)
		if !nextCoords.IsValid(dimensions) {
			return fmt.Errorf("internal error: invalid coordinates %v while executing push", nextCoords)
		}
		nextCell := &game.Board[nextCoords.Row][nextCoords.Col]
		switch *nextCell { //nolint:exhaustive // has default
		case lib.Box:
			*curCell = *nextCell
			iBox := game.BoxesByCoord[nextCoords]
			delete(game.BoxesByCoord, nextCoords)
			game.Boxes[iBox] = pushDest
			game.BoxesByCoord[pushDest] = iBox
		case lib.Empty:
			*curCell = *nextCell
			game.Robot = pushDest
			return nil
		default:
			return fmt.Errorf("internal error: unexpected cell type %v while executing push", *nextCell)
		}
		pushDest = nextCoords
	}
}

func evalPush(game lib.Game, dimensions lib.Coord, move lib.Coord) (*lib.Coord, bool) {
	for nextCoords := game.Robot.Add(move); nextCoords.IsValid(dimensions); nextCoords = nextCoords.Add(move) {
		nextCell := &game.Board[nextCoords.Row][nextCoords.Col]
		switch *nextCell {
		case lib.Box:
			continue
		case lib.Wall:
			return &nextCoords, false
		case lib.Empty:
			return &nextCoords, true
		}
	}

	return nil, false
}

func readInputFile(args Args) (*lib.Game, error) {
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
	values, err := lib.ReadInput(scanner)

	return values, err //nolint:wrapcheck // Toy code
}
