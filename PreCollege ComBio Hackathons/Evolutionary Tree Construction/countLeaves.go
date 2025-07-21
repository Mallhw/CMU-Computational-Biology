package main

// CountLeaves takes a non-nil pointer to a Node object and returns
// the number of leaves in the tree rooted at the node. It returns 1 at a leaf.
func CountLeaves(v *Node) int {
	//base case: we are at a leaf, return 1
	if v.Child1 == nil && v.Child2 == nil {
		return 1
	}
	// what if only one is nil?
	if v.Child1 == nil {
		return CountLeaves(v.Child2) // we know this is OK since we made it here
	}
	if v.Child2 == nil {
		return CountLeaves(v.Child1) // same reasoning
	}
	//if I make it here, the node has two children that are non-nil
	return CountLeaves(v.Child1) + CountLeaves(v.Child2)
}
