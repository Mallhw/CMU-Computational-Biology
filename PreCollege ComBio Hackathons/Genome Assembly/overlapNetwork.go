package main

// MakeOverlapNetwork() takes a slice of reads with match, mismatch, gap and a threshold.
// It returns adjacency lists of reads; edges are only included
// in the overlap graph is the alignment score of the two reads is at least the threshold.
func MakeOverlapNetwork(reads []string, match, mismatch, gap, threshold float64) map[string][]string {
	//Initialize an adjacency list to represent the overlap graph.

	// make the overlap matrix
	overlapMatrix := OverlapScoringMatrix(reads, match, mismatch, gap)

	// binarize it
	binarizedMatrix := BinarizeMatrix(overlapMatrix, threshold)

	//convert it to an adjacency list
	return ConvertAdjacencyMatrixToList(reads, binarizedMatrix)

}

// ConvertAdjacencyMatrixToList
// Input: a collection of strings and a binary matrix representing an adjacency matrix
// Output: an adjacency list corresponding to the network represented by the adjacency matrix
func ConvertAdjacencyMatrixToList(reads []string, b [][]int) map[string][]string {
	adjacencyList := make(map[string]([]string))

	// range over the entire matrix, and any time you see a 1, add the appropriate edge to the adjacency list
	for r := range b {
		for c := range b[r] {
			if b[r][c] == 1 {
				// we are currently entering the reads that are connected from reads[r]
				// we found that reads[r] should be connected to reads[c]
				currentRead := reads[r]
				adjacencyList[currentRead] = append(adjacencyList[currentRead], reads[c])
				// this will create the slice if currentRead is not a key of the map; otherwise, it appends to it
			}
		}
	}

	return adjacencyList
}
