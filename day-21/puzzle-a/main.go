package main

import (
	"bufio"
	"errors"
	"fmt"
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

	numPadCodes, err := readInputFile(args)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("num pad codes: %d", numPadCodes)
	log.Printf("total: %d", len(numPadCodes))

	numPadLayoutMap, actionPadLayoutMap := makeLayoutMaps()
	log.Printf("num pad layout map: %v", numPadLayoutMap)
	log.Printf("action pad layout map: %v", actionPadLayoutMap)

	numPadRevLayoutMap, actionRevLayoutMap := makeRevLayoutMaps()
	log.Printf("num pad rev layout map: %v", numPadRevLayoutMap)
	log.Printf("action rev layout map: %v", actionRevLayoutMap)

	actionSeq := lo.Map([]rune("<vA<AA>>^AvAA<^A>A<v<A>>^AvA^A<vA>^A<v<A>^A>AAvA^A<v<A>A>^AAAvA<^A>A"), func(r rune, _ int) lib.Action {
		return lo.ValueOr(lib.ActionByRune, r, lib.InvalidAction)
	})
	log.Printf("action sequence: %v", actionSeq)

	actionSeq2, err := execSeqOnActionPad(actionSeq, actionPadLayoutMap, actionRevLayoutMap)
	if err != nil {
		log.Panic(err)
	}

	actionSeq3, err := execSeqOnActionPad(actionSeq2, actionPadLayoutMap, actionRevLayoutMap)
	if err != nil {
		log.Panic(err)
	}

	keyPresses, err := execSeqOnNumPad(actionSeq3, numPadLayoutMap, numPadRevLayoutMap)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("key presses: %v", keyPresses)
}

func execSeqOnNumPad(actions []lib.Action, numPadLayoutMap map[lib.Coord]lib.NumPadKey, numPadRevLayoutMap map[lib.NumPadKey]lib.Coord) ([]lib.NumPadKey, error) {
	coord := numPadRevLayoutMap[lib.NumPadKeyA]
	keyPresses := make([]lib.NumPadKey, 0)
	for _, action := range actions {
		switch action {
		case lib.InvalidAction:
			return nil, errors.New("encountered invalid action while executing numpad sequence")
		case lib.Press:
			keyPresses = append(keyPresses, numPadLayoutMap[coord])
		default:
			dir := lib.Actions[action]
			nextCoord := coord.Add(dir)
			if !lo.HasKey(numPadLayoutMap, nextCoord) {
				return nil, fmt.Errorf("out-of-bounds while executing numpad sequence: %v", nextCoord)
			}
			coord = nextCoord
		}
	}

	return keyPresses, nil
}

func execSeqOnActionPad(actions []lib.Action, actionPadLayoutMap map[lib.Coord]lib.Action, actionPadRevLayoutMap map[lib.Action]lib.Coord) ([]lib.Action, error) {
	coord := actionPadRevLayoutMap[lib.Press]
	keyPresses := make([]lib.Action, 0)
	for _, action := range actions {
		switch action {
		case lib.InvalidAction:
			return nil, errors.New("encountered invalid action while executing actionpad sequence")
		case lib.Press:
			keyPresses = append(keyPresses, actionPadLayoutMap[coord])
		default:
			dir := lib.Actions[action]
			nextCoord := coord.Add(dir)
			if !lo.HasKey(actionPadLayoutMap, nextCoord) {
				return nil, fmt.Errorf("out-of-bounds while executing actionpad sequence: %v", nextCoord)
			}
			coord = nextCoord
		}
	}

	return keyPresses, nil
}

func makeRevLayoutMaps() (map[lib.NumPadKey]lib.Coord, map[lib.Action]lib.Coord) {
	numPadRevLayoutMap := make(map[lib.NumPadKey]lib.Coord)
	for iRow, row := range lib.NumPadLayout {
		for iCol, key := range row {
			if key == lib.InvalidNumPadKey {
				continue
			}

			numPadRevLayoutMap[key] = lib.Coord{Row: iRow, Col: iCol}
		}
	}

	actionRevLayoutMap := make(map[lib.Action]lib.Coord)
	for iRow, row := range lib.ActionPadLayout {
		for iCol, key := range row {
			if key == lib.InvalidAction {
				continue
			}

			actionRevLayoutMap[key] = lib.Coord{Row: iRow, Col: iCol}
		}
	}

	return numPadRevLayoutMap, actionRevLayoutMap
}

func makeLayoutMaps() (map[lib.Coord]lib.NumPadKey, map[lib.Coord]lib.Action) {
	numPadLayoutMap := make(map[lib.Coord]lib.NumPadKey)
	for iRow, row := range lib.NumPadLayout {
		for iCol, key := range row {
			if key == lib.InvalidNumPadKey {
				continue
			}

			numPadLayoutMap[lib.Coord{Row: iRow, Col: iCol}] = key
		}
	}

	actionPadLayoutMap := make(map[lib.Coord]lib.Action)
	for iRow, row := range lib.ActionPadLayout {
		for iCol, key := range row {
			if key == lib.InvalidAction {
				continue
			}

			actionPadLayoutMap[lib.Coord{Row: iRow, Col: iCol}] = key
		}
	}

	return numPadLayoutMap, actionPadLayoutMap
}

func readInputFile(args Args) ([]lib.NumPadCode, error) {
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
	numPadCodes, err := lib.ReadInput(scanner)

	return numPadCodes, err //nolint:wrapcheck // Toy code
}
