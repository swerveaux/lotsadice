# lotsadice
A Go program to simulate rolling a bunch of arbitrarily-sided dice. Something I wrote after a truck spilled 216,000 dice on a highway 
and someone at work complained that the 'roll' command couldn't simulate that many rolls. It's pretty simple, spins up a goroutine per CPU
and divides up the rolls between them, collecting their results and accumulating them at the end.

## Running
It has one required argument which is N1dN2 where N1 is the number of rolls and N2 is the number of sides on the die.
e.g, to roll a six-sided die ten times, use 10d6.   To roll a twenty-sided die a million times, do 1000000d20.

