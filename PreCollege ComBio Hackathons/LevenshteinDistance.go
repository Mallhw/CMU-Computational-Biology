package main

import (
    "fmt"
)

func main() {
	
	
}

//LongestCommonSubsequence takes as input two strings
//returns string of all matches between two strings


func LongestCommonSubsequence(str1, str2 string) string {
	dp := LCSScoreMatrix(str1, str2) 
	i := len(str1) //numRows
	j := len(str2) //numCols
    finalString := ""
	if str1[i-1] == str2[j-1] {
        finalString = string(str1[i-1])
	}
	for i > 1 || j > 1 { 
        
		if dp[i - 1][j - 1] >= dp[i - 1][j] && dp[i - 1][j - 1] >= dp[i][j - 1] { 
			i--
			j--
            
		} else if dp[i - 1][j] >= dp[i][j - 1] {
			i--
		} else {
			j--
		}
        if i == 0 || j == 0 {
            break
        }
		if str1[i-1] == str2[j-1] { //if there's a match
            finalString = string(str1[i-1]) + finalString
		}
	}
	return finalString
}



func LCSScoreMatrix()(str1, str2 string) [][]int {
	dp := make([][]int, len(str1)+1)
	for i := range dp {
		dp[i] = make([]int, len(str2)+1)
	}
	for i := 1 ; i <= len(str1); i++ {
		for j := 1; j <= len(str2); j++ {
			if str1[i - 1] == str2[j - 1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				dp[i][j] = max(dp[i-1][j], dp[i][j-1])
			}
		}
	}
	return dp
}


func GlobalScoreTable(str1, str2 string, match, mismatch, gap float64) [][]float64 {
    dp := make([][]int, len(str1)+1)
	for i := range dp {
		dp[i] = make([]int, len(str2)+1)
	}
	for i := 1 ; i <= len(str1); i++ {
		for j := 1; j <= len(str2); j++ {
			if str1[i - 1] == str2[j - 1] {
				dp[i][j] = dp[i-1][j-1] + match
			} else {
				dp[i][j] = max(dp[i-1][j] - gap, dp[i][j-1] - gap, dp[i-1][j-1] - mismatch)
			}
		}
	}
	return dp
}