package main

import (
	"fmt"
)

// Please do not remove package declarations because these are used by the autograder. If you need additional packages, then you may declare them above.

// Alignment is an array of two strings corresponding to the top and bottom
// rows of an alignment of two strings.
type Alignment [2]string

func main() {
	table := GlobalScoreTable("GAGA", "GAT", 1, 1, 2)
	for i := 0; i < len(table); i++ {
		for j := 0; j < len(table[0]); j++ {
			fmt.Println(table[i][j])
		}
		fmt.Println("")
	}
	GlobalAlignment("GAGA", "GAT", 1, 1, 2)
}

// Insert your GlobalAlignment() function here, along with any subroutines that you need. Please note the subroutines indicated in the problem description that are provided for you.
func GlobalAlignment(str1, str2 string, match, mismatch, gap float64) Alignment {
	dp := GlobalScoreTable(str1, str2, match, mismatch, gap)
	var alignment Alignment
	alignment[0] = ""
	alignment[1] = ""
	i := len(str1)
	j := len(str2)
	for i != 0 || j != 0 {
		if i != 0 && j != 0 {
			score := -mismatch
			if str1[i-1] == str2[j-1] {
				score = match
			}
			if dp[i][j] == dp[i-1][j-1]+score {
				alignment[0] = string(str1[i-1]) + alignment[0]
				alignment[1] = string(str2[j-1]) + alignment[1]
				i--
				j--
			} 
		}
		if i != 0 && dp[i][j] == dp[i-1][j]-gap {
			alignment[0] = string(str1[i-1]) + alignment[0]
			alignment[1] = "-" + alignment[1]
			i--
		} else if j != 0 && dp[i][j] == dp[i][j-1]-gap {
			alignment[0] = "-" + alignment[0]
			alignment[1] = string(str2[j-1]) + alignment[1]
			j--
		}
	}
	return alignment
}

func GlobalScoreTable(str1, str2 string, match, mismatch, gap float64) [][]float64 {
	dp := make([][]float64, len(str1)+1)
	for i := range dp {
		dp[i] = make([]float64, len(str2)+1)
	}
	for i := 1; i <= len(str1); i++ {
		dp[i][0] = float64(i) * -gap
	}
	for i := 1; i <= len(str2); i++ {
		dp[0][i] = float64(i) * -gap
	}
	for i := 1; i <= len(str1); i++ {
		for j := 1; j <= len(str2); j++ {
			if str1[i-1] == str2[j-1] {
				dp[i][j] = dp[i-1][j-1] + match
			} else {
				dp[i][j] = max(dp[i-1][j]-gap, dp[i][j-1]-gap, dp[i-1][j-1]-mismatch)
			}
		}
	}
	return dp
}
