package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func GraphToDot(filename string, adjList map[string][]string) {
	outFile, err := os.Create(filename)
	if err != nil {
		panic("Sorry, couldn't create file!")
	}
	fmt.Fprintln(outFile, "digraph g {")
	for node, neigh := range adjList {
		for _, n := range neigh {
			fmt.Fprintln(outFile, node, " -> ", n)
		}
	}
	fmt.Fprintln(outFile, "}")
	outFile.Close()
}

func GraphToFastg(filename string, adjList map[string][]string) {
	outFile, err := os.Create(filename)
	if err != nil {
		panic("Sorry, couldn't create file!")
	}
	idxMap := make(map[string]int, 0)

	i := 0
	for key := range adjList {
		idxMap[key] = i
		i++
	}

	for node, neigh := range adjList {
		fmt.Fprint(outFile, ">", "NODE_", idxMap[node], "_length_", len(node), "_cov_", 1.0)
		if len(neigh) > 0 {
			fmt.Fprint(outFile, ":")
		}
		for i, n := range neigh {
			fmt.Fprint(outFile, "NODE_", idxMap[n], "_length_", len(n), "_cov_", 1.0)
			if i != len(neigh)-1 {
				fmt.Fprint(outFile, ",")
			}
		}
		fmt.Fprintln(outFile, ";")
		fmt.Fprintln(outFile, node)
	}
	outFile.Close()
}

// ReadStringsFromFASTA takes a file name of a FASTA file and
// It reads in all strings in this file if there are multiple FASTA headers
// and places the strings into a slice.
func ReadStringsFromFASTA(filename string) []string {
	file, err := os.Open(filename)

	if err != nil {
		// error in opening file
		panic("Error: something went wrong with file open (probably you gave wrong filename).")
	}

	scanner := bufio.NewScanner(file) // think of this as a "reader bot"
	patterns := make([]string, 0)

	currentString := ""

	// go for as long as the reader bot can still see text
	for scanner.Scan() {
		currentLine := scanner.Text() // grabs one line of text and returns a string
		if len(currentLine) > 0 && currentLine[0] != '>' {
			// append the current line to our growing read
			currentString += currentLine
		} else { // we are at a header
			// the current read is complete! :) append it
			if len(currentString) > 0 {
				patterns = append(patterns, currentString)
			}
			currentString = ""
		}
	}

	//at end, add current string
	if len(currentString) > 0 {
		patterns = append(patterns, currentString)
	}

	// we have read everything in
	if scanner.Err() != nil {
		panic("Error: issue in scanning process.")
	}

	file.Close()

	return patterns
}

func ReadGenomeFromFASTA(filename string) string {
	file, err := os.Open(filename)

	if err != nil {
		// error in opening file
		panic("Error: something went wrong with file open (probably you gave wrong filename).")
	}

	scanner := bufio.NewScanner(file)
	genome := ""
	for scanner.Scan() {
		currentLine := scanner.Text() // grabs one line of text and returns a strings
		if len(currentLine) > 0 && currentLine[0] != '>' {
			genome += currentLine
		}
	}
	return genome
}

func CollectReadsFromFASTA(filename string) []string {
	file, err := os.Open(filename)

	if err != nil {
		// error in opening file
		panic("Error: something went wrong with file open (probably you gave wrong filename).")
	}

	scanner := bufio.NewScanner(file) // think of this as a "reader bot"
	reads := make([]string, 0)

	// let's use same trick of using map and not read in duplicate reads.
	readCount := make(map[string]int)
	currentRead := ""
	counter := 0 // for updating user

	// go for as long as the reader bot can still see text
	for scanner.Scan() {
		currentLine := scanner.Text() // grabs one line of text and returns a strings
		if len(currentLine) == 0 {
			continue
		}
		if currentLine[0] != '>' {
			// append the current line to our growing read
			currentRead += currentLine
		} else { // we are at a header
			// the current read is complete! :) append it
			if currentRead != "" && ValidDNAString(currentRead) {
				readCount[currentRead]++
				counter++
				currentRead = ""
				if counter%20000 == 0 {
					fmt.Println("Update: we have processed", counter, "reads.")
				}
			}
			currentRead = ""
		}
	}

	// we have read everything in
	if scanner.Err() != nil {
		panic(scanner.Err())
		// panic("Error: issue in scanning process.")
	}

	file.Close()

	// our reads are living as the keys of the readCount map. Grab them.
	for read := range readCount {
		reads = append(reads, read)
	}

	return reads
}

func ValidDNAString(dna string) bool {
	// check if dna is composed of A's, C's, G's, T's
	if len(dna) < 20 {
		return false
	}
	for _, symbol := range dna {
		if symbol != 'A' && symbol != 'C' && symbol != 'G' && symbol != 'T' {
			return false
		}
	}
	// we made it!
	return true
}

func WriteGenomeToFile(genome string, outFilename string) {
	outFile, err := os.Create(outFilename)
	if err != nil {
		panic("Sorry, couldn't create file!")
	}
	fmt.Fprintln(outFile, genome)
	outFile.Close()
}

func WriteGenomeToFileFASTA(genome string, outFilename string) {
	outFile, err := os.Create(outFilename)
	if err != nil {
		panic("Sorry, couldn't create file!")
	}
	fmt.Fprintln(outFile, ">genome")
	fmt.Fprintln(outFile, genome)
	outFile.Close()
}

// WriteAlignmentToFASTA takes an alignment and a file name and writes the alignment
// to the file as a FASTA. It uses "string_1" and "string_2" as the headers.
func WriteAlignmentToFASTA(a PairwiseAlignment, filename string) {
	file, err := os.Create(filename)
	if err != nil { // panic if anything went wrong
		panic(err)
	}

	writer := bufio.NewWriter(file)
	// first header
	fmt.Fprintln(writer, ">string_1")

	//first string
	fmt.Fprintln(writer, a[0])

	newSymbols := make([]byte, len(a[0]))
	for i := range a[0] {
		if a[0][i] == '-' || a[1][i] == '-' {
			newSymbols[i] = ' '
		} else if a[0][i] == a[1][i] {
			newSymbols[i] = '|'
		} else {
			// mismatch
			newSymbols[i] = '.'
		}
	}

	fmt.Fprintln(writer, string(newSymbols))

	//second string
	fmt.Fprintln(writer, a[1])

	//second header
	fmt.Fprintln(writer, ">string_2")

	writer.Flush()

	file.Close()
}

func WriteContigsToFile(contigs []string, outFilename string) {
	outFile, err := os.Create(outFilename)
	if err != nil {
		panic("Sorry, couldn't create file!")
	}
	for _, str := range contigs {
		fmt.Fprintln(outFile, str)
	}
	outFile.Close()
}

func WriteContigsToFileFASTA(contigs []string, outFilename string) {
	outFile, err := os.Create(outFilename)
	if err != nil {
		panic("Sorry, couldn't create file!")
	}
	for i, str := range contigs {
		fmt.Fprintln(outFile, ">contig"+strconv.Itoa(i))
		fmt.Fprintln(outFile, str)
	}
	outFile.Close()
}
