package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"main/lib"
	"os"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/hashicorp/go-set/v3"
	"github.com/samber/lo"
	pq "gopkg.in/dnaeon/go-priorityqueue.v1"
)

type Args struct {
	InputFile string `arg:"positional,required" help:"input file"`
}

type NumPadLayoutMap map[lib.Coord]int
type NumPadRevLayoutMap map[int]lib.Coord
type ActionPadLayoutMap map[lib.Coord]int
type ActionPadRevLayoutMap map[int]lib.Coord

type NumPadMaps struct {
	Layout           NumPadLayoutMap
	RevLayout        NumPadRevLayoutMap
	RunesByNumPadKey map[int]rune
}

type ActionPadMaps struct {
	Layout        ActionPadLayoutMap
	RevLayout     ActionPadRevLayoutMap
	RunesByAction map[int]rune
}

type AllMaps struct {
	NumPad    NumPadMaps
	ActionPad ActionPadMaps
}

type State struct {
	// actionPad1      string
	// actionPad2Coord lib.Coord
	actionPad2      string
	actionPad3Coord lib.Coord
	actionPad3      string
	numPadCoord     lib.Coord
	numPad          string
}

var errOutOfBounds = errors.New("out-of-bounds")
var errSolutionNotFound = errors.New("solution not found")

func main() {
	var args Args
	arg.MustParse(&args)

	numPadCodes, err := readInputFile(args)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("num pad codes: %d", numPadCodes)
	log.Printf("total: %d", len(numPadCodes))

	allMaps := AllMaps{NumPad: makeNumPadMaps(), ActionPad: makeActionPadMaps()}
	log.Printf("num pad layout map: %v", allMaps.NumPad.Layout)
	log.Printf("num pad rev layout map: %v", allMaps.NumPad.RevLayout)
	log.Printf("action pad layout map: %v", allMaps.ActionPad.Layout)
	log.Printf("action rev layout map: %v", allMaps.ActionPad.RevLayout)

	testStr := "<vA<AA>>^AvAA<^A>A<v<A>>^AvA^A<vA>^A<v<A>^A>AAvA^A<v<A>A>^AAAvA<^A>A"
	log.Printf("test string: %s", testStr)
	result, err := execUberSequence(testStr, allMaps)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("action sequence: %s", result)

	resultInts := lo.Map([]rune(result), func(r rune, _ int) int {
		return lib.NumPadKeyByRune[r]
	})
	solution, err := doDijkstra(resultInts, allMaps)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("solution: %s", solution)
}

func execUberSequence(s string, allMaps AllMaps) (string, error) {
	actionSeq := lo.Map([]rune(s), func(r rune, _ int) int {
		return lo.ValueOr(lib.ActionByRune, r, lib.InvalidAction)
	})
	log.Printf("action sequence: %v", actionSeq)

	actionSeq2, err := execSeqOnActionPad(actionSeq, allMaps.ActionPad)
	if err != nil {
		return "", err
	}

	actionSeq3, err := execSeqOnActionPad(actionSeq2, allMaps.ActionPad)
	if err != nil {
		return "", err
	}

	keyPresses, err := execSeqOnNumPad(actionSeq3, allMaps.NumPad)
	if err != nil {
		return "", err
	}

	log.Printf("key presses: %v", keyPresses)

	return string(lo.Map(keyPresses, func(key, _ int) rune {
		return allMaps.NumPad.RunesByNumPadKey[key]
	})), nil
}

func doDijkstra(numPadKeyPresses []int, allMaps AllMaps) (string, error) {
	numPadString := string(lo.Map(numPadKeyPresses, func(key int, _ int) rune {
		return allMaps.NumPad.RunesByNumPadKey[key]
	}))
	initState := State{
		// actionPad2Coord: allMaps.ActionPad.RevLayout[lib.Press],
		// actionPad3Coord: allMaps.ActionPad.RevLayout[lib.Press],
		numPadCoord: allMaps.NumPad.RevLayout[lib.NumPadKeyA],
	}
	bestCostByState := make(map[State]int)
	unvisitedQueue := pq.New[State, float64](pq.MinHeap)
	unvisitedQueue.Put(initState, 0)
	removed := set.New[State](0)
	for !unvisitedQueue.IsEmpty() {
		item := unvisitedQueue.Get()
		state, priority := item.Value, item.Priority
		cost := int(priority)

		removed.Insert(state)
		if state.numPad == numPadString {
			return state.actionPad3, nil
		}

		for nextRune := range lib.ActionByRune {
			newState := state
			newCost := cost
			var action int
			var err error
			// action, err = processAction(
			// 	nextRune,
			// 	&newState.actionPad1,
			// 	&newState.actionPad2Coord,
			// 	allMaps.ActionPad.Layout,
			// )
			// if errors.Is(err, errOutOfBounds) {
			// 	continue
			// } else if err != nil {
			// 	return "", err
			// }
			//
			// if action != lib.InvalidAction {
			// 	nextRune = allMaps.ActionPad.RunesByAction[action]
			// action, err = processAction(
			// 	nextRune,
			// 	&newState.actionPad2,
			// 	&newState.actionPad3Coord,
			// 	allMaps.ActionPad.Layout,
			// )
			// if errors.Is(err, errOutOfBounds) {
			// 	continue
			// } else if err != nil {
			// 	return "", err
			// }
			// }
			//
			// if action != lib.InvalidAction {
			// nextRune = allMaps.ActionPad.RunesByAction[action]
			action, err = processAction(
				nextRune,
				&newState.actionPad3,
				&newState.numPadCoord,
				allMaps.NumPad.Layout,
			)
			if errors.Is(err, errOutOfBounds) {
				continue
			} else if err != nil {
				return "", err
			}
			// }

			if action != lib.InvalidNumPadKey {
				nextRune = allMaps.NumPad.RunesByNumPadKey[action]
				newState.numPad += string(nextRune)
				// newCost = len(newState.actionPad1)
				// newCost = len(newState.actionPad2)
			}

			if !strings.HasPrefix(numPadString, newState.numPad) {
				continue
			}

			newCost = len(newState.actionPad3)
			stateKey := State{
				actionPad3:  "",
				numPadCoord: newState.numPadCoord,
				numPad:      newState.numPad,
			}
			if oldBest, ok := bestCostByState[stateKey]; ok {
				switch {
				case newCost < oldBest:
					unvisitedQueue.Update(newState, float64(newCost))
				case newCost > oldBest:
					newCost = oldBest
				}
			} else {
				unvisitedQueue.Put(newState, float64(newCost))
			}
			bestCostByState[stateKey] = newCost
		}
	}

	return "", errSolutionNotFound
}

