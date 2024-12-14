package main

import (
	"bufio"
	"log"
	"os"

	"main/lib"

	"github.com/alexflint/go-arg"
)

type Args struct {
	InputFile   string `arg:"positional,required" help:"input file"`
	X           int64  `arg:"-x, --x-dimension"   default:"101"     help:"X dimension of the board"`
	Y           int64  `arg:"-y, --y-dimension"   default:"103"     help:"Y dimension of the board"`
	SecondsToFF int64  `arg:"-s, --seconds"       default:"100"     help:"seconds to fast-forward"`
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

	// Advance all robots
	for iRobot := range robots {
		robot := &robots[iRobot]
		robot.Pos = robot.Vel.Mul(args.SecondsToFF).Add(robot.Pos).ModOther(dimensions).Add(dimensions).ModOther(dimensions)
	}

	midpoints := dimensions.Div(2)
	divPoints := midpoints.Add(lib.Coord{X: 1, Y: 1})
	var quadrantCounts [2][2]int64
	for _, robot := range robots {
		if robot.Pos.X == midpoints.X || robot.Pos.Y == midpoints.Y {
			continue
		}

		quadrants := robot.Pos.DivOther(divPoints)
		quadrantCounts[quadrants.X][quadrants.Y]++
	}

	log.Printf("quadrant counts: %v", quadrantCounts)

	product := int64(1)
	for x := range 2 {
		for y := range 2 {
			product *= quadrantCounts[x][y]
		}
	}

	log.Printf("product: %d", product)
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
