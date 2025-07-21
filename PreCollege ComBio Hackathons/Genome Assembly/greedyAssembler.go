package main

// GreedyAssembler takes a collection of strings and returns a genome whose
// k-mer composition is these strings. It makes the following assumptions.
// 1. "Perfect coverage" -- every k-mer is detected
// 2. No errors in reads
// 3. Every read has equal length (k)
// 4. DNA is single-stranded
func GreedyAssembler(reads []string) string {
	// greedy algorithm: look for whatever helps me the most (overlap of k-1 symbols).
	if len(reads) == 0 {
		panic("Error: No reads given to GenomeAssembler!")
	}
	//create a copy of reads into reads2
	reads2 := make([]string, len(reads))
	copy(reads2, reads)
	//create a string to hold the genome, starting with the first read
	genome := reads2[0]

	// identify k
	k := len(reads2[0])

	//remove the first read from reads2
	reads2 = reads2[1:]
	//while we can find something to add to genome, keep adding to genome
	foundLeft := true
	foundRight := true
	for foundRight {
		//find the next read to add to genome on right side
		nextRead, nextReadIndex := FindNextReadRight(genome, reads2, k)
		if nextReadIndex == -1 {
			//if we can't find anything to add, we are done
			foundRight = false
		} else { // we found something
			//add nextRead to genome
			genome += nextRead[len(nextRead)-1:]
			//remove nextRead from reads2
			reads2 = append(reads2[:nextReadIndex], reads2[nextReadIndex+1:]...)
		}
	}
	for foundLeft {
		//find the next read to add to genome on left side
		nextRead, nextReadIndex := FindNextReadLeft(genome, reads2, k)
		if nextReadIndex == -1 {
			//if we can't find anything to add, we are done
			foundLeft = false
		} else { // we found something
			//add nextRead to genome
			genome = nextRead[:1] + genome
			//remove nextRead from reads2
			reads2 = append(reads2[:nextReadIndex], reads2[nextReadIndex+1:]...)
		}
	}
	return genome
}

// FindNextReadRight takes as input a genome, a collection of reads, and an integer k corresponding to k-mer length.
// It returns the next read to add to the genome (to the right), as well as the index of this read in the collection of reads.
func FindNextReadRight(genome string, reads []string, k int) (string, int) {
	//find the next read to add to genome
	var nextRead string
	nextReadIndex := -1                //keep track of the index of the next read
	suffix := genome[len(genome)-k+1:] // suffix of length k-1

	//does any read have suffix as a prefix?
	for i, read := range reads {
		if read[:k-1] == suffix {
			//if so, set nextRead to this read
			nextRead = read
			nextReadIndex = i
			return nextRead, nextReadIndex
		}
	}

	return nextRead, nextReadIndex // default: return "", -1
}

// FindNextReadLeft takes as input a genome, a collection of reads, and an integer k corresponding to k-mer length.
// It returns the next read to add to the genome (to the left), as well as the index of this read in the collection of reads.
func FindNextReadLeft(genome string, reads []string, k int) (string, int) {
	//find the next read to add to genome
	var nextRead string
	nextReadIndex := -1    //keep track of the index of the next read
	prefix := genome[:k-1] // prefix of length k-1

	//does any read have prefix as a suffix?
	for i, read := range reads {
		if read[1:] == prefix {
			//if so, set nextRead to this read
			nextRead = read
			nextReadIndex = i
			return nextRead, nextReadIndex
		}
	}

	return nextRead, nextReadIndex // default: return "", -1
}
