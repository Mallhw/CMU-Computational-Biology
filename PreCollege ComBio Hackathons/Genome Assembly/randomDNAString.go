package main

import "math/rand"

// GenerateRandomGenome takes a parameter length and returns
// a random DNA string of this length where each nucleotide has equal probability.
func GenerateRandomGenome(length int) string {
	// generate an array of random symbols
	symbols := make([]byte, length)
	for i := range symbols {
		symbols[i] = RandomDNASymbol()
	}

	/*
		// don't do this -- string concatenations are slow
		genome := ""
		for i := 0; i < length; i++ {
			genome += string(RandomDNASymbol())
		}
	*/

	return string(symbols)
}

func RandomDNASymbol() byte {
	number := rand.Intn(4)
	switch number {
	case 0:
		return 'A'
	case 1:
		return 'C'
	case 2:
		return 'G'
	case 3:
		return 'T'
	}
	panic("Error: something really wrong with RNG.")
}