func processAction(r rune, currentStr *string, nextCoords *lib.Coord, nextLayout map[lib.Coord]int) (int, error) {
	*currentStr += string(r)
	action := lib.ActionByRune[r]
	switch action {
	case lib.Press:
		newAction := nextLayout[*nextCoords]
		return newAction, nil
	default:
		*nextCoords = nextCoords.Add(lib.Actions[action])
		if !lo.HasKey(nextLayout, *nextCoords) {
			return lib.InvalidAction, errOutOfBounds
		}
		return lib.InvalidAction, nil
	}
}

func execSeqOnNumPad(actions []int, numPadMaps NumPadMaps) ([]int, error) {
	coord := numPadMaps.RevLayout[lib.NumPadKeyA]
	keyPresses := make([]int, 0)
	for _, action := range actions {
		switch action {
		case lib.InvalidAction:
			return nil, errors.New("encountered invalid action while executing numpad sequence")
		case lib.Press:
			keyPresses = append(keyPresses, numPadMaps.Layout[coord])
		default:
			dir := lib.Actions[action]
			nextCoord := coord.Add(dir)
			if !lo.HasKey(numPadMaps.Layout, nextCoord) {
				return nil, fmt.Errorf("out-of-bounds while executing numpad sequence: %v", nextCoord)
			}
			coord = nextCoord
		}
	}

	return keyPresses, nil
}

func execSeqOnActionPad(actions []int, actionPadMaps ActionPadMaps) ([]int, error) {
	coord := actionPadMaps.RevLayout[lib.Press]
	keyPresses := make([]int, 0)
	for _, action := range actions {
		switch action {
		case lib.InvalidAction:
			return nil, errors.New("encountered invalid action while executing actionpad sequence")
		case lib.Press:
			keyPresses = append(keyPresses, actionPadMaps.Layout[coord])
		default:
			dir := lib.Actions[action]
			nextCoord := coord.Add(dir)
			if !lo.HasKey(actionPadMaps.Layout, nextCoord) {
				return nil, fmt.Errorf("out-of-bounds while executing actionpad sequence: %v", nextCoord)
			}
			coord = nextCoord
		}
	}

	return keyPresses, nil
}

func makeActionPadMaps() ActionPadMaps {
	actionPadMaps := ActionPadMaps{}
	actionPadMaps.Layout = make(ActionPadLayoutMap)
	for iRow, row := range lib.ActionPadLayout {
		for iCol, key := range row {
			if key == lib.InvalidAction {
				continue
			}

			actionPadMaps.Layout[lib.Coord{Row: iRow, Col: iCol}] = key
		}
	}

	actionPadMaps.RevLayout = make(ActionPadRevLayoutMap)
	for iRow, row := range lib.ActionPadLayout {
		for iCol, key := range row {
			if key == lib.InvalidAction {
				continue
			}

			actionPadMaps.RevLayout[key] = lib.Coord{Row: iRow, Col: iCol}
		}
	}

	actionPadMaps.RunesByAction = lo.Invert(lib.ActionByRune)

	return actionPadMaps
}

func makeNumPadMaps() NumPadMaps {
	numPadMaps := NumPadMaps{}
	numPadMaps.Layout = make(NumPadLayoutMap)
	for iRow, row := range lib.NumPadLayout {
		for iCol, key := range row {
			if key == lib.InvalidNumPadKey {
				continue
			}

			numPadMaps.Layout[lib.Coord{Row: iRow, Col: iCol}] = key
		}
	}

	numPadMaps.RevLayout = make(NumPadRevLayoutMap)
	for iRow, row := range lib.NumPadLayout {
		for iCol, key := range row {
			if key == lib.InvalidNumPadKey {
				continue
			}

			numPadMaps.RevLayout[key] = lib.Coord{Row: iRow, Col: iCol}
		}
	}

	numPadMaps.RunesByNumPadKey = lo.Invert(lib.NumPadKeyByRune)

	return numPadMaps
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
