package main

// GetExtendedNeighbors takes in a pattern (read), the overlap graph and maxK.
// It returns the extendedNeighbors list. For each neighbor *n* in this list,
// distance between n and pattern is between 2 to maxK.
func GetExtendedNeighbors(pattern string, adjList map[string][]string, maxK int) []string {

	if maxK < 2 {
		panic("Error: maxK is too small.")
	}

	// keep track of nodes visited
	// each string has an on/off light
	visited := make(map[string]bool)

	// turn on the light for pattern
	visited[pattern] = true

	//first, jump to everybody that is one away from pattern
	currentNodes := adjList[pattern]

	// now, expand currentNodes 1 step at a time, flipping on the light as you go
	for d := 2; d <= maxK; d++ {
		currentNodes = ExpandFrontier(currentNodes, adjList, visited)
	}

	// at the end of this, I want to process current nodes
	// I know which nodes I have visited as part of the process because they are the keys of the map
	finalList := make([]string, 0, len(visited)-1)
	for node := range visited {
		if node != pattern {
			//append node to a list
			finalList = append(finalList, node)
		}
	}

	return finalList
}

// ExpandFrontier
// Input: currentnodes representing the current neighbors of a string that we have visited at distance d, an adjacency list of our network, and a map of strings to bools representing if a node has been visited
// Output: slice of nodes corresponding to one additional step (distance d+1)
func ExpandFrontier(currentNodes []string, adjList map[string][]string, visited map[string]bool) []string {
	nextNodes := make([]string, 0)

	// pass all the work to a subroutine that works on each individual node :)
	for _, frontierNode := range currentNodes {
		nextNodes = ExpandNode(frontierNode, adjList, visited, nextNodes)
	}

	return nextNodes
}

// ExpandNode
// Input: node, adjacency list, visited map of strings to bools indicating which nodes have been visited, next slice
// Output: updated next slice with newly discovered neighbors from the current node
func ExpandNode(node string, adjList map[string][]string, visited map[string]bool, nextNodes []string) []string {

	// range over all of the neighbors you can find of node
	for _, neighbor := range adjList[node] {
		if !(visited[neighbor]) {
			nextNodes = append(nextNodes, neighbor)

			//flip on the light
			visited[neighbor] = true
		}

	}

	return nextNodes
}
