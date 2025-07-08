package main

import (
	"fmt"
)

func main() {
	fmt.Println("Sequence alignment!")

	//Debug()

	//ShortSARSDemo()

	//Hemoglobin()

	SARSAlignment()
}

func SARSAlignment() {
	sars1 := ReadFASTAFile("Data/Coronaviruses/SARS-CoV_genome.fasta")

	sars2 := ReadFASTAFile("Data/Coronaviruses/SARS-CoV-2_genome.fasta")

	match := 1.0
	mismatch := 1.0
	gap := 3.0

	sarsAlignment := GlobalAlignment(sars1, sars2, match, mismatch, gap)

	WriteAlignmentToFASTA(sarsAlignment, "Output/SARS_genome_alignment.txt")

	spikeSequence := ReadFASTAFile("Data/Coronaviruses/SARS-CoV_genome_spike_protein.fasta")

	spikeAlignment := GlobalAlignment(spikeSequence, sars2, match, mismatch, gap)

	WriteAlignmentToFASTA(spikeAlignment, "Output/SARS_spike_SARS-2_genome_alignment.txt")
}

func Hemoglobin() {
	zebrafish := ReadFASTAFile("Data/Hemoglobin/Danio_rerio_hemoglobin.fasta")

	human := ReadFASTAFile("Data/Hemoglobin/Homo_sapiens_hemoglobin.fasta")

	cow := ReadFASTAFile("Data/Hemoglobin/Bos_taurus_hemoglobin.fasta")

	gorilla := ReadFASTAFile("Data/Hemoglobin/Gorilla_gorilla_hemoglobin.fasta")

	match := 1.0
	mismatch := 1.0
	gap := 5.0

	// call our algorithm
	alignment1 := GlobalAlignment(zebrafish, human, match, mismatch, gap)

	WriteAlignmentToFASTA(alignment1, "Output/zebrafish_human_hemoglobin.txt")

	alignment2 := GlobalAlignment(cow, human, match, mismatch, gap)

	WriteAlignmentToFASTA(alignment2, "Output/cow_human_hemoglobin.txt")

	alignment3 := GlobalAlignment(gorilla, human, match, mismatch, gap)

	WriteAlignmentToFASTA(alignment3, "Output/gorilla_human_hemoglobin.txt")
}

func Debug() {
	str1 := "GACT"
	str2 := "ATG"

	lcs := LongestCommonSubsequence(str1, str2)
	fmt.Println(lcs)
}

func ShortSARSDemo() {
	sars := ReadFASTAFile("Data/Coronaviruses/SARS-CoV_genome.fasta")
	sars2 := ReadFASTAFile("Data/Coronaviruses/SARS-CoV-2_genome.fasta")

	fmt.Println("LCS length between genomes:", LCSLength(sars, sars2))

	fmt.Println("Edit distance between genomes:", EditDistance(sars, sars2))

	fmt.Println("LCS between genomes:", LongestCommonSubsequence(sars, sars2))
}
