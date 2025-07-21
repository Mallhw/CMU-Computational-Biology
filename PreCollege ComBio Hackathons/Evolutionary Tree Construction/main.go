package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {

	fmt.Println("Happy trees.")

	//HemoglobinUPGMA()

	//SARS2UPGMA()

	Process16SUPGMA(2024)
	Process16SUPGMA(2025)
}

func Process16SUPGMA(year int) {
	fmt.Println("Read in files.")

	// Build input & output paths based on year
	yearDir := fmt.Sprintf("16S_%d", year)
	dataDir := filepath.Join("Data", yearDir, "merged_sequences")
	outDir := filepath.Join("Output", yearDir)

	// Make sure the output dir exists (ignore error if already exists)
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		panic(err)
	}

	// Read all 16S sequence files
	bacteria16SDirectory, err := Read16SFilesFromDirectory(dataDir)
	if err != nil {
		panic(err)
	}

	numberOfSequences := len(bacteria16SDirectory)

	// Collect names & sequences
	names := make([]string, 0, numberOfSequences)
	sequences := make([]string, 0, numberOfSequences)
	for name, sequence := range bacteria16SDirectory {
		names = append(names, name)
		sequences = append(sequences, sequence)
	}

	// Clean up names: drop a trailing "-seqF" if present
	for i, n := range names {
		if strings.HasSuffix(n, "-seqF") {
			names[i] = strings.TrimSuffix(n, "-seqF")
		}
	}

	fmt.Println("Files read. Total number of files:", numberOfSequences)
	fmt.Println("Let's print the length of the strings.")

	// Build pairwise distance matrix
	mtx := CalculateDistanceMatrix(sequences)
	fmt.Println("Distance matrix generated.")

	// Build UPGMA tree
	T := UPGMA(mtx, names)
	fmt.Println("UPGMA complete. Writing tree to file.")

	// Write Newick tree
	WriteNewickToFile(T, outDir, "colonized_bacteria.tre")
	fmt.Println("Tree written to file.")
}

func HemoglobinUPGMA() {
	fmt.Println("Reading in hemoglobin files.")

	stringMap := ReadDNAStringsFromFile("Data/HBA1/hemoglobin_protein.fasta")

	speciesNames, proteinStrings := GetKeyValues(stringMap)

	fmt.Println("Making distance matrix.")

	// make a distance matrix
	mtx := CalculateDistanceMatrix(proteinStrings)

	fmt.Println("Distance matrix made.")

	fmt.Println("Starting UPGMA.")

	t := UPGMA(mtx, speciesNames)

	fmt.Println("Tree built! Writing to file.")

	WriteNewickToFile(t, "Output/HBA1", "hba1.tre")

	fmt.Println("Tree written to file.")
}

func SARS2UPGMA() {
	fmt.Println("Reading in patient virus genomes.")

	directory := "Data/UK-Genomes"

	genomeDatabase := ReadGenomesFromDirectory(directory)
	fmt.Println("Database read! Let's grab some spike proteins.")

	numberOfGenomesPerDate := 3
	numberOfDates := len(genomeDatabase)

	// we will have two arrays of equal length: the spike proteins and the labels of the genomes corresponding to each one
	spikeProteins := make([]string, 0, numberOfGenomesPerDate*numberOfDates)

	labels := make([]string, 0, numberOfGenomesPerDate*numberOfDates)

	// range over the dates and get the spike proteins
	for date := range genomeDatabase {
		count := 0

		for i := 0; i < len(genomeDatabase[date]) && count < numberOfGenomesPerDate; i++ {
			currentGenome := genomeDatabase[date][i]
			// just get the spike protein
			spikeDNA := ExciseSpikeProtein(currentGenome)
			if spikeDNA != "" && ValidDNA(spikeDNA) {
				// grab it
				// translate it to amino acids
				rnaStrand := DNAToRNA(spikeDNA)
				proteinSequence := Translate(rnaStrand, 0)

				// I have the protein, so append it to the collection of proteins
				spikeProteins = append(spikeProteins, proteinSequence)
				currentLabel := date + "_" + strconv.Itoa(i)
				labels = append(labels, currentLabel)

				//increment the count so that I don't sample too many spike proteins
				count++
			}
		}
	}

	fmt.Println("Spike proteins now excised.")

	fmt.Println("Making distance matrix.")

	// make a distance matrix
	mtx := CalculateDistanceMatrix(spikeProteins)

	fmt.Println("Distance matrix made.")

	fmt.Println("Starting UPGMA.")

	t := UPGMA(mtx, labels)

	fmt.Println("Tree built! Writing to file.")

	WriteNewickToFile(t, "Output/UK-Genomes", "sars-cov-2.tre")

	fmt.Println("Tree written to file.")
}
