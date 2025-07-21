package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
)

type IO interface {

	/* ReadMatrixFromFile
	   REQUIRES : Valid file path. File provided is a tab-separated matrix.
	              First line is number of samples. Left column are row names.
	              All numbers in matrix can be parsed as floats.
	   ENSURES  : Returns file matrix into float matrix and list of row names.
	*/
	ReadMatrixFromFile(string) (DistanceMatrix, []string)

	/* ReadStringsFromFile
	   REQUIRES : Valid file path. File of strings, one string per line.
		 ENSURES  : String array corresponding to strings in file.
	*/
	ReadStringsFromFile(string) []string

	/* ReadDNAStringsFromFile
	   REQUIRES : Valid file path. File provided annotated as outlined in
		 						README. (each label MUST BE UNIQUE...see ensures below)
		 ENSURES  : String dictionary, where keys are annotation labels and values
		            are right column.
	*/
	ReadDNAStringsFromFile(string) map[string]string

	/* PrintGraphViz
		   REQUIRES : Tree is completed. (a tree is considered "completed" if
		                                      all fields are populated)
	     ENSURES  : Prints a visualization of the given tree.
	*/
	PrintGraphViz(Tree) // void

	/* ToNewick
	   REQUIRES : Tree is completed. (a tree is considered "completed" if
	                                  all fields are populated)
	   ENSURES  : Returns formatted string corresponding to tree.
	              (more specifically Newick format, which is a popular medium
	               used for data visualization software)
	*/
	ToNewick(Tree) string

	/* CreateCSV
	     REQUIRES : Tree is completed. (a tree is considered "completed" if
	                                    all fields are populated)
									Labels MUST be FASTA annotated.
									List of categories for annotation table. For SARS-Cov-2,
									use ["Wuhan","Italy","USA"].
	     ENSURES  : \result is string annotation table.
			 DESCRIP  : Creates CSV annotation table for Newick tree. Use ONLY for
			            data visualization.
	*/
	CreateCSV(Tree, []string) string

	/* CreateDistanceMatrix
	     REQUIRES : File name is valid and setting in {W, F}
		 	              W: File is FASTA format or annotated strings.
		 							  F: File is strings, one string per line.
	     ENSURES  : \result is a valid distance matrix, slice of labels.
			 DESCRIP  : Given raw DNA strings, annotated or unannotated, produces a
			            symmetric pairwise distance matrix. Returns a slice of annotations
									where labels[i] is the label for matrix[i]. If unannotated, produces
									dummy labels for each string.
	*/
	CreateDistanceMatrix(string, string) (DistanceMatrix, []string)
}

func SequenceOrder(T Tree) []string {
	return ReturnSequenceOrder(T[len(T)-1])
}

func ReturnSequenceOrder(node *Node) []string {
	if node.Child1 == nil {
		return []string{node.Label}
	} else {
		return append(ReturnSequenceOrder(node.Child1), ReturnSequenceOrder(node.Child2)...)
	}
}

func CreateCSV(tree Tree, categories []string) string {
	var freqDict = make(map[string]int, 0)
	i, c := 1, 0
	var count = &c
	for _, item := range categories {
		freqDict[item] = i
		i += 1
	}
	return "##,continent,color\n" + subtreeCSV(tree[len(tree)-1], freqDict, count)
}

func subtreeCSV(node *Node, freqDict map[string]int, count *int) string {
	if node.Child1 == nil {
		*count = *count + 1
		var category = getCatFASTA(node.Label)
		return strconv.Itoa(*count) + "," + category + "," + strconv.Itoa(freqDict[category]) + "\n"
	} else {
		return subtreeCSV(node.Child1, freqDict, count) + subtreeCSV(node.Child2, freqDict, count)
	}
}

func getCatFASTA(annotation string) string {
	var bars = strings.Split(annotation, "|")
	return bars[0]
}

func ToNewick(tree Tree) string {
	return "(" + subtreeNewickAges(tree[len(tree)-1]) + ");"
}

func ToNewickL(tree Tree) string {
	return "(" + subtreeNewickL(tree[len(tree)-1]) + ");"
}

func ToNewickAges(tree Tree) string {
	return "(" + subtreeNewickAges(tree[len(tree)-1]) + ");"
}

