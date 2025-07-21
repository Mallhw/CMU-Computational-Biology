package main

// DeleteClusters
// Input: slice of node pointers, row and col indices
// Note: col > row
// Output: updated slice in which we delete clusters[row] and clusters[col]
func DeleteClusters(clusters []*Node, row, col int) []*Node {
	// why I stressed that col > row is that it will allow you to delete clusters[col] first
	//first, delete clusters[col]
	clusters = append(clusters[:col], clusters[col+1:]...)

	//next, delete clusters[row]
	clusters = append(clusters[:row], clusters[row+1:]...)

	return clusters
}
