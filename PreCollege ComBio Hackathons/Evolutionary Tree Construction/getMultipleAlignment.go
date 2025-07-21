package main

//GetMultipleAlignment takes a Tree object as input.
//It returns the alignment associated with the tree, which is the
//multiple alignment labeling the root node.
func GetMultipleAlignment(t Tree) Alignment {

	return t[len(t)-1].Alignment
}