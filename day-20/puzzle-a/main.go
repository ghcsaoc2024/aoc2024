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
	InputFile      string `arg:"positional,required" help:"input file"`
	CheatThreshold int    `arg:"-n,--cheat-threshold,required" help:"cheat threshold"`
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
	log.Printf("current pos: %v", maze.Pos)

	noCheatMaze := *maze
	noCheatMaze.HasCheated = true
	bestNoCheatingPrice, prevs := doDijkstra(noCheatMaze)
	log.Printf("best price without cheating: %d", bestNoCheatingPrice)
	endState := lib.State{
		Pos:        *maze.End,
		HasCheated: true,
	}
	nNoCheatingPaths, _ := collectPaths(*maze, endState, prevs)
	log.Printf("number of paths: %d", nNoCheatingPaths)

	maze.BlockedCheats = set.New[lib.Coord](0)
	nTotalCheatPaths := 0
	cheatCost := lib.Cost(args.CheatThreshold)
	for cheatCost < bestNoCheatingPrice {
		cheatMaze := *maze
		cheatMaze.HasCheated = false
		cheatMaze.CheatCost = lib.Cost(args.CheatThreshold)
		bestPenalizedPrice, prevs := doDijkstra(cheatMaze)
		switch {
		case bestPenalizedPrice > bestNoCheatingPrice:
			log.Panicf("bestPenalizedPrice=%d > bestNoCheatingPrice=%d", bestPenalizedPrice, bestNoCheatingPrice)
		case bestPenalizedPrice == bestNoCheatingPrice:
			cheatCost++
		case bestPenalizedPrice < bestNoCheatingPrice:
			nNewPaths, newBlockedCheats := collectPaths(*maze, endState, prevs)
			for _, pos := range newBlockedCheats {
				maze.BlockedCheats.Insert(pos)
			}
			nTotalCheatPaths += nNewPaths
			cheatCost = bestPenalizedPrice
			log.Printf("number of new paths: %d (cheatCost=%d)", nNewPaths, cheatCost)
		}
	}

	log.Printf("number of total cheat paths: %d", nTotalCheatPaths)
}

func collectPaths(maze lib.Maze, state lib.State, prevs map[lib.State][]lib.State) (int, []lib.Coord) {
	if _, ok := prevs[state]; !ok {
		return 1, nil
	}

	nPaths := 0
	var cheatList []lib.Coord
	isCheat := (maze.Board[state.Pos.Row][state.Pos.Col] == lib.Wall)
	cheatList = make([]lib.Coord, 0)
	if isCheat {
		cheatList = append(cheatList, state.Pos)
	}

	for _, prevState := range prevs[state] {
		nNewPaths, newCheats := collectPaths(maze, prevState, prevs)
		nPaths += nNewPaths
		cheatList = append(cheatList, newCheats...)
	}
	return nPaths, cheatList
}

func doDijkstra(maze lib.Maze) (lib.Cost, map[lib.State][]lib.State) {
	initState := lib.State{
		Pos:        *maze.Start,
		HasCheated: maze.HasCheated,
	}
	bestCostByState := make(map[lib.State]lib.Cost)
	bestCostByState[initState] = 0
	unvisitedQueue := pq.New[lib.State, float64](pq.MinHeap)
	unvisitedQueue.Put(initState, 0)
	var removed *set.Set[lib.State]
	if maze.BlockedCheats != nil {
		removed = set.New[lib.State](maze.BlockedCheats.Size())
		for pos := range maze.BlockedCheats.Items() {
			removed.Insert(lib.State{
				Pos:        pos,
				HasCheated: true,
			})
		}
	} else {
		removed = set.New[lib.State](0)
	}

	prevs := make(map[lib.State][]lib.State)
	for !unvisitedQueue.IsEmpty() {
		item := unvisitedQueue.Get()
		state, priority := item.Value, item.Priority
		cost := lib.Cost(priority)
		removed.Insert(state)
		if state.Pos == *maze.End {
			return lib.Cost(priority), prevs
		}

		for _, move := range lib.Moves {
			maze.Pos = state.Pos
			maze.HasCheated = state.HasCheated
			if !move.Precondition(maze) {
				continue
			}

			nextPos := move.Func(state.Pos)
			nextIsWall := maze.Board[nextPos.Row][nextPos.Col] == lib.Wall
			if nextIsWall && state.HasCheated {
				continue
			}

			nextState := lib.State{
				Pos:        nextPos,
				HasCheated: state.HasCheated || nextIsWall,
			}

			if removed.Contains(nextState) {
				continue
			}

			var nextCost lib.Cost
			if nextIsWall {
				nextCost = cost + maze.CheatCost
			} else {
				nextCost = cost + 1
			}

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

	return NothingFound, nil
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
