package main

//ALL PENALTIES POSITIVE

// ScoreOverlapAlignment takes two strings along with match, mismatch, and gap penalties.
// It returns the maximum score of an overlap alignment in which str1 is overlapped with str2.
// Assume we are overlapping a suffix of str1 with a prefix of str2.
func ScoreOverlapAlignment(str1, str2 string, match, mismatch, gap float64) float64 {
	if len(str1) == 0 || len(str2) == 0 {
		panic("String has empty length.")
	}

	scoreTable := OverlapScoreTable(str1, str2, match, mismatch, gap)

	// get the max in this row
	return MaxArrayFloat(scoreTable[len(scoreTable)-1])
}

func OverlapScoreTable(str1, str2 string, match, mismatch, gap float64) [][]float64 {
	numRows := len(str1) + 1
	numCols := len(str2) + 1

	scoreTable := InitializeFloatMatrix(numRows, numCols)

	// remember: one column is done -- that's the first column

	// let's set everybody in the first row
	for c := 1; c < numCols; c++ {
		scoreTable[0][c] = float64(c) * (-gap)
	}

	// now we can range over score table interior and set values
	for r := 1; r < numRows; r++ {
		for c := 1; c < numCols; c++ {
			//apply the recurrence relation
			// this is our normal recurrence relation from global alignment
			up := scoreTable[r-1][c] - gap
			left := scoreTable[r][c-1] - gap
			diag := scoreTable[r-1][c-1]

			// diag will vary based on whether we have a match or a mismatch

			if str1[r-1] == str2[c-1] {
				// match
				diag += match
			} else {
				// mismatch
				diag -= mismatch
			}

			// now we just set scoring matrix value at (r,c) based on the recurrence
			scoreTable[r][c] = MaxFloat(up, left, diag)
		}
	}
	return scoreTable
}
