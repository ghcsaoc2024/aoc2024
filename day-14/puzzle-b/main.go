package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"main/lib"

	"github.com/alexflint/go-arg"
)

type Args struct {
	InputFile               string `arg:"positional,required" help:"input file"`
	X                       int64  `arg:"-x, --x-dimension"   default:"101"     help:"X dimension of the board"`
	Y                       int64  `arg:"-y, --y-dimension"   default:"103"     help:"Y dimension of the board"`
	SecondsToFF             int64  `arg:"-s, --seconds"       default:"100"     help:"seconds to fast-forward"`
	DisplayAfter            int64  `arg:"-d, --display-after" default:"-1"      help:"display board after this many seconds"`
	MinQuadDisplayThreshold int    `arg:"-i, --min-threshold" default:"-1"      help:"display board if it has a quad count at or below this value"`
	MaxQuadDisplayThreshold int    `arg:"-a, --max-threshold" default:"-1"      help:"display board if it has a quad count at or above this value"`
}

func main() {
	var args Args
	arg.MustParse(&args)

	if args.X < 2 {
		log.Fatalf("X dimension of board must be at least 2; got %d", args.X)
	}
	if args.Y < 2 {
		log.Fatalf("Y dimension of board must be at least 2; got %d", args.Y)
	}

	dimensions := lib.Coord{X: args.X, Y: args.Y}

	robots, err := readInputFile(args)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("%d robots read", len(robots))

	if args.MaxQuadDisplayThreshold < 0 {
		args.MaxQuadDisplayThreshold = len(robots)
	}

	// displayAll(robots, dimensions)

	// Advance all robots
	midpoints := dimensions.Div(2)
	overallMinQuadrantCount := len(robots)
	overallMaxQuadrantCount := 0
	for iSec := range args.SecondsToFF {
		for iRobot := range robots {
			robot := &robots[iRobot]
			robot.Pos = robot.Pos.Add(robot.Vel).Add(dimensions).ModOther(dimensions)
		}

		quadrantCounts := genQuadrantCounts(robots, midpoints)

		minQuadrantCount := len(robots)
		maxQuadrantCount := 0
		for x := range 2 {
			for y := range 2 {
				minQuadrantCount = min(minQuadrantCount, quadrantCounts[x][y])
				maxQuadrantCount = max(maxQuadrantCount, quadrantCounts[x][y])
			}
		}
		overallMinQuadrantCount = min(overallMinQuadrantCount, minQuadrantCount)
		overallMaxQuadrantCount = max(overallMaxQuadrantCount, maxQuadrantCount)

		// log.Printf("quadrant counts after %d seconds: %v", iSec, quadrantCounts)
		if (iSec == args.DisplayAfter) || (minQuadrantCount <= args.MinQuadDisplayThreshold) || (maxQuadrantCount >= args.MaxQuadDisplayThreshold) {
			displayAll(robots, dimensions)
			log.Printf("(this is after %d seconds)", iSec+1)
		}
	}

	log.Printf("min quadrant count: %d", overallMinQuadrantCount)
	log.Printf("max quadrant count: %d", overallMaxQuadrantCount)

	quadrantCounts := genQuadrantCounts(robots, midpoints)

	log.Printf("quadrant counts: %v", quadrantCounts)

	product := int64(1)
	for x := range 2 {
		for y := range 2 {
			product *= int64(quadrantCounts[x][y])
		}
	}

	log.Printf("product: %d", product)
}

func genQuadrantCounts(robots []lib.Robot, midpoints lib.Coord) [2][2]int {
	var quadrantCounts [2][2]int
	divPoints := midpoints.Add(lib.Coord{X: 1, Y: 1})
	for _, robot := range robots {
		if robot.Pos.X == midpoints.X || robot.Pos.Y == midpoints.Y {
			continue
		}

		quadrants := robot.Pos.DivOther(divPoints)
		quadrantCounts[quadrants.X][quadrants.Y]++
	}
	return quadrantCounts
}

func displayAll(robots []lib.Robot, dimensions lib.Coord) {
	board := make([][]bool, dimensions.X)
	for x := range board {
		board[x] = make([]bool, dimensions.Y)
	}

	for _, robot := range robots {
		board[robot.Pos.X][robot.Pos.Y] = true
	}

	for x := range board {
		for y := range board[x] {
			if board[x][y] {
				fmt.Printf("#") //nolint:forbidigo // Intentional console output
			} else {
				fmt.Printf(".") //nolint:forbidigo // Intentional console output
			}
		}
		fmt.Printf("\n") //nolint:forbidigo // Intentional console output
	}
}

func readInputFile(args Args) ([]lib.Robot, error) {
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
