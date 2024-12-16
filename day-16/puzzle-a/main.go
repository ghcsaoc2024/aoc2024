package main

import (
	"bufio"
	"log"
	"main/lib"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/samber/lo"
)

type Args struct {
	InputFile string `arg:"positional,required" help:"input file"`
}

func main() {
	var args Args
	arg.MustParse(&args)

	maze, err := readInputFile(args)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("maze dimensions: %v", maze.Dimensions)
	log.Printf("starting coord: %v", maze.Start)
	log.Printf("ending coord: %v", maze.End)
	log.Printf("current cursor: %v", maze.Cursor)

	bestPaths := make(map[lib.Cursor]lib.Cost)
	cost := traverse(*maze, bestPaths)

	log.Printf("cost: %d", cost)
}

func traverse(state lib.Maze, bestPaths map[lib.Cursor]lib.Cost) lib.Cost {
	cheapest, ok := bestPaths[state.Cursor]
	if ok && cheapest <= state.Cost {
		return -1
	}

	bestPaths[state.Cursor] = state.Cost
	if state.Cursor.Coord == *state.End {
		endCursors := lo.Map(lib.Directions, func(dir lib.Coord, _ int) lib.Cursor {
			return lib.Cursor{Coord: state.Cursor.Coord, Dir: dir}
		})
		successfulEndCursors := lo.Filter(endCursors, func(c lib.Cursor, _ int) bool {
			_, ok := bestPaths[c]
			return ok
		})
		costs := lo.Map(successfulEndCursors, func(c lib.Cursor, _ int) lib.Cost {
			return bestPaths[c]
		})

		return lo.Min(costs)
	}

	costs := make([]lib.Cost, 0, len(lib.Moves))
	for _, move := range lib.Moves {
		if !move.Precondition(state) {
			continue
		}

		cost := traverse(move.Func(state), bestPaths)
		if cost > 0 {
			costs = append(costs, cost)
		}
	}

	if len(costs) < 1 {
		return -1
	}

	return lo.Min(costs)
}

func readInputFile(args Args) (*lib.Maze, error) {
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
