package main

import "fmt"

// OverlapScoringMatrix takes a collection of reads along with alignment penalties.
// It returns a matrix in which mtx[i][j] is the overlap alignment score of
// reads[i] with reads[j].
func OverlapScoringMatrix(reads []string, match, mismatch, gap float64) [][]float64 {
	numReads := len(reads)
	mtx := InitializeFloatMatrix(numReads, numReads)

	for i := range mtx {
		if i%10 == 0 {
			fmt.Println("Currently making row", i, "of overlap matrix.")
		}
		for j := range mtx[i] {
			if i != j { // avoid loops, so don't set main diagonal values
				mtx[i][j] = ScoreOverlapAlignment(reads[i], reads[j], match, mismatch, gap)
			}
		}
	}
	return mtx
}
