package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"main/lib"
	"os"
	"regexp"
	"strconv"
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
	CurrentStr string
	NextCoord  lib.Coord
	NextStr    string
}

var errOutOfBounds = errors.New("out-of-bounds")

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

	total := int64(0)
	for _, numPadCode := range numPadCodes {
		numPadString := string(lo.Map(numPadCode, func(key, _ int) rune {
			return allMaps.NumPad.RunesByNumPadKey[key]
		}))
		value, err := solve(numPadString, allMaps)
		if err != nil {
			log.Panic(err)
		}
		log.Printf("num pad string: %s", numPadString)
		log.Printf("value: %d", value)
		total += value
	}

	log.Printf("total: %d", total)
}

func solve(target string, allMaps AllMaps) (int64, error) {
	solutions3, err := doDijkstra(target, allMaps.NumPad.RunesByNumPadKey, allMaps.NumPad.Layout, allMaps.NumPad.RevLayout[lib.NumPadKeyA])
	if err != nil {
		return -1, err
	}
	log.Printf("number of solutions3: %d", len(solutions3))

	var solutions2 []string
	for _, solution3 := range solutions3 {
		subSolutions2, err := doDijkstra(solution3, allMaps.ActionPad.RunesByAction, allMaps.ActionPad.Layout, allMaps.ActionPad.RevLayout[lib.Press])
		if err != nil {
			return -1, err
		}
		// log.Printf("subSolutions2: %v", subSolutions2)
		solutions2 = append(solutions2, subSolutions2...)
	}
	log.Printf("number of solutions2: %d", len(solutions2))

	minLength2 := lo.Min(lo.Map(solutions2, func(s string, _ int) int {
		return len(s)
	}))
	filteredSolutions2 := lo.Filter(solutions2, func(s string, _ int) bool {
		return len(s) == minLength2
	})
	log.Printf("number of solutions2 after filtering: %d", len(filteredSolutions2))

	var solutions1 []string
	for _, solution2 := range filteredSolutions2 {
		subSolutions1, err := doDijkstra(solution2, allMaps.ActionPad.RunesByAction, allMaps.ActionPad.Layout, allMaps.ActionPad.RevLayout[lib.Press])
		if err != nil {
			return -1, err
		}
		// log.Printf("subSolutions1: %v", subSolutions1)
		solutions1 = append(solutions1, subSolutions1...)
	}
	log.Printf("number of solutions1: %d", len(solutions1))

	minLength1 := lo.Min(lo.Map(solutions1, func(s string, _ int) int {
		return len(s)
	}))
	filteredSolutions1 := lo.Filter(solutions1, func(s string, _ int) bool {
		return len(s) == minLength1
	})
	log.Printf("number of solutions1 after filtering: %d", len(filteredSolutions1))

	for _, solution1 := range filteredSolutions1 {
		numPadString, err := execUberSequence(solution1, allMaps)
		if err != nil {
			return -1, err
		}
		if numPadString != target {
			return -1, fmt.Errorf("internal error: num pad string does not match target: `%s` != `%s`", numPadString, target)
		}
	}
	log.Printf("verified %d solutions", len(filteredSolutions1))

	log.Printf("minLength1: %d", minLength1)

	pattern := `^0*([1-9][0-9]*)A$`
	re, err := regexp.Compile(pattern)
	if err != nil {
		return -1, fmt.Errorf("internal error: failed to compile regex: %w", err)
	}

	match := re.FindStringSubmatch(target)
	if len(match) != 2 {
		return -1, fmt.Errorf("internal error: failed to find match in target string: %s", target)
	}

	value, err := strconv.ParseInt(match[1], 10, 64)
	if err != nil {
		return -1, fmt.Errorf("internal error: failed to parse match: %w", err)
	}

	return value * int64(minLength1), nil
}

func execUberSequence(s string, allMaps AllMaps) (string, error) {
	actionSeq := lo.Map([]rune(s), func(r rune, _ int) int {
		return lo.ValueOr(lib.ActionByRune, r, lib.InvalidAction)
	})

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

	return string(lo.Map(keyPresses, func(key, _ int) rune {
		return allMaps.NumPad.RunesByNumPadKey[key]
	})), nil
}

func doDijkstra(targetString string, runesByKey map[int]rune, nextLayout map[lib.Coord]int, initialCoord lib.Coord) ([]string, error) {
	initState := State{
		NextCoord: initialCoord,
	}
	bestCostByState := make(map[State]int)
	unvisitedQueue := pq.New[State, float64](pq.MinHeap)
	unvisitedQueue.Put(initState, 0)
	removed := set.New[State](0)
	results := make([]string, 0)
	for !unvisitedQueue.IsEmpty() {
		item := unvisitedQueue.Get()
		state, _ := item.Value, item.Priority

		removed.Insert(state)
		if state.NextStr == targetString {
			results = append(results, state.CurrentStr)
			continue
		}

		for nextRune := range lib.ActionByRune {
			newState := state
			var action int
			var err error
			action, err = processAction(
				nextRune,
				&newState.CurrentStr,
				&newState.NextCoord,
				nextLayout,
			)
			if errors.Is(err, errOutOfBounds) {
				continue
			} else if err != nil {
				return nil, err
			}
			// }

			if action != lib.InvalidNumPadKey {
				nextRune = runesByKey[action]
				newState.NextStr += string(nextRune)
			}

			if !strings.HasPrefix(targetString, newState.NextStr) {
				continue
			}

			newCost := len(newState.CurrentStr)
			stateKey := State{
				CurrentStr: "",
				NextCoord:  newState.NextCoord,
				NextStr:    newState.NextStr,
			}
			if oldBest, ok := bestCostByState[stateKey]; ok {
				switch {
				case newCost < oldBest:
					unvisitedQueue.Update(newState, float64(newCost))
				case newCost > oldBest:
					newCost = oldBest
				default:
					unvisitedQueue.Put(newState, float64(newCost))
				}
			} else {
				unvisitedQueue.Put(newState, float64(newCost))
			}
			bestCostByState[stateKey] = newCost
		}
	}

	return results, nil
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
