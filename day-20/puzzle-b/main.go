package main

import (
	"bufio"
	"log"
	"main/lib"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/hashicorp/go-set/v3"
	"github.com/samber/lo"
	pq "gopkg.in/dnaeon/go-priorityqueue.v1"
)

type Args struct {
	InputFile               string `arg:"positional,required" help:"input file"`
	DepthOfCheat            int    `arg:"-d,--depth,required" help:"depth of cheat window"`
	ThresholdForImprovement int    `arg:"-t,--threshold,required" help:"threshold of improvement to consider"`
}

const (
	nothingFound = lib.Cost(-1)
)

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

	revMaze := *maze
	revMaze.Start, revMaze.End = revMaze.End, revMaze.Start
	bestNoCheatingPrice, prevs, dijkstraBoard := doDijkstra(revMaze)
	// // if endState == nil {
	// // 	log.Panic("endState is nil")
	// // }
	log.Printf("best price without cheating: %d", bestNoCheatingPrice)
	paths := collectPaths(lib.State{Pos: *revMaze.End}, prevs)
	log.Printf("number of paths: %d", len(paths))

	if len(paths) != 1 {
		log.Panicf("len(paths)=%d != 1", len(paths))
	}

	path := lo.Reverse(paths[0])
	log.Printf("path: %v", path)

	improverCounts := walkPath(revMaze, path, dijkstraBoard, args.DepthOfCheat, lib.Cost(args.ThresholdForImprovement))
	nImprover := lo.Sum(lo.Values(improverCounts))
	log.Printf("number of improvers: %d", nImprover)
}

func walkPath(maze lib.Maze, path []lib.Coord, dijkstraBoard [][]lib.Cost, depthOfCheat int, thresholdForImprovement lib.Cost) map[lib.Cost]int {
	improverCounts := make(map[lib.Cost]int)
	for _, pos := range path {
		curCell := dijkstraBoard[pos.Row][pos.Col]
		for dRow := -depthOfCheat; dRow <= depthOfCheat; dRow++ {
			colDepth := depthOfCheat - lib.Abs(dRow)
			for dCol := -colDepth; dCol <= colDepth; dCol++ {
				dijkstraCoord := lib.Coord{Row: pos.Row + dRow, Col: pos.Col + dCol}
				if !dijkstraCoord.IsValid(maze.Dimensions) {
					continue
				}

				dijkstraCell := dijkstraBoard[dijkstraCoord.Row][dijkstraCoord.Col]
				if dijkstraCell == nothingFound {
					continue
				}

				improvement := dijkstraCell - curCell - lib.Cost(lib.Distance(pos, dijkstraCoord))
				if improvement < thresholdForImprovement {
					continue
				}

				improverCounts[improvement]++
			}
		}
	}

	return improverCounts
}

func collectPaths(state lib.State, prevs map[lib.State][]lib.State) [][]lib.Coord {
	if _, ok := prevs[state]; !ok {
		return [][]lib.Coord{{state.Pos}}
	}

	allPaths := make([][]lib.Coord, 0)
	for _, prevState := range prevs[state] {
		paths := collectPaths(prevState, prevs)
		for _, path := range paths {
			amendedPath := path
			amendedPath = append(amendedPath, state.Pos)
			allPaths = append(allPaths, amendedPath)
		}
	}

	return allPaths
}

func doDijkstra(maze lib.Maze) (lib.Cost, map[lib.State][]lib.State, [][]lib.Cost) {
	initState := lib.State{
		Pos: *maze.Start,
	}

	bestCostByState := make(map[lib.State]lib.Cost)
	bestCostByState[initState] = 0
	unvisitedQueue := pq.New[lib.State, float64](pq.MinHeap)
	unvisitedQueue.Put(initState, 0)
	removed := set.New[lib.State](0)
	prevs := make(map[lib.State][]lib.State)
	for !unvisitedQueue.IsEmpty() {
		item := unvisitedQueue.Get()
		state, priority := item.Value, item.Priority
		cost := lib.Cost(priority)

		removed.Insert(state)
		for _, move := range maze.Moves {
			if !move.Precondition(maze, state) {
				continue
			}

			nextState, costIncrease := move.Func(state)
			if removed.Contains(nextState) {
				continue
			}

			nextCost := cost + costIncrease
			if oldBest, ok := bestCostByState[nextState]; ok {
				switch {
				case nextCost < oldBest:
					unvisitedQueue.Update(nextState, float64(nextCost))
					prevs[nextState] = []lib.State{state}
				case nextCost == oldBest:
					prevs[nextState] = append(prevs[nextState], state)
				case nextCost > oldBest:
					nextCost = oldBest
				}
			} else {
				unvisitedQueue.Put(nextState, float64(nextCost))
				prevs[nextState] = []lib.State{state}
			}
			bestCostByState[nextState] = nextCost
		}
	}

	dijkstraBoard := make([][]lib.Cost, maze.Dimensions.Row)
	for iRow := range maze.Dimensions.Row {
		dijkstraBoard[iRow] = make([]lib.Cost, maze.Dimensions.Col)
		for iCol := range maze.Dimensions.Col {
			state := lib.State{Pos: lib.Coord{Row: iRow, Col: iCol}}
			var depth lib.Cost
			if maze.Board[iRow][iCol] == lib.Wall {
				depth = nothingFound
			} else {
				depth = bestCostByState[state]
			}

			dijkstraBoard[iRow][iCol] = depth
		}
	}

	return bestCostByState[lib.State{Pos: *maze.End}], prevs, dijkstraBoard
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
	maze, err := lib.ReadInput(scanner)

	return maze, err //nolint:wrapcheck // Toy code
}
