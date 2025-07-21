package main

import (
	"fmt"
)

func main() {
	fmt.Println("Genome assembly!")

	SARSOverlapNetworkMinimizersTrimMessyContigs()

	fmt.Println("Program exiting.")
}

func SARSOverlapNetworkMinimizersTrimMessyContigs() {
	fmt.Println("Read in the SARS-CoV-2 genome.")

	genome := ReadGenomeFromFASTA("Data/SARS-CoV-2_genome.fasta")

	fmt.Println("Genome read. Sampling reads.")

	//sample some reads
	readLength := 150
	probability := 0.1
	substitutionErrorRate := 0.01
	insertionErrorRate := 0.01
	deletionErrorRate := 0.01
	reads := SimulateReadsMessy(genome, readLength, probability, substitutionErrorRate, insertionErrorRate, deletionErrorRate)

	fmt.Println("Reads generated! We have", len(reads), "total reads.")

	fmt.Println("Now, make the minimizer map.")

	windowLength := 20
	k := 10

	minimizerDictionary := BuildMinimizerMap(reads, k, windowLength)

	fmt.Println("Minimizer map made. It contains", len(minimizerDictionary), "total keys.")

	fmt.Println("Building overlap network.")

	match := 1.0
	mismatch := 1.0
	gap := 5.0

	threshold := 40.0

	adjList := MakeOverlapNetworkMinimizers(reads, minimizerDictionary, match, mismatch, gap, threshold)

	fmt.Println("Overlap network made!")

	fmt.Println("The network has", len(adjList), "total keys.")

	fmt.Println("The average outdegree is", AverageOutDegree(adjList))

	maxK := 3

	fmt.Println("Now, trimming network by removing transitivity with a maxK value of", maxK)

	trimmedAdjacencyList := TrimNetwork(adjList, maxK)

	fmt.Println("Graph has been trimmed.")

	fmt.Println("Graph now has", len(trimmedAdjacencyList), "total reads.")

	fmt.Println("Average outdegree is", AverageOutDegree(trimmedAdjacencyList))

	fmt.Println("Inferring the maximal non-branching paths.")

	contigs := GenerateContigs(trimmedAdjacencyList)

	fmt.Println("Graph has", len(contigs), "total contigs.")

	fmt.Println("Pipeline finished.")
}

func SARSOverlapNetworkMinimizersTrim() {
	fmt.Println("Read in the SARS-CoV-2 genome.")

	genome := ReadGenomeFromFASTA("Data/SARS-CoV-2_genome.fasta")

	fmt.Println("Genome read. Sampling reads.")

	//sample some reads
	readLength := 150
	probability := 0.1
	reads := SimulateReadsClean(genome, readLength, probability)

	fmt.Println("Reads generated! We have", len(reads), "total reads.")

	fmt.Println("Now, make the minimizer map.")

	windowLength := 20
	k := 10

	minimizerDictionary := BuildMinimizerMap(reads, k, windowLength)

	fmt.Println("Minimizer map made. It contains", len(minimizerDictionary), "total keys.")

	fmt.Println("Building overlap network.")

	match := 1.0
	mismatch := 1.0
	gap := 5.0

	threshold := 40.0

	adjList := MakeOverlapNetworkMinimizers(reads, minimizerDictionary, match, mismatch, gap, threshold)

	fmt.Println("Overlap network made!")

	fmt.Println("The network has", len(adjList), "total keys.")

	fmt.Println("The average outdegree is", AverageOutDegree(adjList))

	maxK := 3

	fmt.Println("Now, trimming network by removing transitivity with a maxK value of", maxK)

	trimmedAdjacencyList := TrimNetwork(adjList, maxK)

	fmt.Println("Graph has been trimmed.")

	fmt.Println("Graph now has", len(trimmedAdjacencyList), "total reads.")

	fmt.Println("Average outdegree is", AverageOutDegree(trimmedAdjacencyList))

}

func SARSOverlapNetwork() {
	fmt.Println("Read in the SARS-CoV-2 genome.")

	genome := ReadGenomeFromFASTA("Data/SARS-CoV-2_genome.fasta")

	fmt.Println("Genome read. Sampling reads.")

	//sample some reads
	readLength := 150
	probability := 0.1
	reads := SimulateReadsClean(genome, readLength, probability)

	fmt.Println("Reads generated! We have", len(reads), "total reads.")

	fmt.Println("Building overlap network.")

	match := 1.0
	mismatch := 1.0
	gap := 5.0

	threshold := 40.0

	adjList := MakeOverlapNetwork(reads, match, mismatch, gap, threshold)

	fmt.Println("Overlap network made!")

	fmt.Println("The network has", len(adjList), "total keys.")

	fmt.Println("The average outdegree is", AverageOutDegree(adjList))
}

func SARSOverlapNetworkMinimizers() {
	fmt.Println("Read in the SARS-CoV-2 genome.")

	genome := ReadGenomeFromFASTA("Data/SARS-CoV-2_genome.fasta")

	fmt.Println("Genome read. Sampling reads.")

	//sample some reads
	readLength := 150
	probability := 0.1
	reads := SimulateReadsClean(genome, readLength, probability)

	fmt.Println("Reads generated! We have", len(reads), "total reads.")

	fmt.Println("Now, make the minimizer map.")

	windowLength := 20
	k := 10

	minimizerDictionary := BuildMinimizerMap(reads, k, windowLength)

	fmt.Println("Minimizer map made. It contains", len(minimizerDictionary), "total keys.")

	fmt.Println("Building overlap network.")

	match := 1.0
	mismatch := 1.0
	gap := 5.0

	threshold := 40.0

	adjList := MakeOverlapNetworkMinimizers(reads, minimizerDictionary, match, mismatch, gap, threshold)

	fmt.Println("Overlap network made!")

	fmt.Println("The network has", len(adjList), "total keys.")

	fmt.Println("The average outdegree is", AverageOutDegree(adjList))
}
