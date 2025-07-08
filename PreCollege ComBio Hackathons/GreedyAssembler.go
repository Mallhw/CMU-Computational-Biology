package main

import (
	"fmt"
)

func main(reads string[]) {
	reads := []string{"CC", "CC", "CC", "CC", "CC"}
	fmt.Println(GreedyAssembler(reads))
}



// Insert your GreedyAssembler() function here, along with any subroutines that you need. Please note the subroutines indicated in the problem description that are provided for you.
func GreedyAssembler(reads []string) string {
	ans := make([]string, 0)
	ans = string[0]
	i := 0
	// right
	while i < len(reads) {
		fmt.Println(reads[i])
		i++
		if ans[:2] == reads[i][:2] {
			ans += reads[i][2:]
			i = 0
		}
	}
	// left 
	while i < len(reads) {
		fmt.Println(reads[i])
		i++
		if ans[-2:] == reads[i][-2:] {
			ans = reads[i][:-2] + ans
			i = 0
		}
	}
	return ans
}