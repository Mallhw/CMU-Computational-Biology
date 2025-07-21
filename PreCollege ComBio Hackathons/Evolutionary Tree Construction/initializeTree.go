package main

import "strconv"

// InitializeTree takes a slice of n species names as input.
// It returns a rooted binary tree with with 2n-1 total nodes, where
// the leaves are the first n and have the associated species names.
func InitializeTree(speciesNames []string) Tree {
	var t Tree // []*Node, initially nil
	numLeaves := len(speciesNames)

	t = make(Tree, 2*numLeaves-1)
	// all of these node pointers are nil

	// we don't want this long term

	for i := range t {
		// we need a Node in order to point to it
		var vx Node

		// let's label it with the species name if it's a leaf

		if i < numLeaves {
			vx.Label = speciesNames[i]
		} else {
			// give it an unspecific name
			vx.Label = "Ancestor species " + strconv.Itoa(i)
		}

		vx.Num = i

		//final and most important step: link t[i] to the node that we created
		t[i] = &vx

	}

	return t
}
