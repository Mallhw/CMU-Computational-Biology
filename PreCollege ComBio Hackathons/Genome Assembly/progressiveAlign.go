package main

import (
	"math"
	"strings"
)

// ProgressiveOverlapAlign aligns a new read against an existing alignment
// using the progressive alignment method.
// It performs an overlap alignment, in which we overlap a suffix of the existing alignment
// with a prefix of the new read, and we take the alignment achieving the maximum score
// over all such alignments.
// It computes the score of the alignment using the provided scoring parameters, where
// we score the current symbol against all elements of the indicated column in the alignment.
// match pairs receive the match reward; mismatch pairs receive the mismatch penalty;
// gaps receive the gap penalty.
func ProgressiveOverlapAlign(a Alignment, read string, match, mismatch, gap float64) Alignment {
	//------------------------------------------------------------------(A)
	// Freeze all but the right-most len(read) columns
	window := len(read)
	frozen, aSuf := SplitFrozenPrefix(a, window) // <— new helper
	//-----------------------------------------------------------------------

	// ------------ original code, but using aSuf instead of a --------------
	profile := MakeProfile(aSuf)
	scoringMatrix := OverlapProfileScoringMatrix(profile, read,
		match, mismatch, gap)

	maxScore := float64(math.MinInt64)
	maxCol := -1
	for j := 1; j < len(scoringMatrix[0]); j++ {
		if scoringMatrix[len(scoringMatrix)-1][j] > maxScore {
			maxScore = scoringMatrix[len(scoringMatrix)-1][j]
			maxCol = j
		}
	}

	// rows now []byte instead of string
	alignBytes := make([][]byte, len(aSuf)+1)
	for i := range alignBytes {
		alignBytes[i] = nil
	}

	// Traceback (unchanged except we index aSuf)
	row := len(scoringMatrix) - 1
	col := maxCol
	for row > 0 || col > 0 {
		if row == 0 { // LEFT
			alignBytes[len(aSuf)] = prepend(alignBytes[len(aSuf)], read[col-1])
			for i := range aSuf {
				alignBytes[i] = prepend(alignBytes[i], '-')
			}
			col--
			continue
		}
		if col == 0 { // UP
			for i := range aSuf {
				alignBytes[i] = prepend(alignBytes[i], aSuf[i][row-1])
			}
			alignBytes[len(aSuf)] = prepend(alignBytes[len(aSuf)], '-')
			row--
			continue
		}

		up := scoringMatrix[row-1][col] -
			gap*float64(numberofNonGapSymbols(profile[row-1]))
		left := scoringMatrix[row][col-1] -
			gap*float64(numberofNonGapSymbols(profile[row-1]))

		switch {
		case scoringMatrix[row][col] == up:
			for i := range aSuf {
				alignBytes[i] = prepend(alignBytes[i], aSuf[i][row-1])
			}
			alignBytes[len(aSuf)] = prepend(alignBytes[len(aSuf)], '-')
			row--
		case scoringMatrix[row][col] == left:
			for i := range aSuf {
				alignBytes[i] = prepend(alignBytes[i], '-')
			}
			alignBytes[len(aSuf)] = prepend(alignBytes[len(aSuf)], read[col-1])
			col--
		default:
			for i := range aSuf {
				alignBytes[i] = prepend(alignBytes[i], aSuf[i][row-1])
			}
			alignBytes[len(aSuf)] = prepend(alignBytes[len(aSuf)], read[col-1])
			row--
			col--
		}
	}

	// unchanged suffix loop
	for j := maxCol; j < len(read); j++ {
		alignBytes[len(aSuf)] = append(alignBytes[len(aSuf)], read[j])
		for i := range aSuf {
			alignBytes[i] = append(alignBytes[i], '-')
		}
	}

	// convert []byte rows → string
	suffixAligned := make(Alignment, len(alignBytes))
	for i := range alignBytes {
		suffixAligned[i] = string(alignBytes[i])
	}

	//------------------------------------------------------------------(B)
	// Stitch the frozen prefix back in front
	out := make(Alignment, len(suffixAligned))
	for i := 0; i < len(a); i++ { // existing rows
		out[i] = frozen[i] + suffixAligned[i]
	}
	prefixGaps := strings.Repeat("-", len(frozen[0])) // same width for all
	out[len(a)] = prefixGaps + suffixAligned[len(a)]  // new-read row
	//------------------------------------------------------------------(C)

	return out
}

// SplitFrozenPrefix returns two things:
//
//  1. frozen[i]  — the left-hand part of row i that we will **never** touch
//  2. suffix     — the alignment consisting of just the last <window> columns
func SplitFrozenPrefix(a Alignment, window int) (frozen []string, suffix Alignment) {
	nCols := NumColumns(a)
	if nCols <= window { // nothing to freeze
		frozen = make([]string, len(a))
		suffix = a
		return
	}

	cut := nCols - window
	frozen = make([]string, len(a))
	suffix = make(Alignment, len(a))
	for i := range a {
		frozen[i] = a[i][:cut] // left part (fixed)
		suffix[i] = a[i][cut:] // right part (we’ll overlap on this)
	}
	return
}

// helper: prepend one byte to a slice (amortised O(len), but avoids strings)
func prepend(b []byte, c byte) []byte {
	b = append(b, 0)          // grow by 1
	copy(b[1:], b[:len(b)-1]) // shift right
	b[0] = c
	return b
}

// MakeSuffixProfile returns a profile of the last w columns of a.
func MakeSuffixProfile(a Alignment, w int) (Profile, int) {
	totalCols := NumColumns(a)
	if w >= totalCols {
		return MakeProfile(a), 0 // whole alignment fits in the window
	}
	start := totalCols - w
	suffix := make(Alignment, len(a))
	for i := range a {
		suffix[i] = a[i][start:] // slice, no copy
	}
	return MakeProfile(suffix), start // start = # frozen columns
}

// MakeProfile creates a profile matrix from the given alignment.
// If a symbol is a gap symbol such that no other non-gap symbol exists to its right in the same row of the alignment, then it does not contribute to the profile matrix.
func MakeProfile(a Alignment) Profile {
	AssertRectangular(a)

	// create a profile matrix with the same number of columns as the alignment
	numCols := NumColumns(a)
	profile := make(Profile, numCols)

	// fill the profile matrix
	for col := 0; col < numCols; col++ {
		profile[col] = make(map[rune]float64)
		for row := 0; row < len(a); row++ {
			char := rune(a[row][col])
			if char != '-' { // ignore gap symbols
				profile[col][char]++
			}
			// if the character is a gap, then it only counts if there is a non-gap character to its right in the same row (i.e., it's an internal gap in its string)
			if char == '-' && !LastGapSymbol(a[row], col) {
				profile[col]['-']++
			}
		}
	}

	return profile
}

// LastGapSymbol checks if the character at the given column in the row is a gap
// and if it is the last character in the row or if there are no non-gap characters to its right.
func LastGapSymbol(row string, col int) bool {
	if col >= len(row) || row[col] != '-' {
		return false // not a gap or out of bounds
	}
	// check if there are any non-gap characters to the right
	for i := col + 1; i < len(row); i++ {
		if row[i] != '-' {
			return false // found a non-gap character
		}
	}
	return true // no non-gap characters to the right
}
