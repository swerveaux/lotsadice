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
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func main() {
	// Starting a timer so we can display execution time when done.
	start := time.Now()

	// Check length of os.Args.   First arg should be the actual command, the second is the one argument we accept.
	if len(os.Args) != 2 {
		fmt.Println("Need an argument on how many rolls of how many sided dice.   Like, to roll 100 20-sided dice, use '100d20'")
		os.Exit(1)
	}

	// Split on 'd'.
	pair := strings.Split(os.Args[1], "d")
	// If that didn't split into exactly two items, then it wasn't formatted right.
	if len(pair) != 2 {
		fmt.Printf("Argument does not appear to be in the right format.   It should be a number, a d, and another number, not %s\n", os.Args[1])
		os.Exit(1)
	}

	// Try to convert the first part to the number of times to roll the dice.
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

	// We're going to try to spread the workload across the processors available,
	// so get the number of CPUs.
	numCPUs := runtime.NumCPU()

	// It's unlikely that the number of rolls will line up evenly with the
	// number of CPUs, and we can't roll a die a fractional number of times,
	// so get the whole number and then the remainder, which we'll use to
	// add an extra roll to some of the goroutines.
	rollsPerCPU := numRolls / numCPUs // This will likely give us a remainder.
	extraRolls := numRolls % numCPUs  // This should be the extra number of rolls.

	// Set up a buffered channel to get the incoming total from each goroutine
	// after it's added up its rolls.   We want buffered so the goroutines don't
	// block when they're done.   It shouldn't really matter that much, but if
	// we decide to make a timeout or something later, it'll keep us from
	// stranding goroutines.
	totalsCh := make(chan int, numCPUs)

	// Seed the random number generator.
	rand.Seed(time.Now().UnixNano())

	// Start a goroutine per CPU responsible for its number of rolls.
	for i := 0; i < numCPUs; i++ {
		go func(index int) {
			rollsForThisCPU := rollsPerCPU
			// Look at the remainder to decide if there is an extra roll we should
			// take care of to get to the total.
			if index < extraRolls {
				rollsForThisCPU++
			}
			var total int
			// sum the total from all the rolls this goroutine is responsible for.
			for j := 0; j < rollsForThisCPU; j++ {
				total = total + rand.Int()%dieSize + 1
			}
			totalsCh <- total
		}(i)
	}

	total := 0
	// Collect the totals per goroutine into one single sum.
	for i := 0; i < numCPUs; i++ {
		total += <-totalsCh
	}

	fmt.Printf("Rolled %d from %d rolls across %d goroutines in %v\n", total, numRolls, numCPUs, time.Since(start))
}
