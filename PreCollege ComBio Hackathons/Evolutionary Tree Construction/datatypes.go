package main

//DistanceMatrix is a 2D slice of floats
type DistanceMatrix [][]float64

//Tree is a slice of pointers to nodes
type Tree []*Node

//Alignment is a multiple alignment object corresponding to a slice of strings
type Alignment []string

//Node is an object that represents a node of a tree.
//We also think of a node as a "cluster" when building a UPGMA tree.
type Node struct {
	Alignment      Alignment
	Num            int
	Age            float64
	Label          string
	Child1, Child2 *Node // if at leaf, both will be nil
}
