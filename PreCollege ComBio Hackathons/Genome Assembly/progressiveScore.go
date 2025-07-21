package main

import "math"

// OverlapProfileScoringMatrix computes the scoring matrix for the overlap alignment of a read against a profile matrix.
func OverlapProfileScoringMatrix(profile Profile, read string, match, mismatch, gap float64) [][]float64 {
	// initialize the scoring matrix
	numRows := len(profile) + 1
	numCols := len(read) + 1

	scoringMatrix := InitializeFloatMatrix(numRows, numCols)

	// we assume that str1 is coded along the rows of the matrix.
	// Thus, we want the 0-th column to all be zero to allow for "free rides" starting
	// at any position in the string. These values are zero by default.

	// next, set the values in the first row to negative infinity to only allow
	// prefix of str2 to match with suffix of str1.
	for j := 1; j < numCols; j++ {
		scoringMatrix[0][j] = float64(math.MinInt64)
	}

	// now I am ready to range row by row over the interior of table and apply the recurrence relation
	// for overlap alignment.
	for i := 1; i < numRows; i++ {
		for j := 1; j < numCols; j++ {
			// the up value should correspond to penalizing the gap penalty times the number of non gaps in the current column of the profile matrix
			upValue := scoringMatrix[i-1][j] - gap*float64(numberofNonGapSymbols(profile[i-1]))

			// the left value should correspond to penalizing the gap penalty times the number of gaps in the current column of the profile matrix
			leftValue := scoringMatrix[i][j-1] - gap*float64(numberofNonGapSymbols(profile[i-1]))

			// the diagonal value should correspond to the match or mismatch score
			// the current character in the read is read[j-1] and we reward it with the match score for every symbol in the profile matrix matching it; we penalize it with the mismatch score for every symbol in the profile matrix not matching it; and we penalize it with the gap score for every gap in the profile matrix of the current column.

			var diagonalWeight float64

			// range over the symbols in the profile matrix for the current column and compute the diagonal value accordingly
			for symbol, count := range profile[i-1] {
				if symbol == rune(read[j-1]) { // match!
					diagonalWeight += match * count
				} else if symbol == '-' { // gap!
					diagonalWeight -= gap * count
				} else { // mismatch!
					diagonalWeight -= mismatch * count
				}
			}

			diagValue := scoringMatrix[i-1][j-1] + diagonalWeight
			// now, we take the maximum of the three values
			scoringMatrix[i][j] = MaxFloat(upValue, leftValue, diagValue)

		}
	}

	return scoringMatrix
}

// numberofNonGapSymbols counts the number of non-gap symbols in a profile matrix column.
func numberofNonGapSymbols(column map[rune]float64) int {
	count := 0
	for symbol, freq := range column {
		if symbol != '-' && freq > 0 {
			count += int(freq)
		}
	}
	return count
}

// numberOfGapSymbols counts the number of gap symbols in a profile matrix column.
func numberOfGapSymbols(column map[rune]float64) int {
	count := 0
	for symbol, freq := range column {
		if symbol == '-' && freq > 0 {
			count += int(freq)
		}
	}
	return count
}