func subtreeNewick(node *Node) string {
	if node.Child1 == nil {
		return node.Label
	} else {
		return "(" + subtreeNewick(node.Child1) + "," + subtreeNewick(node.Child2) + ")"
	}
}

func subtreeNewickL(node *Node) string {
	if node.Child1 == nil {
		return node.Label
	} else {
		return "(" + subtreeNewickL(node.Child2) + "," + subtreeNewickL(node.Child1) + ")"
	}
}

func subtreeNewickAges(node *Node) string {
	if node.Child1 == nil {
		return node.Label + ":" + fmt.Sprintf("%.2f", node.Age)
	} else {
		return "(" + subtreeNewickAges(node.Child1) + "," + subtreeNewickAges(node.Child2) + "):" + fmt.Sprintf("%.2f", node.Age)
	}
}

func ReadDNAStringsFromFile(fileName string) map[string]string {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error: couldn't open the file")
		os.Exit(1)
	}
	var lines []string = make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if scanner.Err() != nil {
		fmt.Println("Sorry: there was some kind of error during the file reading")
		os.Exit(1)
	}
	file.Close()

	var dnaMap = make(map[string]string, 0)
	var curLabel = ""

	for idx, str := range lines {
		if idx%2 == 0 {
			curLabel = str
		} else {
			dnaMap[curLabel] = str
		}
	}
	return dnaMap
}

func ReadStringsFromFile(fileName string) []string {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error: couldn't open the file")
		os.Exit(1)
	}
	var lines []string = make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if scanner.Err() != nil {
		fmt.Println("Sorry: there was some kind of error during the file reading")
		os.Exit(1)
	}
	file.Close()

	var dnaMap = make([]string, 0)
	for _, str := range lines {
		dnaMap = append(dnaMap, str)
	}
	return dnaMap
}

// Read16SFilesFromDirectory reads all FASTA files in a directory and maps each sequence header to its sequence.
func Read16SFilesFromDirectory(directory string) (map[string]string, error) {
	sequenceMap := make(map[string]string)

	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".fasta") {
			f, err := os.Open(directory + "/" + file.Name())
			if err != nil {
				return nil, err
			}
			defer f.Close()

			scanner := bufio.NewScanner(f)
			var currentHeader string
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, ">") {
					currentHeader = line[1:strings.Index(line, " ")]
				} else {
					sequenceMap[currentHeader] += line
				}
			}
			if err := scanner.Err(); err != nil {
				return nil, err
			}
		}
	}
	return sequenceMap, nil
}

// ReadMatrixFromFile takes a file name and reads the information in this file to produce
// a distance matrix and a slice of strings holding the species names.  The first line of the
// file should contain the number of species.  Each other line contains a species name
// and its distance to each other species.
func ReadMatrixFromFile(fileName string) (DistanceMatrix, []string) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error: couldn't open the file")
		os.Exit(1)
	}
	var lines []string = make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if scanner.Err() != nil {
		fmt.Println("Sorry: there was some kind of error during the file reading")
		os.Exit(1)
	}
	file.Close()

	mtx := make(DistanceMatrix, 0)
	speciesNames := make([]string, 0)

	for idx, _ := range lines {
		if idx >= 1 {
			row := make([]float64, 0)
			nums := strings.Split(lines[idx], "\t")
			for i, num := range nums {
				if i == 0 {
					speciesNames = append(speciesNames, num)
				} else {
					n, err := strconv.ParseFloat(num, 64)
					if err != nil {
						fmt.Println("Error: Wrong format of matrix!")
						os.Exit(1)
					}
					row = append(row, n)
				}
			}
			mtx = append(mtx, row)
		}
	}
	return mtx, speciesNames
}

