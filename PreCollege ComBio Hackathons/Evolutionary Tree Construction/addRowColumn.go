package main

// AddRowCol takes a distance matrix Given a DistanceMatrix, a slice of current clusters,
// and a row/col index (NOTE: col > row).
// It returns the matrix corresponding to "gluing" together clusters[row] and clusters[col]
// forming a new row/col of the matrix for the new cluster, computing
// distances to other elements of the matrix weighted according to the sizes
// of clusters[row] and clusters[col].
func AddRowCol(row, col, clusterSize1, clusterSize2 int, mtx DistanceMatrix) DistanceMatrix {
	numRows := len(mtx)

	newRow := make([]float64, numRows+1)

	// set values of newRow
	// all values are 0.0 by default
	// final value should be 0.0
	// newRow[row] and newRow[col] should be 0.0 too
	for j := 0; j < len(newRow)-1; j++ {
		if j != row && j != col {
			// set the value of newRow[j]
			newRow[j] = (float64(clusterSize1)*mtx[row][j] + float64(clusterSize2)*mtx[col][j]) / float64(clusterSize1+clusterSize2)
		}
	}

	// append newRow to mtx
	mtx = append(mtx, newRow)

	//add last column to the matrix according to symmetry
	for i := 0; i < numRows; i++ {
		mtx[i] = append(mtx[i], newRow[i])
	}

	return mtx
}
