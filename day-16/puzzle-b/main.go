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

const NothingFound = lib.Cost(-1)
const CutoffReached = lib.Cost(-2)

func main() {
	var args Args
	arg.MustParse(&args)

	maze, err := readInputFile(args)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("maze dimensions: %v", maze.Dimensions)
	log.Printf("starting coord: %v", *maze.Start)
	log.Printf("ending coord: %v", *maze.End)
	log.Printf("current cursor: %v", maze.Cursor)

	bestCostByCursor := make(map[lib.Cursor]lib.Cost)

	bestCost := CutoffReached
	for cutoff := lib.Cost(1); bestCost == CutoffReached; cutoff *= 2 {
		log.Printf("attempting traversal with cutoff: %v", cutoff)
		cost := traverse(*maze, bestCostByCursor, cutoff)
		if cost > 0 {
			bestCost = cost
		}
	}

	log.Printf("phase 1 complete; bestCost: %d", bestCost)
	for _, cursor := range maze.EndCursors {
		bestCostByCursor[cursor] = bestCost
	}

	goodSeats := findGoodSeats(*maze, bestCostByCursor)
	log.Printf("good seats: %d", goodSeats.Size())
}

func traverse(state lib.Maze, bestCostByCursor map[lib.Cursor]lib.Cost, costCutoff lib.Cost) lib.Cost {
	if costCutoff > 0 && state.Cost >= costCutoff {
		return CutoffReached
	}

	cheapest, ok := bestCostByCursor[state.Cursor]
	if ok && cheapest < state.Cost {
		return NothingFound
	}

	if state.Cursor.Coord == *state.End {
		successfulEndCursors := lo.Filter(state.EndCursors, func(c lib.Cursor, _ int) bool {
			_, ok := bestCostByCursor[c]
			return ok
		})
		costs := lo.Map(successfulEndCursors, func(c lib.Cursor, _ int) lib.Cost {
			return bestCostByCursor[c]
		})

		cheapest = lo.Min(costs)
		if state.Cost < cheapest || len(costs) < 1 {
			cheapest = state.Cost
			log.Printf("new cheapest threshold: %v", cheapest)
			bestCostByCursor[state.Cursor] = cheapest
		}

		return cheapest
	}

	bestCostByCursor[state.Cursor] = state.Cost
	costs := make([]lib.Cost, 0, len(lib.Moves))
	for _, move := range lib.Moves {
		if !move.Precondition(state) {
			continue
		}

		cost := traverse(move.Func(state), bestCostByCursor, costCutoff)
		costs = append(costs, cost)
	}

	realCosts := lo.Filter(costs, func(c lib.Cost, _ int) bool {
		return c > 0
	})
	if len(realCosts) == 0 {
		return lo.Min(costs)
	}

	return lo.Min(realCosts)
}

func findGoodSeats(state lib.Maze, bestCostByCursor map[lib.Cursor]lib.Cost) *set.Set[lib.Coord] {
	cheapest, ok := bestCostByCursor[state.Cursor]
	if ok && cheapest < state.Cost {
		return nil
	}

	if state.Cursor.Coord == *state.End {
		cheapest = bestCostByCursor[state.Cursor]
		if state.Cost == cheapest {
			log.Printf("found a path (cost: %v)", cheapest)
			bestCostByCursor[state.Cursor] = cheapest
			goodSeats := set.New[lib.Coord](1)
			goodSeats.Insert(state.Cursor.Coord)
			return goodSeats
		}

		return nil
	}

	bestCostByCursor[state.Cursor] = state.Cost
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
