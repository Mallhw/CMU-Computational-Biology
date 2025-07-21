package main

// GenerateContigs takes an adjacency list for an overlap network.
// It returns a collection of contigs, or a slice containing lists of
// strings, where each list is the reads of the contig in order.
func GenerateContigs(adjList map[string][]string) [][]string {

	//Get all suffixes
	allSuffixes := []string{}
	allPrefixes := []string{}
	for prefix, suffixList := range adjList {
		for i := 0; i < len(suffixList); i++ {
			allSuffixes = append(allSuffixes, suffixList[i])
		}

		allPrefixes = append(allPrefixes, prefix)
	}

	inDegreesAllNodes := make(map[string]int)
	outDegreesAllNodes := make(map[string]int)
	numPaths := 0
	graphStarts := make([]string, 0)
	for node, _ := range adjList {
		// found is a boolean. True if prefix is found in allSuffixes

		outDegreesAllNodes[node] = len(adjList[node])
		inDegreesAllNodes[node] = countOccurrences(allSuffixes, node)

		if outDegreesAllNodes[node] != inDegreesAllNodes[node] {
			if outDegreesAllNodes[node] > 0 {
				graphStarts = append(graphStarts, node)
				for q := 0; q < outDegreesAllNodes[node]; q++ {
					numPaths = numPaths + 1
				}
			}

		}
		if outDegreesAllNodes[node] == inDegreesAllNodes[node] {
			if outDegreesAllNodes[node] > 1 {
				graphStarts = append(graphStarts, node)
				for q := 0; q < outDegreesAllNodes[node]; q++ {
					numPaths = numPaths + 1
				}
			}

		}
	}

	//make our collection of maximal paths
	maximalPaths := make([][]string, numPaths)
	for i := 0; i < numPaths; i++ {
		maximalPaths[i] = make([]string, 0)
	}

	counter := 0
	OutDegree := 0
	InDegree := 0

	currentNode := ""

	// Iterate over starting nodes
	for _, startNode := range graphStarts {

		OutDegree = len(adjList[startNode])                 // # suffixes this prefix node has
		InDegree = countOccurrences(allSuffixes, startNode) // # times this node appears as a suffix

		// For each path beginning at this starting node, extend it until
		// the in degree or the out degree of the node added is not 1
		for b := 0; b < OutDegree; b++ {

			// Add the start node to this maximal path
			maximalPaths[counter] = append(maximalPaths[counter], startNode)

			currentNode = adjList[startNode][b]

			// WHILE loop. EXTEND PATH until ...
			z := 0
			for z == 0 {
				// Add current node to path
				maximalPaths[counter] = append(maximalPaths[counter], currentNode)

				// The next node can have outDegree of either 0 or >0.
				if len(adjList[currentNode]) == 0 {
					// There is no next node
					z = 1
					counter = counter + 1
				} else if len(adjList[currentNode]) == 1 {
					// There is 1 next node
					currentNode = adjList[currentNode][0]

					InDegree = countOccurrences(allSuffixes, currentNode)
					if InDegree != 1 {
						maximalPaths[counter] = append(maximalPaths[counter], currentNode)
						z = 1
						// We finished extending this maximal path. We can update counter for the
						// next maximal path.
						counter = counter + 1
					}

				} else {
					// There is more than 1 next node (path branches)
					//OutDegree = len(Graph[currentNode])
					//InDegree = countOccurrences(allSuffixes,currentNode)
					z = 1
					// We finished extending this maximal path. We can update counter for the
					// next maximal path.
					counter = counter + 1
				}

			}
		}
	}

	return maximalPaths
}

// MT: https://golangcode.com/check-if-element-exists-in-slice/
// Find takes a slice and looks for an element in it. If found it will
// return it's key, otherwise it will return -1 and a bool of false.
func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func countOccurrences(slice []string, val string) int {
	counter := 0
	for _, item := range slice {
		if item == val {
			counter = counter + 1
		}
	}
	return counter
}
