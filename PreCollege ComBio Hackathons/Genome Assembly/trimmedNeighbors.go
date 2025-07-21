package main

// GetTrimmedNeighbors takes in a string pattern (read), an adjacency list and maxK.
// It returns all n-transitivity edges in the adjList of the current read (pattern) for n <= maxK.
func GetTrimmedNeighbors(pattern string, adjList map[string][]string, maxK int) []string {
	neighbors := adjList[pattern]

	//expand the neighbor set and find the extended neighborhood
	extendedNeighbors := GetExtendedNeighbors(pattern, adjList, maxK)

	// range over the neighbors, and if the current neighbor is found in the extended neighborhood, get rid of it
	for i := len(neighbors) - 1; i >= 0; i-- {
		if Contains(extendedNeighbors, neighbors[i]) {
			// length of neighbors keeps changing
			neighbors = append(neighbors[:i], neighbors[i+1:]...)
		}
	}

	return neighbors
}

func Contains(patterns []string, element string) bool {
	for _, pattern := range patterns {
		if pattern == element {
			return true
		}
	}
	return false
}
