package main

import (
	"fmt"
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

	ChannelBasics()
	ParallelFactorial()

	fmt.Println("Finished")

}

func ParallelFactorial() {
	n := 40
	c := make(chan int)
	go Perm(1, n / 2, c)
	go Perm(n / 2, n, c)

	p1 := <- c
	p2 := <- c

	fmt.Println("n! is ", p1*p2)
	fmt.Println("Finished")

}

func ChannelBasics() {
	c := make(chan string)

	// This channel is "synchronous", when you send or recieve a message into or out of the channel, the rest of the function "blocks" or refuses to run until the other party is preseent
	
	go SayHi(c)
	// function blocks here

	msg := <- c
	
	fmt.Println(msg)
}

func SayHi(c chan string) {
	c <- "Hello"
}

func Perm(k, n int, c chan int) {
	p := 1

	for i := k; i <= n; i++ {
		p *= i
	}

	c <- p
}

func printFactorial(n int) {
	p := 1

	for i := 1; i <= n; i++ {
		fmt.Println(p)
		p *= i		
	}
}

func Factorial(n int) int {
	p := 1

	for i := 1; i <= n; i++ {
		p *= i
	}

	return p
}


