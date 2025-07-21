package main

import (
	"fmt"
)

// MakeOverlapNetworkMinimizers takes a slice of reads, and integer k, a minimizer dictionary, along with parameters match, mismatch, gap and a threshold.
// It returns adjacency list of reads; edges are only included
// in the overlap graph is the alignment score of the two reads is at least the threshold AND if the two reads share some k-mer minimizer according to the dictionary.
func MakeOverlapNetworkMinimizers(reads []string, minimizerDictionary StringIndex, match, mismatch, gap, threshold float64) map[string][]string {
	numReads := len(reads)

	// make the matrix
	overlapMatrix := InitializeFloatMatrix(numReads, numReads)

	// everything is zero by default, which is kind of good
	bigNegative := threshold - 10000000.0
	// this flag will tell us have I overlapped two reads before?
	// at this point in time, that answer is NO for all pairs of reads
	for r := range overlapMatrix {
		for c := range overlapMatrix[r] {
			overlapMatrix[r][c] = bigNegative
		}
	}

	counter := 0

	// engine of function: range over minimizer map, and ask:
	// (1) do two strings share the same minimizer?
	// (2) have I overlapped them already?
	for _, readIndices := range minimizerDictionary {
		if counter%100 == 0 {
			fmt.Println("Now considering element", counter, "of minimizer map.")
		}

		//readIndices is a slice of integers
		//range over this slice, and grab all possible pairs of elements that are in the slice
		for i := range readIndices {
			for j := i + 1; j < len(readIndices); j++ {
				// note: mtx[i][j] != mtx[j][i], so I may need to do two alignments
				index1 := readIndices[i]
				index2 := readIndices[j]

				read1 := reads[index1]
				read2 := reads[index2]

				//perform alignment of read1 and read2
				// IF I haven't already
				if overlapMatrix[index1][index2] == bigNegative {
					overlapMatrix[index1][index2] = ScoreOverlapAlignment(read1, read2, match, mismatch, gap)
				}
				if overlapMatrix[index2][index1] == bigNegative {
					overlapMatrix[index2][index1] = ScoreOverlapAlignment(read2, read1, match, mismatch, gap)
				}

			}
		}

		// update counter
		counter++
	}

	// overlap matrix is now made. Binarize it!
	b := BinarizeMatrix(overlapMatrix, threshold)

	return ConvertAdjacencyMatrixToList(reads, b)
}
