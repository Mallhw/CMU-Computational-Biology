package main

// LCSLength takes two strings as input. It returns the length of a longest common
// subsequence of the two strings.
func LCSLength(str1, str2 string) int {
	if len(str1) == 0 || len(str2) == 0 {
		return 0
	}

	// use subroutine to form the scoring matrix
	mtx := LCSScoreMatrix(str1, str2)

	//numRows := len(str1)+1
	//numCols := len(str2)+1

	return mtx[len(str1)][len(str2)]
}

// LCSScoreMatrix
// Input: two strings(
// Output: the scoring matrix for the symbol matching problem (i.e., LCS length) according to dynamic programming
func LCSScoreMatrix(str1, str2 string) [][]int {
	if len(str1) == 0 || len(str2) == 0 {
		panic("Error: zero string length.")
	}

	numRows := len(str1) + 1
	numCols := len(str2) + 1

	scoringMatrix := InitializeMatrix(numRows, numCols)

	// fill in values according to dynamic programming
	// the first row and column are equal to zero by default :)
	// start indexing at (1, 1)
	for row := 1; row < numRows; row++ {
		for col := 1; col < numCols; col++ {
			// apply the recurrence
			// consult 3 matrix values
			up := scoringMatrix[row-1][col]
			left := scoringMatrix[row][col-1]
			diag := scoringMatrix[row-1][col-1]

			// check if there's a match between str1[row-1] and str2[col-1]
			if str1[row-1] == str2[col-1] {
				//add 1 to diag value
				diag++
			}

			// take max of all three values
			scoringMatrix[row][col] = MaxInts(up, left, diag)
		}
	}

	return scoringMatrix
}
