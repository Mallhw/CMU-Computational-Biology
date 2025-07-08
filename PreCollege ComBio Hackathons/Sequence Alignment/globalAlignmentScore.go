package main

// GlobalScoreTable takes two strings and alignment penalties. It returns a 2-D array
// holding dynamic programming scores for global alignment with these penalties.
func GlobalScoreTable(str1, str2 string, match, mismatch, gap float64) [][]float64 {
	if len(str1) == 0 || len(str2) == 0 {
		panic("Empty string given.")
	}

	numRows := len(str1) + 1
	numCols := len(str2) + 1

	scoringMatrix := InitializeFloatMatrix(numRows, numCols)

	// set values of table
	// first, set 0-th row and column
	for r := 1; r < numRows; r++ {
		scoringMatrix[r][0] = float64(r) * (-gap)
	}

	for c := 1; c < numCols; c++ {
		scoringMatrix[0][c] = float64(c) * (-gap)
	}

	// now, range over the interior
	for r := 1; r < numRows; r++ {
		for c := 1; c < numCols; c++ {
			up := scoringMatrix[r-1][c] - gap
			left := scoringMatrix[r][c-1] - gap
			diag := scoringMatrix[r-1][c-1]

			// diag will vary based on whether we have a match or a mismatch

			if str1[r-1] == str2[c-1] {
				// match
				diag += match
			} else {
				// mismatch
				diag -= mismatch
			}

			// now we just set scoring matrix value at (r,c) based on the recurrence
			scoringMatrix[r][c] = MaxFloat(up, left, diag)
		}
	}

	return scoringMatrix
}