// PrintGraphViz prints the tree in GraphViz format, where directed = true
// if we desire to print a directed graph and directed = false for an
// undirected graph.
func PrintGraphViz(t Tree) {
	fmt.Println("strict digraph {")
	for i := range t {
		if t[i].Child1 != nil && t[i].Child2 != nil {
			//print first edge
			fmt.Print("\"", t[i].Label, "\"")
			fmt.Print("->")
			fmt.Print("\"", t[i].Child1.Label, "\"")
			fmt.Print("[label = \"")
			fmt.Printf("%.2f", t[i].Age-t[i].Child1.Age)
			fmt.Print("\"]")
			fmt.Println()

			//print second edge
			fmt.Print("\"", t[i].Label, "\"")
			fmt.Print("->")
			fmt.Print("\"", t[i].Child2.Label, "\"")
			fmt.Print("[label = \"")
			fmt.Printf("%.2f", t[i].Age-t[i].Child2.Age)
			fmt.Print("\"]")
			fmt.Println()
		}
	}
	fmt.Println("}")
}

// /Users/kunaljoshi/go/src/EvolutionaryTrees/Datasets/Input/SARS-Cov-2/distance.mtx

func WriteToDistanceMatrix(filePath string, fileDest string, setting string) {

	var mtx DistanceMatrix
	var labels []string

	switch setting {
	case "W":
		var dnaMap = ReadDNAStringsFromFile(filePath)
		var dnaLabels, dnaStrings = GetKeyValues(dnaMap)
		labels = dnaLabels
		mtx = CalculateDistanceMatrix(dnaStrings)
	case "F":
		var dnaMap = CreateFrequencyDNAMap(ReadStringsFromFile(filePath))
		var dnaLabels, dnaStrings = GetKeyValues(dnaMap)
		labels = dnaLabels
		mtx = CalculateDistanceMatrix(dnaStrings)
	default:
		panic("Not valid setting.")
	}

	F, err := os.Create(fileDest + "/distance.mtx")
	if err != nil {
		fmt.Println(err)
		return
	}

	accum := ""
	for i := range mtx {
		accum += labels[i] + "\t"
		for j := range mtx {
			accum += fmt.Sprintf("%.2f", mtx[i][j]) + "\t"
		}
		accum = accum[0 : len(accum)-1]
		accum += "\n"
	}

	_, err = F.WriteString(strconv.Itoa(len(mtx)) + "\n" + accum)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = F.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

}

