package main

import (
	"bufio"
	"log"
	"main/lib"
	"os"

	"github.com/alexflint/go-arg"
)

type Args struct {
	InputFile string `arg:"positional,required" help:"input file"`
}

const TurnCost = 1000

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

	budget := lib.Cost(1)
	states := []lib.Maze{*maze}
	for {
		log.Printf("budget: %d, len(states): %d", budget, len(states))
		states = traverse(states, budget)
		if len(states) == 0 {
			log.Printf("no solution found")
			break
		}

		if states[0].Solved {
			log.Printf("solution found; cost: %d", states[0].Cost)
			break
		}

		budget += 1000
	}
}

func traverse(states []lib.Maze, budget lib.Cost) []lib.Maze {
	returnStates := make([]lib.Maze, 0, len(states))
	var winner *lib.Maze
	for _, maze := range states {
		c := lib.Coord{Row: maze.Cursor.Row, Col: maze.Cursor.Col}
		if c == *maze.End {
			maze.Solved = true
			return []lib.Maze{maze}
		}

		if maze.Cost >= budget {
			returnStates = append(returnStates, maze)
			continue
		}

		for _, move := range lib.Moves {
			if !move.Precondition(maze) {
				continue
			}
			localStates := traverse([]lib.Maze{move.Func(maze)}, budget)
			if len(localStates) > 0 && localStates[0].Solved {
				if winner == nil || localStates[0].Cost < winner.Cost {
					winner = &localStates[0]
				}
				continue
			}

			if winner != nil {
				continue
			}

			returnStates = append(returnStates, localStates...)
		}
	}

	if winner != nil {
		return []lib.Maze{*winner}
	}

	return returnStates
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
