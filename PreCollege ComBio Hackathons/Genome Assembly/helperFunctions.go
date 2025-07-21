package main

import "sort"

type PairwiseAlignment [2]string

// InitializeMatrix takes a 2-D matrix of float64 values as well as an integer.
// It returns an integer matrix of the same dimensions.
func InitializeMatrix(mtx [][]float64) [][]int {
	b := make([][]int, len(mtx))
	for i := range mtx {
		b[i] = make([]int, len(mtx[i]))
	}

	return b
}

// CopyGraph returns a copy of the input graph.
func CopyGraph(graph map[string][]string) map[string][]string {
	newGraph := make(map[string][]string)

	for key, value := range graph {
		newArr := make([]string, len(value))
		n := copy(newArr, value)
		if n != len(value) {
			panic("Something wrong happend when copying graph.")
		}
		newGraph[key] = newArr
	}
	return newGraph
}

// MtxEq takes two matrices and returns true if they are equal
// and false otherwise.
func MtxEq(mtx1, mtx2 [][]int) bool {
	if len(mtx1) != len(mtx2) {
		return false
	}
	for i := range mtx1 {
		if len(mtx1[i]) != len(mtx2[i]) {
			return false
		}
		for j := range mtx1[i] {
			if mtx1[i][j] != mtx2[i][j] {
				return false
			}
		}
	}
	return true
}

// GraphEq takes 2 graphs graph1 and graph2.
// It returns if graph1 and graph2 has the same keys and values.
func GraphEq(graph1, graph2 map[string][]string) bool {
	for key1, value1 := range graph1 {
		value2, found := graph2[key1]
		if !found {
			return false
		}
		if !StrSliceEq(value1, value2) {
			return false
		}
	}
	for key2, value2 := range graph2 {
		value1, found := graph1[key2]
		if !found {
			return false
		}
		if !StrSliceEq(value1, value2) {
			return false
		}
	}
	return true
}

// StrSliceEq takes 2 string slices slice1 and slice2.
// It returns if slice1 == slice2.
func StrSliceEq(slice1, slice2 []string) bool {
	sort.Strings(slice1)
	sort.Strings(slice2)
	if len(slice1) != len(slice2) {
		return false
	}
	for i := 0; i < len(slice1); i++ {
		if slice1[i] != slice2[i] {
			return false
		}
	}
	return true
}

// SumLength takes a collection of strings and returns their total length.
func SumLength(patterns []string) int {
	c := 0

	for i := range patterns {
		c += len(patterns[i])
	}

	return c
}

// GlobalAlignment takes two strings, along with match, mismatch, and gap scores.
// It returns a maximum score global alignment of the strings corresponding to these penalties.
func GlobalAlignment(str1, str2 string, match, mismatch, gap float64) PairwiseAlignment {
	var alignment PairwiseAlignment // variable declaration of array will have two empty strings

	// grab the table

	scoringMatrix := GlobalScoreTable(str1, str2, match, mismatch, gap)

	// backtrack to find the best global alignment

	r := len(str1)
	c := len(str2)

	// backtrack to the zero-th row or column
	for r > 0 && c > 0 {
		// we have four cases
		// first, match
		if scoringMatrix[r][c] == scoringMatrix[r-1][c]-gap {
			// UP
			// symbol from string 1 against gap
			alignment[0] = string(str1[r-1]) + alignment[0]
			alignment[1] = "-" + alignment[1]
			r--
		} else if scoringMatrix[r][c] == scoringMatrix[r][c-1]-gap {
			// LEFT
			// gap against symbol of string 2
			alignment[0] = "-" + alignment[0]
			alignment[1] = string(str2[c-1]) + alignment[1]
			c--
		} else if str1[r-1] != str2[c-1] && scoringMatrix[r][c] == scoringMatrix[r-1][c-1]-mismatch {
			// mismatch case
			alignment[0] = string(str1[r-1]) + alignment[0]
			alignment[1] = string(str2[c-1]) + alignment[1]
			// travel diagonally up to the left
			r--
			c--
		} else if str1[r-1] == str2[c-1] && scoringMatrix[r][c] == scoringMatrix[r-1][c-1]+match {
			// I need to add elements to the alignment
			alignment[0] = string(str1[r-1]) + alignment[0]
			alignment[1] = string(str2[c-1]) + alignment[1]
			// travel diagonally up to the left
			r--
			c--
		} else {
			panic("I don't know why I am here.")
		}
	}

	// we may be in the zero-th row or column
	// backtrack until I hit the origin
	for r > 0 {
		// move up until the origin
		// align symbol from string 1 against gap
		alignment[0] = string(str1[r-1]) + alignment[0]
		alignment[1] = "-" + alignment[1]
		r--
	}

	for c > 0 {
		// move left until origin
		// gap against symbol of string 2
		alignment[0] = "-" + alignment[0]
		alignment[1] = string(str2[c-1]) + alignment[1]
		c--
	}

	return alignment
}

