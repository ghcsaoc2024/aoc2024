package main

import (
	"bufio"
	"log"
	"main/lib"
	"os"

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

	maze, err := readInputFile(args)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("maze dimensions: %v", maze.Dimensions)
	log.Printf("starting coord: %v", maze.Start)
	log.Printf("ending coord: %v", maze.End)
	log.Printf("current cursor: %v", maze.Cursor)

	bestCostByCursor := make(map[lib.Cursor]lib.Cost)
	cost := traverse(*maze, bestCostByCursor)

	log.Printf("phase 1 complete; best cost: %d", cost)

	goodSeats := findGoodSeats(*maze, bestCostByCursor)
	log.Printf("good seats: %d", goodSeats.Size())
}

func traverse(state lib.Maze, bestCostByCursor map[lib.Cursor]lib.Cost) lib.Cost {
	cheapest, ok := bestCostByCursor[state.Cursor]
	if ok && cheapest <= state.Cost {
		return -1
	}

	bestCostByCursor[state.Cursor] = state.Cost
	if state.Cursor.Coord == *state.End {
		endCursors := lo.Map(lib.Directions, func(dir lib.Coord, _ int) lib.Cursor {
			return lib.Cursor{Coord: state.Cursor.Coord, Dir: dir}
		})
		successfulEndCursors := lo.Filter(endCursors, func(c lib.Cursor, _ int) bool {
			_, ok := bestCostByCursor[c]
			return ok
		})
		costs := lo.Map(successfulEndCursors, func(c lib.Cursor, _ int) lib.Cost {
			return bestCostByCursor[c]
		})

		return lo.Min(costs)
	}

	costs := make([]lib.Cost, 0, len(lib.Moves))
	for _, move := range lib.Moves {
		if !move.Precondition(state) {
			continue
		}

		cost := traverse(move.Func(state), bestCostByCursor)
		if cost > 0 {
			costs = append(costs, cost)
		}
	}

	if len(costs) < 1 {
		return -1
	}

	return lo.Min(costs)
}

func findGoodSeats(state lib.Maze, bestCostByCursor map[lib.Cursor]lib.Cost) *set.Set[lib.Coord] {
	cheapest, ok := bestCostByCursor[state.Cursor]
	if ok && cheapest < state.Cost {
		return nil
	}

	bestCostByCursor[state.Cursor] = state.Cost
	if state.Cursor.Coord == *state.End {
		endCursors := lo.Map(lib.Directions, func(dir lib.Coord, _ int) lib.Cursor {
			return lib.Cursor{Coord: state.Cursor.Coord, Dir: dir}
		})
		successfulEndCursors := lo.Filter(endCursors, func(c lib.Cursor, _ int) bool {
			_, ok := bestCostByCursor[c]
			return ok
		})
		costs := lo.Map(successfulEndCursors, func(c lib.Cursor, _ int) lib.Cost {
			return bestCostByCursor[c]
		})

		cheapest = lo.Min(costs)
		if cheapest == state.Cost {
			log.Printf("found a new path to cheapest cost: %d", cheapest)
			goodSeats := set.New[lib.Coord](1)
			goodSeats.Insert(state.Cursor.Coord)
			return goodSeats
		}

		return nil
	}

	var goodSeats *set.Set[lib.Coord]
	for _, move := range lib.Moves {
		if !move.Precondition(state) {
			continue
		}

		moreSeats := findGoodSeats(move.Func(state), bestCostByCursor)
		if moreSeats == nil {
			continue
		}

		if goodSeats == nil {
			goodSeats = moreSeats
			continue
		}

		goodSeats.InsertSet(moreSeats)
	}

	if goodSeats != nil {
		goodSeats.Insert(state.Cursor.Coord)
	}

	return goodSeats
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
