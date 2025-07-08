package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// ReadFASTAFile takes a file name with a single FASTA header and reads out all elements
// as a string other than header elements.
func ReadFASTAFile(filename string) string {
	file, err := os.Open(filename)

	if err != nil {
		// error in opening file
		panic("Error: something went wrong with file open (probably you gave wrong filename).")
	}

	scanner := bufio.NewScanner(file)
	genome := ""
	for scanner.Scan() {
		currentLine := scanner.Text() // grabs one line of text and returns a string
		if len(currentLine) > 0 && currentLine[0] != '>' {
			genome += currentLine
		}
	}
	return genome
}

// WriteAlignmentToFASTA takes an alignment and a file name and writes the alignment
// to the file as a FASTA. It uses "string_1" and "string_2" as the headers.
func WriteAlignmentToFASTA(a Alignment, filename string) {
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

// WriteLocalAlignmentToFasta takes an alignment and a filename.
// It writes the alignment to the file as a FASTA, with indices for local alignment.
func WriteLocalAlignmentToFASTA(a Alignment, filename string, start1, end1, start2, end2 int) {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	// Write headers
	fmt.Fprintln(writer, ">string_1: "+strconv.Itoa(start1)+" to "+strconv.Itoa(end1))
	fmt.Fprintln(writer, a[0]) // Print the complete first string

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

	fmt.Fprintln(writer, a[1]) // Print the complete second string
	fmt.Fprintln(writer, ">string_2: "+strconv.Itoa(start2)+" to "+strconv.Itoa(end2))

}

// PrintAlignment takes an alignment and prints the two corresponding rows of the alignment.
func PrintAlignment(a Alignment) {
	fmt.Println(a[0])
	fmt.Println(a[1])
}
