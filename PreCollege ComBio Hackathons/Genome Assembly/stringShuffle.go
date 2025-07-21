package main

import "math/rand"

// ShuffleStrings takes a collection of strings patterns as input.
// It returns a random shuffle of the strings.
func ShuffleStrings(patterns []string) []string {
	n := len(patterns)

	perm := rand.Perm(n)
	// this will give us a permutation of (0, 1, ..., n-1)

	patterns2 := make([]string, n) // we will copy our strings into a new slice

	for i := range patterns {
		// e.g., say permutation is (3, 1, 0, 2, 4)
		//using perm[i] as the index in patterns
		index := perm[i]
		patterns2[i] = patterns[index]
	}

	return patterns2
}
