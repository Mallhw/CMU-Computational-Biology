package main

import "fmt"

// AverageOutDegree takes the adjacency list of a directed network.
// It returns the average outdegree of each node in the network.
// It does not include nodes with no outgoing edges in the average.
func AverageOutDegree(adjList map[string][]string) float64 {
	count := 0
	totalNonZeroNodes := 0

	for currentRead := range adjList {
		numNeighbors := len(adjList[currentRead])

		// only count the current node if length of current slice associated with current read is > 0
		if numNeighbors > 0 {
			count += numNeighbors
			totalNonZeroNodes++ // add 1 to whatever we finally divide by
		}

	}

	if totalNonZeroNodes == 0 {
		fmt.Println("Note: you gave me a network that doesn't have any edges.")
		return 0.0
	}

	return float64(count) / float64(totalNonZeroNodes)
}

// AverageOutDegreeAllNodes takes as input the adjacency list of a directed network.
// It returns the average outdegree of each node in the network, including nodes with outdegree zero.
func AverageOutDegreeAllNodes(adjList map[string][]string) float64 {
	return 0.0
}
