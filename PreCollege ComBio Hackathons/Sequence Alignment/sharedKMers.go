package main

// CountSharedKmers takes two strings and an integer k. It returns the number of
// k-mers that are shared by the two strings.
func CountSharedKmers(str1, str2 string, k int) int {
	count := 0

	// form two frequency tables
	freqMap1 := FrequencyMap(str1, k)
	freqMap2 := FrequencyMap(str2, k)

	// we only need to range over the keys of one map
	for pattern, val1 := range freqMap1 {
		// does it occur in the second map?
		val2 := freqMap2[pattern]
		// this will be zero if pattern is not a key of freqMap2
		// I am interested in the minimum of val1 and val2
		count += Min2(val1, val2)
	}

	return count
}

// FrequencyMap
// Input: string text, integer k
// Output: map associating k-mers of text with their number of occurrences as substrings
func FrequencyMap(text string, k int) map[string]int {
	freq := make(map[string]int)

	// range over every k-mer of the string and add 1 to the current k-mer's value in the table

	n := len(text)

	// remember: there are n-k+1 k-mers in a string of length n
	for i := 0; i < n-k+1; i++ {
		pattern := text[i : i+k]
		// check if pattern exists as a key of map.
		// if not, create it and set it to 1.
		// if so, increase its value by 1.
		freq[pattern]++
		//if freq[pattern] doesn't exist, this creates it (0) and then adds 1 to it

		/*
			// don't need this long code
			_, exists := freq[pattern]
			if exists == false {
				freq[pattern] = 1
			} else {
				freq[pattern]++
			}
		*/
	}

	return freq
}

func Min2(x, y int) int {
	if x < y {
		return x
	}
	return y
}
