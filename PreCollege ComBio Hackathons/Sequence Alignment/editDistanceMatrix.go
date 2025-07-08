package main

// EditDistanceMatrix takes as input a slice of strings patterns.
// It returns a matrix whose (i,j)th entry is the edit distance between
// the i-th and j-th strings in patterns.
func EditDistanceMatrix(patterns []string) [][]int {
	numPatterns := len(patterns)
	mtx := InitializeMatrix(numPatterns, numPatterns)

	// range over the entire matrix and set edit distances
	for r := range mtx {
		for c := r + 1; c < len(mtx[r]); c++ {
			mtx[r][c] = EditDistance(patterns[r], patterns[c])
			// set same value across the main diagonal
			mtx[c][r] = mtx[r][c]
		}
	}

	return mtx
}
