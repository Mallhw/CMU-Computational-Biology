package main

import "sort"

// BetaDiversityMatrix takes a map of frequency maps along with a distance metric
// ("Bray-Curtis" or "Jaccard") as input.
// It returns a slice of strings corresponding to the sorted names of the keys
// in the map, along with a matrix of distances whose (i,j)-th
// element is the distance between the i-th and j-th samples using the input metric.
// Input: a collection of frequency maps samples and a distance metric
// Output: a list of sample names and a distance matrix where D_i,j is the distance between
// sample_i and sample_j according to the given distance metric
func BetaDiversityMatrix(allMaps map[string](map[string]int), distanceMetric string) ([]string, [][]float64) {
	numSamples := len(allMaps)
	sampleNames := make([]string, 0, numSamples)

	//range over allMaps and grab all sample IDs
	for currentSampleName := range allMaps {
		sampleNames = append(sampleNames, currentSampleName)
	}

	// sort the sample IDs!
	sort.Strings(sampleNames)

	mtx := InitializeSquareMatrix(numSamples)

	// all values are zero

	// two optimizations
	// 1. don't compute main diagonal values
	// 2. don't recompute symmetric values

	// range over matrix and set all values
	for r := range mtx {
		for c := r + 1; c < len(mtx[r]); c++ { // start c at first element right of main diagonal
			sampleID1 := sampleNames[r]
			sampleID2 := sampleNames[c]
			sample1 := allMaps[sampleID1]
			sample2 := allMaps[sampleID2]

			dist := 0.0

			if distanceMetric == "Bray-Curtis" {
				dist = BrayCurtisDistance(sample1, sample2)
			} else if distanceMetric == "Jaccard" {
				dist = JaccardDistance(sample1, sample2)
			} else {
				panic("Oh no")
			}

			// let's set mtx[r][c] and mtx[c][r] both equal to our distance
			mtx[r][c] = dist
			mtx[c][r] = dist
		}
	}

	return sampleNames, mtx
}

func InitializeSquareMatrix(n int) [][]float64 {
	mtx := make([][]float64, n)
	for r := range mtx {
		mtx[r] = make([]float64, n)
	}
	return mtx
}
