package main

// UPGMA takes a distance matrix and a collection of species names as input.
// It returns a Tree (an array of nodes) resulting from applying
// UPGMA to this dataset.
func UPGMA(mtx DistanceMatrix, speciesNames []string) Tree {
	// it would be nice to assert
	// 1. mtx is square
	// 2. numrows of mtx is len(speciesNames)
	t := InitializeTree(speciesNames)

	clusters := InitializeClusters(t)

	// now for the engine of UPGMA ... apply the steps of the algorithm numLeaves - 1 times
	numLeaves := len(speciesNames)
	for p := numLeaves; p < 2*numLeaves-1; p++ {
		row, col, val := FindMinElement(mtx)
		// remember: col > row

		// set the age of current node
		t[p].Age = val / 2.0

		//connect to the kids
		t[p].Child1 = clusters[row] // copying over the "key"
		t[p].Child2 = clusters[col]

		//how many leaves are in each cluster?
		clusterSize1 := CountLeaves(t[p].Child1)
		clusterSize2 := CountLeaves(t[p].Child2)

		mtx = AddRowCol(row, col, clusterSize1, clusterSize2, mtx)

		clusters = append(clusters, t[p]) // add new node to end of clusters

		mtx = DeleteRowCol(mtx, row, col)

		// delete clusters[row] and clusters[col]
		clusters = DeleteClusters(clusters, row, col)

	}

	return t
}
