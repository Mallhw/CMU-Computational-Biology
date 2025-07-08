package main

import (
	"fmt"
)

func main() {
	fmt.Println("Metagenomics!")

	AnalyzeYear("2019")
}

// AnalyzeYear
// Input: a string representing the year or other identifier
// Output: nothing, but we will run a metagenomics analysis and write the results of this analysis to files
func AnalyzeYear(year string) {
	//step 1: let's read in a single file
	filename := fmt.Sprintf("Data/%s_Samples/Fall_Allegheny_1.txt", year)

	freqMap := ReadFrequencyMapFromFile(filename)

	fmt.Println("File read successfully!")

	fmt.Println("File contains:", Richness(freqMap), "total keys.")

	// let's print some statistics
	fmt.Println("Simpson's Index:", SimpsonsIndex(freqMap))

	// let's get another sample
	filename2 := fmt.Sprintf("Data/%s_Samples/Fall_Allegheny_2.txt", year)
	freqMap2 := ReadFrequencyMapFromFile(filename2)

	fmt.Println("Simpsons Index of freq map 2:", SimpsonsIndex((freqMap2)))

	fmt.Println("Bray Curtis distance:", BrayCurtisDistance(freqMap, freqMap2))

	// we don't care about 1 or 2 samples. We want all of them!
	// step 2: read all samples from a directory
	path := fmt.Sprintf("Data/%s_Samples/", year)
	allMaps := ReadSamplesFromDirectory(path)

	fmt.Println("Samples read successfully!")

	fmt.Println("We have", len(allMaps), "total samples.")

	fmt.Println("Now writing richness and evenness maps to file.")

	// step 3: process richness and evenness of all samples
	richness := RichnessMap(allMaps)
	richnessFile := fmt.Sprintf("Matrices/RichnessMatrix_%s.csv", year)
	WriteRichnessMapToFile(richness, richnessFile)

	simpson := SimpsonsMap(allMaps)
	simpsonFile := fmt.Sprintf("Matrices/SimpsonMatrix_%s.csv", year)
	WriteSimpsonsMapToFile(simpson, simpsonFile)

	fmt.Println("Richness and evenness maps written to file successfully.")

	// let's do two beta diversity matrices too
	distanceMetric := "Bray-Curtis"
	sampleNames, mtxBC := BetaDiversityMatrix(allMaps, distanceMetric)
	fmt.Println("Writing BC matrix to file.")
	brayCurtisFile := fmt.Sprintf("Matrices/BrayCurtisBetaDiversityMatrix_%s.csv", year)
	WriteBetaDiversityMatrixToFile(mtxBC, sampleNames, brayCurtisFile)
	fmt.Println("Bray Curtis matrix written to file!")

	fmt.Println("Now writing Jaccard matrix to file.")
	distanceMetric = "Jaccard"
	sampleNames, mtxJac := BetaDiversityMatrix(allMaps, distanceMetric)
	fmt.Println("Writing Jaccard matrix to file.")
	jaccardFile := fmt.Sprintf("Matrices/JaccardBetaDiversityMatrix_%s.csv", year)
	WriteBetaDiversityMatrixToFile(mtxJac, sampleNames, jaccardFile)

	fmt.Println("Success! We are now ready to visualize our data.")

}
