package main

// InitializeClusters takes a tree and returns a slice of
// pointers to the leaves of the tree.
func InitializeClusters(t Tree) []*Node {
	numNodes := len(t)

	// we don't want numNodes clusters
	// numNodes = 2*numLeaves - 1
	// (numNodes + 1)/2 = numLeaves
	numLeaves := (numNodes + 1) / 2

	clusters := make([]*Node, numLeaves)

	// make copies of all the pointers to leaf nodes
	for i := range clusters {
		clusters[i] = t[i] // t[i] is ultimately just an address, and so this is the same thing as x = y
	}

	return clusters
}
