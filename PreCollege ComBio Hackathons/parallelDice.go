package main

import (
	"fmt"
	"math/rand"
	"runtime"
)

func main() {
	fmt.Println("Parallel and concurrent programming.")

	fmt.Println("This computer has", runtime.NumCPU(), "cores available.")

	/*
		n := 40
		p1 := go Perm(1, n / 2)
		p2 := go Perm(n / 2, n)

		fmt.Println("n! is ", p1*p2)
		fmt.Println("Finished")
	*/
	// both parties hae to be presetn for a message to be exchanged

	fmt.Println(CrapsHouseEdgeTimes(10))

	fmt.Println("Finished")

}

func ComputeCrapsHouseEdge(numTrials, numProcs int) float64 {
	wins := 0
	c := make(chan int)

	for i := 0; i < numProcs; i++ {
		go TotalWinOneProc(numTrials/numProcs+numTrials%numProcs, c)
	}

	for i := 0; i < numProcs; i++ {
		wins += <-c
	}
	return float64(wins) / float64(numTrials)
}

func CrapsHouseEdgeTimes(numTrials int) float64 {
	numTrial := 10000000

	start := time.Now()
	numProcs := runtime.NumCPU()
	edges := ComputeCrapsHouseEdge(numTrial, numProcs)
	elapsed := time.Since(start)
	log.Printf("Serial took %s", elapsed)

	fmt.Println(edges)
	return edges
}

func TotalWinOneProc(numTrials int, c chan int) {
	wins := 0
	for i := 0; i < numTrials; i++ {
		if PlayCrapsOnce() {
			wins++
		} else {
			wins--
		}
	}
	c <- wins
}

func PlayCrapsOnce() bool {
	roll := SumDice(2)
	if roll == 7 || roll == 11 {
		return true
	} else if roll == 2 || roll == 3 || roll == 12 {
		return false
	} else {
		for {
			nRoll := SumDice(2)
			if nRoll == roll {
				return true
			} else if nRoll == 7 {
				return false
			}
		}
	}
}

func SumDice(numDice int) int {
	ans := 0
	for i := 0; i < numDice; i++ {
		ans += (rand.Intn(6) + 1)
	}
	return ans
}
