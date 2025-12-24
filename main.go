// This is a reasonably contrived example of how to span a CPU-bound
// problem over all the available cores on a machine.   There was a news
// story about a truck that spilled 216,000 dice and someone on my
// work slack complained that the 'roll' command maxed out at 100 rolls,
// and someone else worried that trying to do 216000d20 would "break
// something", so I figured I'd whip that up in go.

// Since it's CPU bound, spinning up a goroutine per roll will likely just
// cause lots of thrashing.   It makes sense to split up the load as
// evenly as we can across cores, start up a goroutine per core, then have
// them each take their portion of the work and return the result for
// aggregation at the end.

package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

var repeatTimes int
var sortType string

func main() {
	flag.IntVar(&repeatTimes, "repeat", 1, "Number of times to repeat")
	flag.StringVar(&sortType, "sort", "sum", "Sort type. Options: [sum, count]")
	flag.Parse()
	fmt.Printf("Repeating %d times\n", repeatTimes)
	fmt.Printf("Also got %v\n", flag.Args())
	resultsTable := make(map[int]int)

	// Check length of os.Args.   First arg should be the actual command, the second is the one argument we accept.
	if len(flag.Args()) != 1 {
		fmt.Println("Need an argument on how many rolls of how many sided dice.   Like, to roll 100 20-sided dice, use '100d20'")
		os.Exit(1)
	}

	// Split on 'd'.
	pair := strings.Split(flag.Args()[0], "d")
	// If that didn't split into exactly two items, then it wasn't formatted right.
	if len(pair) != 2 {
		fmt.Printf("Argument does not appear to be in the right format.   It should be a number, a d, and another number, not %s\n", os.Args[1])
		os.Exit(1)
	}

	numRolls, err := strconv.Atoi(pair[0])
	if err != nil {
		fmt.Printf("Part of arg before 'd', %q, is not numeric\n", pair[0])
		os.Exit(1)
	}

	// Try to convert the second part to the number of sides on each die.
	dieSize, err := strconv.Atoi(pair[1])
	if err != nil {
		fmt.Printf("Part of arg after 'd', %q, is not numeric\n", pair[1])
		os.Exit(1)
	}

	// Starting a timer so we can display execution time when done.
	for i := 0; i < repeatTimes; i++ {
		total := performRolls(numRolls, dieSize)
		if _, exists := resultsTable[total]; exists {
			resultsTable[total] += 1
		} else {
			resultsTable[total] = 1
		}
		// Try to convert the first part to the number of times to roll the dice.
	}

	sortMap(resultsTable, sortType)
}

func performRolls(numRolls, dieSize int) int {
	total := rollDice(numRolls, dieSize)
	return total
}

func sortMap(resultsTable map[int]int, sortType string) {
	keys := make([]int, 0, 0)
	switch sortType {
	case "sum":
		for key := range resultsTable {
			keys = append(keys, key)
		}
		sort.Ints(keys)
		fmt.Println("Sum,Occurrences")
		for _, key := range keys {
			fmt.Printf("%d,%d\n", key, resultsTable[key])
		}
	default:
		fmt.Printf("Unsupported sort type %s\n", sortType)
	}
}
