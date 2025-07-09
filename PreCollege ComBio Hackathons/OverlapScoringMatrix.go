func OverlapScoringMatrix(reads []string, match, mismatch, gap float64) [][]float64 {
	table := make([][]float64, len(reads))
	for i := 0; i < len(reads); i++ {
		table[i] = make([]float64, len(reads))
	}
	for i := 0; i < len(reads); i++ {
		for j := 0; j < len(reads); j++ {
			table[i][j] = ScoreOverlapAlignment(reads[i], reads[j], match, mismatch, gap)
            if i == j {
				table[i][j] = 0
			}
		}
	}
	return table
}
