package main

//ExciseSpikeProtein grabs the DNA sequence of the spike protein in a coronavirus genome.
func ExciseSpikeProtein(genome string) string {
	beginning := 21000
	end := 26000
	minLength := 3000
	leftFlankingSequence := "MFVF"
	rightFlankingSequence := "LHYT"

	spikeProtein := FindViralGene(genome, beginning, end, minLength, leftFlankingSequence, rightFlankingSequence)

	return spikeProtein
}

func FindViralGene(genome string, beginning, end, minLength int, leftFlankingSequence, rightFlankingSequence string) string {
	// there is an issue if end - beginning isn't longer than minLength
	if end-beginning < minLength {
		panic("Error: Invalid parameters given to FindRegion.")
	}
	for i := beginning; i < end-minLength; i++ {
		if Translate(DNAToRNA(genome[i:i+3*len(leftFlankingSequence)]), 0) == leftFlankingSequence {
			//hit found on left, so try the right, maximizing length of match
			for j := end - 3*len(rightFlankingSequence) - 1; j > i+minLength-3*len(rightFlankingSequence); j-- {
				if Translate(DNAToRNA(genome[j:j+3*len(rightFlankingSequence)]), 0) == rightFlankingSequence {
					// success!
					return genome[i : j+3*len(rightFlankingSequence)-1]
				}
			}
		}
	}

	return ""
}
