package main

//ProgressiveAlignmentScoreTable takes two multiple alignments as well as a collection
//of Clustal scoring parameters. It returns a 2D matrix corresponding to
//the Clustal dynamic programming table for combining the two alignments heuristically
//into a single multiple alignment.
func ProgressiveAlignmentScoreTable(align1 Alignment, align2 Alignment,
	match float64, mismatch float64, gap float64, supergap float64) [][]float64 {
	
	return [][]float64{}
}