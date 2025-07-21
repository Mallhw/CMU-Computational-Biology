package main

import (
	"strings"
)

type Alignment []string

type Profile []map[rune]float64

func BuildContig(reads []string, match, mismatch, gap float64) string {
	if len(reads) == 0 {
		return ""
	}

	// build a huge alignment of all the reads, overlapping one at a time

	if len(reads) == 1 {
		return reads[0] // if only one read, return it as the contig
	}
	// build an alignment consisting just of reads[0]
	var a Alignment = Alignment{reads[0]}

	// for all remaining reads, align them against the current alignment
	for i := 1; i < len(reads); i++ {
		a = ProgressiveOverlapAlign(a, reads[i], match, mismatch, gap)
	}

	// trim any columns that have fewer than some number of non-gap symbols
	minNonGaps := 2
	a = TrimSparseCols(a, minNonGaps) // << new line

	// write the alignment to a file for debugging purposes
	WriteMultipleAlignment(a, "Output/contig_alignment.txt")
	return Consensus(a)

}

// TrimSparseCols removes every column that has fewer than minNonGaps
// non-gap symbols across all rows.
func TrimSparseCols(a Alignment, minNonGaps int) Alignment {
	if len(a) == 0 {
		return a
	}

	keep := make([]bool, NumColumns(a))
	for col := 0; col < NumColumns(a); col++ {
		nz := 0
		for row := 0; row < len(a); row++ {
			if a[row][col] != '-' {
				nz++
			}
		}
		keep[col] = nz >= minNonGaps
	}

	trimmed := make(Alignment, len(a))
	for row := 0; row < len(a); row++ {
		var b strings.Builder
		for col, ok := range keep {
			if ok {
				b.WriteByte(a[row][col])
			}
		}
		trimmed[row] = b.String()
	}
	return trimmed
}

// Consensus returns a consensus string for an alignment.
// At each column it chooses the non-gap rune with the highest count.
// If there is a tie or no non-gap symbols, it writes 'N'.
func Consensus(a Alignment) string {
	var consensus strings.Builder

	numCols := NumColumns(a)
	for col := 0; col < numCols; col++ {
		counts := make(map[rune]int) // tallies for non-gap symbols
		nonGapTotal := 0             // total symbols that aren’t '-'

		for row := 0; row < len(a); row++ {
			if col >= len(a[row]) {
				continue // this row is shorter than the current column
			}
			r := rune(a[row][col])
			if r == '-' {
				continue // ignore gaps when voting
			}
			counts[r]++
			nonGapTotal++
		}

		minCov := 3 // say, require ≥ 3 reads
		if nonGapTotal < minCov {
			consensus.WriteRune('N')
			continue
		}

		// pick the rune with the largest count (> 50 % of non-gaps)
		var (
			consensusChar rune = 'N'
			maxCount           = 0
		)
		for r, c := range counts {
			if c > maxCount {
				maxCount = c
				consensusChar = r
			}
		}

		// require strict majority among non-gaps
		if maxCount*2 <= nonGapTotal {
			consensusChar = 'N'
		}

		consensus.WriteRune(consensusChar)
	}
	return consensus.String()
}

// NumColumns returns the number of columns in the alignment.
func NumColumns(a Alignment) int {
	AssertRectangular(a)
	if len(a) == 0 {
		return 0
	}
	// return the length of the first read in the alignment
	return len(a[0])
}

// AssertRectangular checks if the alignment is rectangular (all rows have the same length).
func AssertRectangular(a Alignment) {
	if len(a) == 0 {
		return // empty alignment is considered rectangular
	}
	firstLength := len(a[0])
	for _, row := range a {
		if len(row) != firstLength {
			panic("Alignment is not rectangular")
		}
	}
}
