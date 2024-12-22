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
	InputFile              string `arg:"positional,required" help:"input file"`
	NumIntermediateKeypads int    `arg:"-n,required" help:"number of intermediate keypads"`
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
	solutionsCache := make(map[string][]string)

	for r := range lib.NumPadKeyByRune {
		solutionLength, err := solveKeyPad(string([]rune{r}), solutionsCache, allMaps, 2)
		if err != nil {
			log.Panic(err)
		}
		log.Printf("num pad string: %s", string([]rune{r}))
		log.Printf("solutionLength: %d", solutionLength)
	}

	for r := range lib.NumPadKeyByRune {
		solutionLength, err := solveKeyPad(string([]rune{r}), solutionsCache, allMaps, 3)
		if err != nil {
			log.Panic(err)
		}
		log.Printf("num pad string: %s", string([]rune{r}))
		log.Printf("solutionLength: %d", solutionLength)
	}

	for _, numPadCode := range numPadCodes {
		numPadString := string(lo.Map(numPadCode, func(key, _ int) rune {
			return allMaps.NumPad.RunesByNumPadKey[key]
		}))
		value, err := solveAndMultiply(numPadString, solutionsCache, allMaps, args.NumIntermediateKeypads)
		if err != nil {
			log.Panic(err)
		}
		log.Printf("num pad string: %s", numPadString)
		log.Printf("value: %d", value)
		total += value
	}

	log.Printf("total: %d", total)
}

func solveAndMultiply(target string, solutionsCache map[string][]string, allMaps AllMaps, numIntermediateKeypads int) (int64, error) {
	solutionLength, err := solveKeyPad(target, solutionsCache, allMaps, numIntermediateKeypads)
	if err != nil {
		return -1, err
	}

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

	return value * int64(solutionLength), nil
}

func solveKeyPad(target string, solutionsCache map[string][]string, allMaps AllMaps, numIntermediateKeypads int) (int, error) {
	solutionsKeypad, err := doDijkstra(target, allMaps.NumPad.RunesByNumPadKey, allMaps.NumPad.Layout, allMaps.NumPad.RevLayout[lib.NumPadKeyA])
	if err != nil {
		return -1, err
	}
	log.Printf("number of solutionsKeypad: %d", len(solutionsKeypad))

	solutionLength, err := solveActionPad(solutionsCache, allMaps, numIntermediateKeypads, solutionsKeypad)
	if err != nil {
		return -1, err
	}
	return solutionLength, nil
}

func solveActionPad(solutionsCache map[string][]string, allMaps AllMaps, numIntermediateKeypads int, solutionsKeypad []string) (int, error) {
	minLengthKeypad := lo.Min(lo.Map(solutionsKeypad, func(s string, _ int) int {
		return len(s)
	}))
	filteredSolutionsKeypad := lo.Filter(solutionsKeypad, func(s string, _ int) bool {
		return len(s) == minLengthKeypad
	})
	log.Printf("number of solutionsAction after filtering: %d", len(filteredSolutionsKeypad))

	minLengthPrev := minLengthKeypad
	filteredSolutionsPrev := filteredSolutionsKeypad
	sizeOfCache := len(solutionsCache)
	for iKeyPad := numIntermediateKeypads; iKeyPad > 0; iKeyPad-- {
		solutionsNext := set.New[string](0)
		for _, solutionPrev := range filteredSolutionsPrev {
			subSolutionsNext, err := doActionDijkstra(solutionPrev, solutionsCache, allMaps.ActionPad.RunesByAction, allMaps.ActionPad.Layout, allMaps.ActionPad.RevLayout[lib.Press])
			if err != nil {
				return -1, err
			}

			subSolutionsNextWithCache := set.New[string](len(subSolutionsNext))
			subSolutionsNextWithCache.InsertSlice(subSolutionsNext)
			subSolutionsNextWithCache.InsertSlice(lo.ValueOr(solutionsCache, solutionPrev, []string{}))
			subSolutionsNextWithCacheSlice := subSolutionsNextWithCache.Slice()
			minLengthNext := lo.Min(lo.Map(subSolutionsNextWithCacheSlice, func(s string, _ int) int {
				return len(s)
			}))
			filteredSubSolutionsNext := lo.Filter(subSolutionsNextWithCacheSlice, func(s string, _ int) bool {
				return len(s) == minLengthNext
			})

			solutionsCache[solutionPrev] = filteredSubSolutionsNext
			newSizeOfCache := len(solutionsCache)
			if newSizeOfCache > sizeOfCache {
				sizeOfCache = newSizeOfCache
				log.Printf("round [%d]: size of solutionsCache: %d", iKeyPad, sizeOfCache)
			}
			solutionsNext.InsertSlice(filteredSubSolutionsNext)
		}
		log.Printf("round [%d]: number of solutionsNext: %d", iKeyPad, solutionsNext.Size())

		solutionsNextSlice := solutionsNext.Slice()
		minLengthNext := lo.Min(lo.Map(solutionsNextSlice, func(s string, _ int) int {
			return len(s)
		}))
		filteredSolutionsNext := lo.Filter(solutionsNextSlice, func(s string, _ int) bool {
			return len(s) == minLengthNext
		})
		log.Printf("round [%d]: number of solutionsNext after filtering: %d", iKeyPad, len(filteredSolutionsNext))

		minLengthPrev = minLengthNext
		filteredSolutionsPrev = filteredSolutionsNext
	}

	return minLengthPrev, nil
}

func doActionDijkstra(targetString string, solutionsCache map[string][]string, runesByKey map[int]rune, nextLayout map[lib.Coord]int, initialCoord lib.Coord) ([]string, error) {
	// origTargetString := targetString
	var preSolutions [][]string
	prefixFound := true
	for prefixFound && len(targetString) > 0 {
		prefixFound = false
		for prefix, solutions := range solutionsCache {
			if strings.HasPrefix(targetString, prefix) {
				preSolutions = append(preSolutions, solutions)
				targetString = targetString[len(prefix):]
				prefixFound = true
			}
		}
	}

	allSolutions := preSolutions
	if len(targetString) > 0 {
		solutions, err := doDijkstra(targetString, runesByKey, nextLayout, initialCoord)
		if err != nil {
			return nil, err
		}

		allSolutions = append(allSolutions, solutions)
	}

	product := lo.Reduce(allSolutions, func(agg []string, item []string, _ int) []string {
		return cartesianConcat(agg, item)
	}, []string{""})

	// if len(allSolutions) > 1 {
	// 	solutions, err := doDijkstra(origTargetString, runesByKey, nextLayout, initialCoord)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	//
	// 	slices.Sort(product)
	// 	slices.Sort(solutions)
	// 	if !reflect.DeepEqual(solutions, product) {
	// 		return nil, fmt.Errorf("internal error: answers are not equal: %v != %v", solutions, product)
	// 	}
	// }

	return product, nil
}

func cartesianConcat(slice1, slice2 []string) []string {
	var result []string
	for _, s1 := range slice1 {
		for _, s2 := range slice2 {
			result = append(result, s1+s2)
		}
	}
	return result
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
