package main

// KmerComposition returns the k-mer composition (all k-mer substrings) of a given genome.
func KmerComposition(genome string, k int) []string {
	n := len(genome)

	kmers := make([]string, n-k+1)

	//range over kmers and identify current substring to add to kmers
	for i := 0; i < n-k+1; i++ {
		kmers[i] = genome[i : i+k]
	}

	return kmers
}
