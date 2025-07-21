package main

// FindMinElement takes a distance matrix as input.
// It returns a pair (row, col, val) where (row, col) corresponds to the minimum
// value of the matrix, and val is the minimum value.
// NOTE: you should have that row < col.
func FindMinElement(mtx DistanceMatrix) (int, int, float64) {
	if len(mtx) <= 1 || len(mtx[0]) <= 1 {
		panic("One row or one column!")
	}

	// can now assume that matrix is at least 2 x 2
	row := 0
	col := 1
	minVal := mtx[row][col]

	// range over matrix, and see if we can do better than minVal.
	for i := 0; i < len(mtx)-1; i++ {
		// start column ranging at i + 1
		for j := i + 1; j < len(mtx[i]); j++ {
			// do we have a winner?
			if mtx[i][j] < minVal {
				// update all three variables
				minVal = mtx[i][j]
				row = i
				col = j
				// col will still always be > row.
			}
		}
	}
	return row, col, minVal
}
