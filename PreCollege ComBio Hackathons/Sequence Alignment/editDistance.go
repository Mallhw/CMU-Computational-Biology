package main

// EditDistance takes two strings as input. It returns the Levenshtein distance
// between the two strings; that is, the minimum number of substitutions, insertions, and deletions
// needed to transform one string into the other.
func EditDistance(str1, str2 string) int {
	if len(str1) == 0 || len(str2) == 0 {
		panic("Empty strings given.")
	}

	// make our edit distance table
	scoringMatrix := EditTable(str1, str2)

	// grab the bottom right value

	return scoringMatrix[len(str1)][len(str2)]
}

// EditTable
// Input: two strings
// Output: the dynamic programming matrix associated with edit distance for the two input strings.
func EditTable(str1, str2 string) [][]int {
	if len(str1) == 0 || len(str2) == 0 {
		panic("Empty strings given.")
	}

	numRows := len(str1) + 1
	numCols := len(str2) + 1

	scoringMatrix := InitializeMatrix(numRows, numCols)

	// set the values in the first row and column.
	// let's start with the 0-th column.
	for r := range scoringMatrix {
		scoringMatrix[r][0] = r
	}

	for c := range scoringMatrix[0] {
		scoringMatrix[0][c] = c
	}

	// range over interior of the matrix and set values
	for row := 1; row < numRows; row++ {
		for col := 1; col < numCols; col++ {
			// now we're ready to apply the recurrence relation
			up := scoringMatrix[row-1][col] + 1
			left := scoringMatrix[row][col-1] + 1
			diag := scoringMatrix[row-1][col-1]

			// if corresponding symbols mismatch, add one to diag
			if str1[row-1] != str2[col-1] {
				diag++
			}

			// finally, set scoring matrix at (row, col)
			scoringMatrix[row][col] = Min(up, left, diag)

		}
	}

	return scoringMatrix
}
