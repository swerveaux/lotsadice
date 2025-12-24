package main

import (
	"math/rand"
	"runtime"
	"time"
)

func rollDice(numRolls, dieSize int) int {
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

	return total
}
