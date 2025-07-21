package main

// BinarizeMatrix takes a matrix of values and a threshold.
// It binarizes the matrix according to the threshold.
// If entries across main diagonal are both above threshold, only retain the bigger one.
func BinarizeMatrix(mtx [][]float64, threshold float64) [][]int {
	b := InitializeMatrix(mtx)

	//range over matrix and set values to 1 if they are above the threshold.
	//If entries across main diagonal are both above threshold, only retain the bigger one.
	//no need to take any action if not above threshold, since these values are 0 by default.
	for i := range mtx {
		for j := range mtx[i] {
			//Only consider filling in if mtx[i][j] is greater/equal to threshold.
			if mtx[i][j] >= threshold {
				if mtx[i][j] > mtx[j][i] {
					b[i][j] = 1
					//Edge case mtx[i][j] == mtx[j][i]
				} else if mtx[i][j] == mtx[j][i] && i < j { //or: ... && i > j
					b[i][j] = 1
				}
			}
		}
	}

	return b
}
