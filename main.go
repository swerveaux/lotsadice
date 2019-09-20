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
	if len(os.Args) != 2 {
		fmt.Println("Need an argument on how many rolls of how many sided dice.   Like, to roll 100 20-sided dice, use '100d20'")
		os.Exit(1)
	}

	pair := strings.Split(os.Args[1], "d")
	if len(pair) != 2 {
		fmt.Printf("Argument does not appear to be in the right format.   It should be a number, a d, and another number, not %s\n", os.Args[1])
		os.Exit(1)
	}
	numRolls, err := strconv.Atoi(pair[0])
	if err != nil {
		fmt.Printf("Part of arg before 'd', %q, is not numeric\n", pair[0])
		os.Exit(1)
	}
	dieSize, err := strconv.Atoi(pair[1])
	if err != nil {
		fmt.Printf("Part of arg after 'd', %q, is not numeric\n", pair[1])
		os.Exit(1)
	}

	numCPUs := runtime.NumCPU()
	rollsPerCPU := numRolls / numCPUs // This will likely give us a remainder.
	extraRolls := numRolls % numCPUs  // This should be the extra number of rolls.

	totalsCh := make(chan int, numCPUs)
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < numCPUs; i++ {
		go func(index int) {
			rollsForThisCPU := rollsPerCPU
			if index < extraRolls {
				rollsForThisCPU++
			}
			var total int
			for j := 0; j < rollsForThisCPU; j++ {
				total = total + rand.Int()%dieSize + 1
			}
			totalsCh <- total
		}(i)
	}

	total := 0
	for i := 0; i < numCPUs; i++ {
		total += <-totalsCh
	}

	fmt.Printf("Rolled %d from %d rolls across %d goroutines\n", total, numRolls, numCPUs)
}
