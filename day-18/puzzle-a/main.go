package main

import (
	"bufio"
	"log"
	"os"

	"main/lib"

	"github.com/alexflint/go-arg"
	"github.com/hashicorp/go-set/v3"
	"github.com/samber/lo"
	pq "gopkg.in/dnaeon/go-priorityqueue.v1"
)

type Args struct {
	InputFile    string `arg:"positional,required" help:"input file"`
	BoardDimRows int    `arg:"-r,--board-dim-rows" default:"7" help:"board dimension rows"`
	BoardDimCols int    `arg:"-c,--board-dim-cols" default:"7" help:"board dimension cols"`
	StartRow     int    `arg:"-s,--start-row" default:"0" help:"starting row"`
	StartCol     int    `arg:"-t,--start-col" default:"0" help:"starting col"`
	EndRow       int    `arg:"-e,--end-row" default:"-1" help:"ending row"`
	EndCol       int    `arg:"-f,--end-col" default:"-1" help:"ending col"`
	NumSteps     int    `arg:"-n,--num-steps,required" help:"number of steps to execute before eval"`
}

func main() {
	var args Args
	arg.MustParse(&args)

	game, err := readInputFile(args)
	if err != nil {
		log.Panic(err)
	}

	if args.EndRow < 0 {
		args.EndRow = ((args.EndRow % args.BoardDimRows) + args.BoardDimRows) % args.BoardDimRows
	}
	if args.EndCol < 0 {
		args.EndCol = ((args.EndCol % args.BoardDimCols) + args.BoardDimCols) % args.BoardDimCols
	}

	game.StartPos = lib.Coord{Row: args.StartRow, Col: args.StartCol}
	game.EndPos = lib.Coord{Row: args.EndRow, Col: args.EndCol}

	log.Printf("game: %v", game)

	board := lib.Board(lo.Map(lo.Range(args.BoardDimRows), func(_, _ int) []bool {
		row := make([]bool, args.BoardDimCols)
		return row
	}))
	for step := 1; step <= args.NumSteps; step++ {
		loc, ok := game.BlockSched[step]
		if !ok {
			log.Printf("step %d: no block scheduled", step)
		}
		board[loc.Row][loc.Col] = true
	}

	getBoardState := func(_ int) *lib.Board {
		return &board
	}

	pathLength, _ := doDijkstra(*game, getBoardState)

	log.Printf("path length: %d", pathLength)
}

// func deepcopyBoard(board lib.Board) *lib.Board {
// 	var copyBoard lib.Board
// 	for _, row := range board {
// 		copyBoard = append(copyBoard, row[:]) //nolint:gocritic // This is intentional; we're doing a deepcopy.
// 	}
//
// 	return &copyBoard
// }

// func traverse(wholeMaze lib.Maze) (int, int) {
// 	var prevsByEndCoord map[lib.Coord]map[lib.Coord][]lib.Coord
// 	var bestPrice int
// 	for _, dir := range lib.Directions {
// 		actualEndCoord := lib.Coord{Coord: *wholeMaze.End, Dir: dir}
// 		revMaze := wholeMaze
// 		revMaze.End, revMaze.Start = wholeMaze.Start, wholeMaze.End
// 		revMaze.Coord = actualEndCoord
// 		revMaze.EndCoord = lib.Coord{Coord: *revMaze.End, Dir: lib.StartDirection.Mul(-1)}
// 		cost, prevs := doDijkstra(revMaze)
// 		if cost == NothingFound {
// 			continue
// 		}
// 		if len(prevsByEndCoord) < 1 || cost < bestPrice {
// 			bestPrice = cost
// 			prevsByEndCoord = make(map[lib.Coord]map[lib.Coord][]lib.Coord)
// 			prevsByEndCoord[revMaze.EndCoord] = prevs
// 		}
// 	}
//
// 	goodSeats := set.New[lib.Coord](0)
// 	for cursor, prevs := range prevsByEndCoord {
// 		goodSeats.InsertSet(collectGoodSeats(prevs, cursor))
// 	}
//
// 	return bestPrice, goodSeats.Size()
// }

func doDijkstra(game lib.Game, getBoardState func(int) *lib.Board) (int, map[lib.Coord][]lib.Coord) {
	bestCostByCoord := make(map[lib.Coord]int)
	bestCostByCoord[game.StartPos] = 0
	unvisitedQueue := pq.New[lib.Coord, float64](pq.MinHeap)
	unvisitedQueue.Put(game.StartPos, 0)
	removed := set.New[lib.Coord](1)

	prevs := make(map[lib.Coord][]lib.Coord)
	for !unvisitedQueue.IsEmpty() {
		item := unvisitedQueue.Get()
		coord, priority := item.Value, item.Priority
		cost := int(priority)
		removed.Insert(coord)
		if coord == game.EndPos {
			return int(priority), prevs
		}

		for _, dir := range lib.Directions {
			nextCoord := coord.Add(dir)
			if !nextCoord.IsValid(game.Dims) {
				continue
			}

			if removed.Contains(nextCoord) {
				continue
			}

			boardState := getBoardState(cost)
			if (*boardState)[nextCoord.Row][nextCoord.Col] {
				continue
			}

			nextCost := cost + 1
			if oldBest, ok := bestCostByCoord[nextCoord]; ok {
				switch {
				case nextCost < oldBest:
					unvisitedQueue.Update(nextCoord, float64(nextCost))
					prevs[nextCoord] = []lib.Coord{coord}
				case nextCost == oldBest:
					prevs[nextCoord] = append(prevs[nextCoord], coord)
				case nextCost > oldBest:
					nextCost = oldBest
				}
			} else {
				unvisitedQueue.Put(nextCoord, float64(nextCost))
				prevs[nextCoord] = []lib.Coord{coord}
			}
			bestCostByCoord[nextCoord] = nextCost
		}
	}

	return -1, nil
}

// func collectGoodSeats(prevs map[lib.Coord][]lib.Coord, cursor lib.Coord) *set.Set[lib.Coord] {
// 	goodSeats := set.New[lib.Coord](1)
// 	goodSeats.Insert(cursor.Coord)
// 	for _, prevCoord := range prevs[cursor] {
// 		goodSeats.InsertSet(collectGoodSeats(prevs, prevCoord))
// 	}
//
// 	return goodSeats
// }

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
	game, err := lib.ReadInput(scanner, lib.Coord{Row: args.BoardDimRows, Col: args.BoardDimCols})

	return game, err //nolint:wrapcheck // Toy code
}
