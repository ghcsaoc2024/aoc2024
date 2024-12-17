package main

import (
	"bufio"
	"log"
	"main/lib"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/hashicorp/go-set/v3"
	pq "gopkg.in/dnaeon/go-priorityqueue.v1"
)

type Args struct {
	InputFile string `arg:"positional,required" help:"input file"`
}

const NothingFound = lib.Cost(-1)

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

	bestPrice, nGoodSeats := traverse(*maze)
	log.Printf("best cost: %d", bestPrice)
	log.Printf("number of good seats: %d", nGoodSeats)
}

func traverse(wholeMaze lib.Maze) (lib.Cost, int) {
	var prevsByEndCursor map[lib.Cursor]map[lib.Cursor][]lib.Cursor
	var bestPrice lib.Cost
	for _, dir := range lib.Directions {
		actualEndCursor := lib.Cursor{Coord: *wholeMaze.End, Dir: dir}
		revMaze := wholeMaze
		revMaze.End, revMaze.Start = wholeMaze.Start, wholeMaze.End
		revMaze.Cursor = actualEndCursor
		revMaze.EndCursor = lib.Cursor{Coord: *revMaze.End, Dir: lib.StartDirection.Mul(-1)}
		cost, prevs := doDijkstra(revMaze)
		if cost == NothingFound {
			continue
		}
		if len(prevsByEndCursor) < 1 || cost < bestPrice {
			bestPrice = cost
			prevsByEndCursor = make(map[lib.Cursor]map[lib.Cursor][]lib.Cursor)
			prevsByEndCursor[revMaze.EndCursor] = prevs
		}
	}

	goodSeats := set.New[lib.Coord](0)
	for cursor, prevs := range prevsByEndCursor {
		goodSeats.InsertSet(collectGoodSeats(prevs, cursor))
	}

	return bestPrice, goodSeats.Size()
}

func doDijkstra(wholeMaze lib.Maze) (lib.Cost, map[lib.Cursor][]lib.Cursor) {
	bestCostByCursor := make(map[lib.Cursor]lib.Cost)
	bestCostByCursor[wholeMaze.Cursor] = 0
	unvisitedQueue := pq.New[lib.Cursor, float64](pq.MinHeap)
	unvisitedQueue.Put(wholeMaze.Cursor, 0)
	removed := set.New[lib.Cursor](1)

	prevs := make(map[lib.Cursor][]lib.Cursor)
	for !unvisitedQueue.IsEmpty() {
		item := unvisitedQueue.Get()
		cursor, priority := item.Value, item.Priority
		cost := lib.Cost(priority)
		removed.Insert(cursor)
		if cursor == wholeMaze.EndCursor {
			return lib.Cost(priority), prevs
		}

		for _, move := range lib.Moves {
			if !move.Precondition(cursor, wholeMaze) {
				continue
			}

			nextCursor, costIncr := move.Func(cursor)
			if removed.Contains(nextCursor) {
				continue
			}

			nextCost := cost + costIncr
			if oldBest, ok := bestCostByCursor[nextCursor]; ok {
				switch {
				case nextCost < oldBest:
					unvisitedQueue.Update(nextCursor, float64(nextCost))
					prevs[nextCursor] = []lib.Cursor{cursor}
				case nextCost == oldBest:
					prevs[nextCursor] = append(prevs[nextCursor], cursor)
				case nextCost > oldBest:
					nextCost = oldBest
				}
			} else {
				unvisitedQueue.Put(nextCursor, float64(nextCost))
				prevs[nextCursor] = []lib.Cursor{cursor}
			}
			bestCostByCursor[nextCursor] = nextCost
		}
	}

	return NothingFound, nil
}

func collectGoodSeats(prevs map[lib.Cursor][]lib.Cursor, cursor lib.Cursor) *set.Set[lib.Coord] {
	goodSeats := set.New[lib.Coord](1)
	goodSeats.Insert(cursor.Coord)
	for _, prevCursor := range prevs[cursor] {
		goodSeats.InsertSet(collectGoodSeats(prevs, prevCursor))
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
