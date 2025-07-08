package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// WriteRichnessMapToFile writes the richness data from a map to a CSV file.
// Input: A map containing sample names as keys and richness values as integers, and a filename string for the output file.
// Output: None. The function writes directly to a file and panics on failure.
func WriteRichnessMapToFile(richness map[string]int, filename string) {
	file, err := os.Create(filename)
	if err != nil { // panic if anything went wrong
		panic(err)
	}

	writer := bufio.NewWriter(file)

	//print headers
	fmt.Fprint(writer, "Sample")
	fmt.Fprint(writer, ",")
	fmt.Fprint(writer, "Richness")
	fmt.Fprintln(writer, "")

	// print sample name and value on each line
	for sampleName, val := range richness {
		fmt.Fprint(writer, sampleName)
		fmt.Fprint(writer, ",")
		fmt.Fprint(writer, val)

		//print new line
		fmt.Fprintln(writer, "")
	}

	writer.Flush()

	file.Close() // the "defer" statement says "do this at the end of the file"
}

// WriteSimpsonsMapToFile writes the Simpson's diversity index data from a map to a CSV file.
// Input: A map containing sample names as keys and Simpson's index values as floats, and a filename string for the output file.
// Output: None. The function writes directly to a file and panics on failure.
func WriteSimpsonsMapToFile(simpson map[string]float64, filename string) {
	file, err := os.Create(filename)
	if err != nil { // panic if anything went wrong
		panic(err)
	}

	writer := bufio.NewWriter(file)

	//print headers
	fmt.Fprint(writer, "Sample")
	fmt.Fprint(writer, ",")
	fmt.Fprint(writer, "SimpsonsIndex")
	fmt.Fprintln(writer, "")

	// print sample name and value on each line
	for sampleName, val := range simpson {
		fmt.Fprint(writer, sampleName)
		fmt.Fprint(writer, ",")
		fmt.Fprint(writer, val)

		//print new line
		fmt.Fprintln(writer, "")
	}

	writer.Flush()

	file.Close() // the "defer" statement says "do this at the end of the file"
}

// WriteBetaDiversityMatrixToFile writes a matrix of beta diversity values and their corresponding sample names to a CSV file.
// Input: A 2D slice of float64 representing the beta diversity matrix, a slice of strings for sample names, and a filename string for the output file.
// Output: None. The function writes directly to a file and panics on failure.
func WriteBetaDiversityMatrixToFile(mtx [][]float64, sampleNames []string, filename string) {
	file, err := os.Create(filename)
	if err != nil { // panic if anything went wrong
		panic(err)
	}

	writer := bufio.NewWriter(file)
	// add gap at start of file
	fmt.Fprint(writer, ",")

	//print all sample names
	for _, name := range sampleNames {
		fmt.Fprint(writer, name)
		fmt.Fprint(writer, ",")
	}
	fmt.Fprintln(writer, "")

	for i := range mtx {
		fmt.Fprint(writer, sampleNames[i])
		fmt.Fprint(writer, ",")
		for j := range mtx[i] {
			fmt.Fprint(writer, mtx[i][j])
			fmt.Fprint(writer, ",")
		}
		fmt.Fprintln(writer, "")
	}

	writer.Flush()

	file.Close() // the "defer" statement says "do this at the end of the file"

}

// ReadSamplesFromDirectory has the following input and output
// Input: a collection of filenames. Each file has one string on each line
// Output: a map whose keys are the sample names and whose values are the frequency map at that site.
func ReadSamplesFromDirectory(directory string) map[string](map[string]int) {

	allMaps := make(map[string](map[string]int))

	dirContents, err := ioutil.ReadDir(directory)
	if err != nil {
		panic("Error reading directory!")
	}

	for _, fileData := range dirContents {
		// what is the file name?
		fileName := fileData.Name()

		// Remove the file extension (".txt") to obtain the name of the sample
		sampleName := strings.Replace(fileName, ".txt", "", 1)

		freqMap := ReadFrequencyMapFromFile(filepath.Join(directory, fileName))

		allMaps[sampleName] = freqMap
	}

	return allMaps
}

// ReadFrequencyMapFromFile
// Input: name of file which contains one string on each line
// Output: the frequency map of strings in the file.
func ReadFrequencyMapFromFile(filename string) map[string]int {

	//first, create our frequency map
	freqMap := make(map[string]int)

	// now, try to open the file.
	file, err := os.Open(filename)

	// was the file open successful?
	if err != nil {
		panic(err)
	}

	defer file.Close() // the "defer" statement says "do this at the end of the file"

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		currentString := scanner.Text() // current line of the text
		freqMap[currentString]++        // increment the current line's # of occurrences
	}

	err2 := scanner.Err()

	if err2 != nil {
		panic(err2)
	}

	return freqMap
}
