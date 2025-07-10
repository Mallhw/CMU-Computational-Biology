package main
import (
    "fmt"
    "bufio"
    "os"
    "strconv"
    "strings"
    "sort"
    "slices"
)

type StringIndex map[string][]int

// Please do not remove package declarations because these are used by the autograder. If you need additional packages, then you may declare them above.

// Insert your BuildMinimizerMap() function here, along with any subroutines that you need. Please note the subroutines indicated in the problem description that are provided for you.
func BuildMinimizerMap(reads []string, k int, windowLength int) StringIndex {
	dict := make(StringIndex)
	//for loop for each read in reads
	for i := 0; i < len(reads); i++ {
		//for loop for each window in read
		for j := 0; j < len(reads[i]) - windowLength + 1; j++ {
            lexiMin := Minimizer(reads[i][j : j + windowLength], k) 
            if !slices.Contains(dict[lexiMin], i) {
                
                dict[lexiMin] = append(dict[lexiMin], i)
            }
		}
		
	}
	return dict
}
