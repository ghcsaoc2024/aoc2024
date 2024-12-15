package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"

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

		pushDest, boxesToPush, ok := evalPush(*game, dimensions, move)
		if !ok {
			continue
		}

		err := execPush(game, dimensions, move, *pushDest, boxesToPush)
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

func execPush(game *lib.Game, dimensions, move, pushDest lib.Coord, boxesToPush []int) error {
	if move.Row != 0 {
		return execPushVert(game, dimensions, move, pushDest, boxesToPush)
	}

	return execPushHoriz(game, dimensions, move, pushDest, boxesToPush)
}

func execPushHoriz(game *lib.Game, dimensions, move, pushDest lib.Coord, _ []int) error {
	for {
		curCell := &game.Board[pushDest.Row][pushDest.Col]
		nextCoords := pushDest.Sub(move)
		if !nextCoords.IsValid(dimensions) {
			return fmt.Errorf("internal error: invalid coordinates %v while executing push", nextCoords)
		}
		nextCell := &game.Board[nextCoords.Row][nextCoords.Col]
		switch *nextCell { //nolint:exhaustive // has default
		case lib.BoxL:
			iBox := game.BoxesByCoord[nextCoords]
			delete(game.BoxesByCoord, nextCoords)
			game.Boxes[iBox] = pushDest
			game.BoxesByCoord[pushDest] = iBox
			fallthrough
		case lib.BoxR:
			*curCell = *nextCell
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

func execPushVert(game *lib.Game, _, move, _ lib.Coord, boxesToPush []int) error {
	slices.Reverse(boxesToPush)
	for _, iBox := range boxesToPush {
		coord := game.Boxes[iBox]
		if err := moveBoxVert(game, coord, coord.Add(move)); err != nil {
			return err
		}
	}

	game.Robot = game.Robot.Add(move)

	return nil
}

func moveBoxVert(game *lib.Game, nextCoords lib.Coord, newCoords lib.Coord) error {
	iBox := game.BoxesByCoord[nextCoords]
	delete(game.BoxesByCoord, nextCoords)
	game.Boxes[iBox] = newCoords
	game.BoxesByCoord[newCoords] = iBox
	if game.Board[newCoords.Row][newCoords.Col] != lib.Empty {
		return fmt.Errorf("internal error: box %v is not empty", newCoords)
	}
	if game.Board[newCoords.Row][newCoords.Col+1] != lib.Empty {
		return fmt.Errorf("internal error: box %v is not empty", newCoords.Add(lib.Coord{Row: 0, Col: 1}))
	}
	game.Board[newCoords.Row][newCoords.Col] = lib.BoxL
	game.Board[nextCoords.Row][nextCoords.Col] = lib.Empty
	game.Board[newCoords.Row][newCoords.Col+1] = lib.BoxR
	game.Board[nextCoords.Row][nextCoords.Col+1] = lib.Empty

	return nil
}

func evalPush(game lib.Game, dimensions lib.Coord, move lib.Coord) (*lib.Coord, []int, bool) {
	if move.Row != 0 {
		return evalPushVert(game, dimensions, move)
	}

	return evalPushHoriz(game, dimensions, move)
}

func evalPushHoriz(game lib.Game, dimensions lib.Coord, move lib.Coord) (*lib.Coord, []int, bool) {
	for nextCoords := game.Robot.Add(move); nextCoords.IsValid(dimensions); nextCoords = nextCoords.Add(move) {
		nextCell := &game.Board[nextCoords.Row][nextCoords.Col]
		switch *nextCell {
		case lib.BoxL, lib.BoxR:
			continue
		case lib.Wall:
			return &nextCoords, nil, false
		case lib.Empty:
			return &nextCoords, nil, true
		}
	}

	return nil, nil, false
}

func evalPushVert(game lib.Game, dimensions lib.Coord, move lib.Coord) (*lib.Coord, []int, bool) {
	nextCoords := game.Robot.Add(move)
	nextCell := &game.Board[nextCoords.Row][nextCoords.Col]
	lastRowBoxes := make([]int, 1)
	switch *nextCell { //nolint:exhaustive // has default
	case lib.BoxL:
		lastRowBoxes[0] = game.BoxesByCoord[nextCoords]
	case lib.BoxR:
		lastRowBoxes[0] = game.BoxesByCoord[lib.Coord{Row: nextCoords.Row, Col: nextCoords.Col - 1}]
	default:
		log.Panicf("internal error: unexpected cell type %v while evaluating first row of vert-push", *nextCell) //nolint:revive // Toy code
	}

	boxesToMove := make([]int, 0, 1)
	for nextCoords = nextCoords.Add(move); nextCoords.IsValid(dimensions); nextCoords = nextCoords.Add(move) {
		boxesToMove = append(boxesToMove, lastRowBoxes...)
		nextRowBoxes := set.New[int](0)
		for _, iBox := range lastRowBoxes {
			box := game.Boxes[iBox]
			boxShadowCoordL := box.Add(move)
			boxShadowCoordR := lib.Coord{Row: boxShadowCoordL.Row, Col: box.Col + 1}
			boxShadowL := game.Board[boxShadowCoordL.Row][boxShadowCoordL.Col]
			boxShadowR := game.Board[boxShadowCoordR.Row][boxShadowCoordR.Col]
			if boxShadowL == lib.Wall || boxShadowR == lib.Wall {
				return nil, nil, false
			}
			if boxShadowL == lib.Empty && boxShadowR == lib.Empty {
				continue
			}

			if boxShadowL == lib.BoxR {
				nextRowBoxes.Insert(game.BoxesByCoord[lib.Coord{Row: boxShadowCoordL.Row, Col: boxShadowCoordL.Col - 1}])
			}
			if boxShadowR == lib.BoxL {
				nextRowBoxes.Insert(game.BoxesByCoord[boxShadowCoordR])
			}
			if boxShadowL == lib.BoxL {
				nextRowBoxes.Insert(game.BoxesByCoord[boxShadowCoordL])
			}
		}

		if nextRowBoxes.Empty() {
			return &nextCoords, boxesToMove, true
		}

		lastRowBoxes = nextRowBoxes.Slice()
	}

	return nil, nil, false
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
