package main

import "fmt"

// LongestCommonSubsequence takes two strings as input.
// It returns a longest common subsequence of the two strings.
func LongestCommonSubsequence(str1, str2 string) string {
	if len(str1) == 0 || len(str2) == 0 {
		panic("Empty string given.")
	}

	//form a string

	lcs := ""

	// generate the scoring matrix
	scoringMatrix := LCSScoreMatrix(str1, str2)

	// start at sink (bottom right) and backtrack
	r := len(str1)
	c := len(str2)

	// when will we stop? At (0, 0)
	// the opposite is when we continue, which is while they are both not equal to zero
	for r != 0 && c != 0 {
		// consult our three values
		up := scoringMatrix[r-1][c]
		fmt.Println("r and c are", r, c)
		left := scoringMatrix[r][c-1]
		diag := scoringMatrix[r-1][c-1]

		// if we're at a match, add 1 to diag
		if str1[r-1] == str2[c-1] {
			diag++
		}

		// as we backtrack, every time we hit a match diagonal, prepend it to the start our string

		// which of the three values is equal to scoring matrix value at (r, c)?
		if scoringMatrix[r][c] == up {
			// move up. Subtract one from r
			r--
		} else if scoringMatrix[r][c] == left {
			// move left
			c--
		} else if scoringMatrix[r][c] == diag {
			//move diagonally
			//add current symbol to start of our growing string IF there's a match
			if str1[r-1] == str2[c-1] {
				lcs = string(str1[r-1]) + lcs
			}
			r--
			c--

		} else {
			panic("Something very wrong with table.")
		}

	}

	// return our string
	return lcs
}