func WriteAlignmentToFile(algn Alignment, labels []string, fileDest string, fileName string) {

	accum := ""
	for i, string := range algn {
		accum += ">" + labels[i] + "\n"
		accum += string + "\n"
	}

	F, err := os.Create(fileDest + "/" + fileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = F.WriteString(accum)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = F.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

}

func WriteNewickToFile(t Tree, fileDest string, fileName string) {

	newickString := ToNewick(t)

	F, err := os.Create(fileDest + "/" + fileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = F.WriteString(newickString)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = F.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

}

func WriteCSVToFile(t Tree, cats []string, fileDest string, fileName string) {

	csvString := CreateCSV(t, cats)

	F, err := os.Create(fileDest + "/" + fileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = F.WriteString(csvString)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = F.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

}

func CreateDistanceMatrix(filePath string, setting string) (DistanceMatrix, []string) {

	var mtx DistanceMatrix
	var labels []string

	switch setting {
	case "W":
		var dnaMap = ReadDNAStringsFromFile(filePath)
		var dnaLabels, dnaStrings = GetKeyValues(dnaMap)
		labels = dnaLabels
		mtx = CalculateDistanceMatrix(dnaStrings)
	case "F":
		var dnaMap = CreateFrequencyDNAMap(ReadStringsFromFile(filePath))
		var dnaLabels, dnaStrings = GetKeyValues(dnaMap)
		labels = dnaLabels
		mtx = CalculateDistanceMatrix(dnaStrings)
	default:
		panic("Not valid setting.")
	}

	return mtx, labels

}

func ProcessBacteria(filePath string) {

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error: couldn't open the file")
		os.Exit(1)
	}
	var lines []string = make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if scanner.Err() != nil {
		fmt.Println("Sorry: there was some kind of error during the file reading")
		os.Exit(1)
	}
	file.Close()

	accum := ""
	freqDict := make(map[string]int, 0)
	for _, line := range lines {

		r := csv.NewReader(strings.NewReader(line))
		r.Comma = ' ' // space
		fields, err := r.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		label := fields[2] + "|" + fields[5]
		if label == "NA|NA" || strings.Contains(label, "|NA") {
			continue
		}

		num := 1
		_, exists := freqDict[label]
		if exists {
			num = freqDict[label]
			freqDict[label]++
		} else {
			freqDict[label] = 2
		}

		accum = accum + fields[2] + "|" + fields[5] + "|" + strconv.Itoa(num) + "\n"
		accum = accum + fields[0] + "\n"
	}
	fmt.Println(accum)

}

/*
Everything below here is useful for our variant finding code.
*/

func WriteVariantPercentagesToCSV(frequenciesByDate map[string](map[string]float64), filename string) {

	file, err := os.Create(filename)
	if err != nil { // panic if anything went wrong
		panic(err)
	}

	writer := bufio.NewWriter(file)

	//print headers
	fmt.Fprint(writer, "Date")
	fmt.Fprint(writer, ",")
	fmt.Fprint(writer, "Percentage Alpha")
	fmt.Fprint(writer, ",")
	fmt.Fprint(writer, "Percentage Delta")
	fmt.Fprint(writer, ",")
	fmt.Fprint(writer, "Percentage Omicron BA.1")
	fmt.Fprint(writer, ",")
	fmt.Fprint(writer, "Percentage Other")
	fmt.Fprint(writer, ",")
	fmt.Fprint(writer, "Percentage Unidentified")
	fmt.Fprintln(writer, "")

	// grab keys in map so we can sort them
	dates := make([]string, 0, len(frequenciesByDate))
	for date := range frequenciesByDate {
		dates = append(dates, date)
	}

	sort.Strings(dates)

	//write each line of the file

	for _, date := range dates {
		fmt.Fprint(writer, date)
		fmt.Fprint(writer, ",")
		fmt.Fprint(writer, frequenciesByDate[date]["alpha"])
		fmt.Fprint(writer, ",")
		fmt.Fprint(writer, frequenciesByDate[date]["delta"])
		fmt.Fprint(writer, ",")
		fmt.Fprint(writer, frequenciesByDate[date]["omicron BA.1"])
		fmt.Fprint(writer, ",")
		fmt.Fprint(writer, frequenciesByDate[date]["other"])
		fmt.Fprint(writer, ",")
		fmt.Fprint(writer, frequenciesByDate[date]["unidentified"])

		//print new line
		fmt.Fprintln(writer, "")
	}

	//close out file

	writer.Flush()

	file.Close()
}

// ReadGenomesFromDirectory has the following input and output
// Input: a directory name containing genomes within a collection of dates.
// Output: a map of dates to the genomes contained at these dates.
func ReadGenomesFromDirectory(directory string) map[string]([]string) {

	genomeDatabase := make(map[string]([]string))

	dirContents, err := ioutil.ReadDir(directory)
	if err != nil {
		panic("Error reading directory!")
	}

	for _, folderName := range dirContents {
		//folder name will give us the name of the folder, which is the date
		if folderName.IsDir() {

			date := folderName.Name()

			// now we need to read out the appropriate file in the directory

			fileName := date + ".fasta"

			genomeDatabase[date] = ReadStringsFromFASTA(directory + "/" + date + "/" + fileName)
		}
	}

	return genomeDatabase
}

func ReadStringsFromFASTA(filename string) []string {
	file, err := os.Open(filename)

	if err != nil {
		// error in opening file
		panic("Error: something went wrong with file open (probably you gave wrong filename).")
	}

	scanner := bufio.NewScanner(file) // think of this as a "reader bot"
	reads := make([]string, 0)

	currentRead := ""
	counter := 0 // for updating user

	// go for as long as the reader bot can still see text
	for scanner.Scan() {
		currentLine := scanner.Text() // grabs one line of text and returns a strings
		if currentLine[0] != '>' {
			// append the current line to our growing read
			currentRead += currentLine
		} else { // we are at a header
			// the current read is complete! :) append it
			if currentRead != "" {
				reads = append(reads, currentRead)
				counter++
				currentRead = ""
				if counter%20000 == 0 {
					fmt.Println("Update: we have processed", counter, "reads.")
				}
			}
		}
	}

	if currentRead != "" {
		reads = append(reads, currentRead)
		counter++
	}

	// we have read everything in
	if scanner.Err() != nil {
		panic("Error: issue in scanning process.")
	}

	file.Close()

	return reads
}

// Writes a beta diversity matrix to a file
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
