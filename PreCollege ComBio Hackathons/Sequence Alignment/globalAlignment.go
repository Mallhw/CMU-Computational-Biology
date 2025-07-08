package main

type Alignment [2]string

// GlobalAlignment takes two strings, along with match, mismatch, and gap scores.
// It returns a maximum score global alignment of the strings corresponding to these penalties.
func GlobalAlignment(str1, str2 string, match, mismatch, gap float64) Alignment {
	var alignment Alignment

	// grab the table
	scoringMatrix := GlobalScoreTable(str1, str2, match, mismatch, gap)

	// get the bottom right element
	r := len(str1)
	c := len(str2)

	//backtrack from bottom right, forming alignment as we go
	// strings start off as empty
	// backtrack first to the first row or column
	for r != 0 && c != 0 {
		if scoringMatrix[r][c] == scoringMatrix[r-1][c]-gap {
			//UP
			// take a symbol from string 1 against a gap in the lower column
			alignment[0] = string(str1[r-1]) + alignment[0]
			alignment[1] = "-" + alignment[1]
			r--
		} else if scoringMatrix[r][c] == scoringMatrix[r][c-1]-gap {
			// LEFT

			//take a gap symbol and align against a symbol from string 2
			alignment[0] = "-" + alignment[0]
			alignment[1] = string(str2[c-1]) + alignment[1]
			c--
		} else {
			// only works if scoring table is correct :)
			// take a symbol from both strings
			alignment[0] = string(str1[r-1]) + alignment[0]
			alignment[1] = string(str2[c-1]) + alignment[1]
			r--
			c--

			// this is bad! I should technically check that the current value is equal to diag value + match or - mismatch depending on case
		}
	}

	// form anything left in the alignment

	return alignment
}