// GlobalScoreTable takes two strings and alignment penalties. It returns a 2-D array
// holding dynamic programming scores for global alignment with these penalties.
func GlobalScoreTable(str1, str2 string, match, mismatch, gap float64) [][]float64 {
	if len(str1) == 0 || len(str2) == 0 {
		panic("Empty string given.")
	}

	numRows := len(str1) + 1
	numCols := len(str2) + 1

	scoringMatrix := InitializeFloatMatrix(numRows, numCols)

	// set values of table
	// first, set 0-th row and column
	for r := 1; r < numRows; r++ {
		scoringMatrix[r][0] = float64(r) * (-gap)
	}

	for c := 1; c < numCols; c++ {
		scoringMatrix[0][c] = float64(c) * (-gap)
	}

	// now, range over the interior
	for r := 1; r < numRows; r++ {
		for c := 1; c < numCols; c++ {
			up := scoringMatrix[r-1][c] - gap
			left := scoringMatrix[r][c-1] - gap
			diag := scoringMatrix[r-1][c-1]

			// diag will vary based on whether we have a match or a mismatch

			if str1[r-1] == str2[c-1] {
				// match
				diag += match
			} else {
				// mismatch
				diag -= mismatch
			}

			// now we just set scoring matrix value at (r,c) based on the recurrence
			scoringMatrix[r][c] = MaxFloat(up, left, diag)
		}
	}

	return scoringMatrix
}

// MaxArrayFloat takes a slice of integers as input and returns the maximum value in the slice.
func MaxArrayFloat(a []float64) float64 {
	if len(a) == 0 {
		panic("Error: array given to MaxArray has zero length.")
	}
	m := a[0]
	for i := 1; i < len(a); i++ {
		if a[i] > m {
			m = a[i]
		}
	}
	return m
}

// MaxFloat takes an arbitary number of integers as input and returns their maximum.
func MaxFloat(nums ...float64) float64 {
	m := 0.0
	// nums gets converted to an array
	for i, val := range nums {
		if val > m || i == 0 {
			// update m
			m = val
		}
	}
	return m
}

// StringSliceEquals compares if 2 string slices contain the same amount of the same patterns.
func StringSliceEquals(patterns1, patterns2 []string) bool {
	//sort first
	sort.Strings(patterns1)
	sort.Strings(patterns2)

	if len(patterns1) != len(patterns2) {
		return false
	}
	for i := range patterns1 {
		if patterns1[i] != patterns2[i] {
			return false
		}
	}
	return true
}

// InitializeFloatMatrix takes a number of rows and columns as input.
// It returns a matrix of zeroes as floats with appropriate dimensions.
func InitializeFloatMatrix(numRows, numCols int) [][]float64 {
	scoreTable := make([][]float64, numRows)
	for i := range scoreTable {
		scoreTable[i] = make([]float64, numCols)
	}
	return scoreTable
}
