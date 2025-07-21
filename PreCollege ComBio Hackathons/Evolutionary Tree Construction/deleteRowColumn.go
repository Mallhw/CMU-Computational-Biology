package main

// DeleteRowCol takes a distance matrix along with two indices.
// It returns the matrix after deleting the row and column corresponding
// to each of the indices.
// NOTE: you should assume that row < col.
func DeleteRowCol(mtx DistanceMatrix, row, col int) DistanceMatrix {
	// first, let's delete appropriate rows
	// remember that col > row, we should delete col-th row first
	mtx = append(mtx[:col], mtx[col+1:]...)
	mtx = append(mtx[:row], mtx[row+1:]...)

	//now, delete columns row and col as well. (going backwards)
	for i := len(mtx) - 1; i >= 0; i-- {
		mtx[i] = append(mtx[i][:col], mtx[i][col+1:]...)
		mtx[i] = append(mtx[i][:row], mtx[i][row+1:]...)
	}

	return mtx
}
