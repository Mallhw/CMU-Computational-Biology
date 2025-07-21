package main

// StringIndex is a type that will map a minimizer string to its list of indices
// in a read set corresponding to reads with this minimizer.
type StringIndex map[string][]int

// BuildMinimizerMap takes a collection of reads, integer k and integer windowLength.
// It returns a map of k-mers to the indices of the reads in the list having this k-mer minimizer.
func BuildMinimizerMap(reads []string, k int, windowLength int) StringIndex {
	dict := make(StringIndex)

	// range over reads
	for i, read := range reads {
		// range over all windows of current read and take minimizers
		n := len(read)
		// how many windows? n-windowLength+1
		for j := 0; j < n-windowLength+1; j++ {
			currentWindow := read[j : j+windowLength]
			//grab the minimizer of the window
			m := Minimizer(currentWindow, k)
			// now that I have the minimizer, append i (read index) to the minimizer's integer slice
			// but only add it if I haven't seen it before
			if !IsElementOf(dict[m], i) {
				dict[m] = append(dict[m], i)
			}
		}
	}

	return dict
}

// IsElementOf
// Input: slice of integers and an integer value
// Output: true if givenValue appears in slice, and false otherwises
func IsElementOf(a []int, givenValue int) bool {
	// range over a and if you find givenValue in a, return true

	for _, val := range a {
		if val == givenValue {
			return true
		}
	}

	// we survive all tests, and so the default is false because givenValue is not in a
	return false
}
