func ScoreOverlapAlignment(str1, str2 string, match, mismatch, gap float64) float64 {
	dp := OverlapScoreTable(str1, str2, match, mismatch, gap)
	i := len(str1)
	j := len(str2)
	// Find the Max value on the edges
    maxVal := dp[i][j]
	for k := 0; k < len(str1); k++ {
		if dp[len(str2) - 1][k] > maxVal {
			j = k
			i = len(str1)
			maxVal = dp[len(str2) - 1][k]
		}
	}

	for k := 0; k < len(str2); k++ {
		if dp[k][len(str1) - 1] > maxVal {
			i = k
			j = len(str2)
			maxVal = dp[k][len(str1) - 1]
		}
	}

	return float64(maxVal)
}

func OverlapScoreTable(str1, str2 string, match, mismatch, gap float64) [][]float64 {
    dp := make([][]float64, len(str1)+1)
	for i := range dp {
		dp[i] = make([]float64, len(str2)+1)
	}
    for i := 1; i <= len(str1); i++ {
        dp[i][0] = 0
    }
    for i := 1; i <= len(str2); i++ {
        dp[0][i] = 0
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