package main
import (
    "fmt"
    "bufio"
    "os"
    "strconv"
)

// Please do not remove package declarations because these are used by the autograder. If you need additional packages, then you may declare them above.

// Insert your Minimizer() function here, along with any subroutines that you need. 
func Minimizer(text string, k int) string {
    kmer := text[0 : k]
	for i:=0; i < len(text) - k + 1; i ++ {
		if text[i : i + k] < kmer {
			kmer = text[i : i + k]
		}
	}
	return kmer
}