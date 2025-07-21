package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// GetKeyValues takes a map of strings to strings as input.
// It returns two slices corresponding to the keys and values of the
// map, respectively.
func GetKeyValues(dnaMap map[string]string) ([]string, []string) {
	var keys = make([]string, 0, len(dnaMap))
	var values = make([]string, 0, len(dnaMap))
	for k, v := range dnaMap {
		keys = append(keys, k)
		values = append(values, v)
	}
	return keys, values
}

// GetKeyValuesI takes a map of strings to ints as input.
// It returns two slices corresponding to the keys and values of the
// map, respectively.
func GetKeyValuesI(dnaMap map[string]int) ([]string, []int) {
	var keys = make([]string, 0, len(dnaMap))
	var values = make([]int, 0, len(dnaMap))
	for k, v := range dnaMap {
		keys = append(keys, k)
		values = append(values, v)
	}
	return keys, values
}

/* CalcDistancematrix : []string -> DistanceMatrix
   REQUIRES           : true
	 ENSURES            : len(\result) == len(\result[0]) ; \result[i][j] == \result[j][i]
												len(\result) == len(patterns)
	 DESCRIP            : Given a string slice, creates a pairwise distance matrix
	                      between all strings. Uses alignment distance heuristic, better
												known as Levenshtein.
*/

// CalculateDistanceMatrix takes a collection of strings as input.
// It returns a distance matrix in which element (i,j)
// is the edit (Levenshtein) distance between strings i and j.
func CalculateDistanceMatrix(patterns []string) DistanceMatrix {
	numSamples := len(patterns)
	mtx := InitializeMatrix(numSamples, numSamples)

	// Progress Indicator
	n := ((numSamples * numSamples) + numSamples) / 2
	counter := 0

	for i := 0; i < numSamples; i++ {
		for j := i; j < numSamples; j++ {
			d := EditDistance(patterns[i], patterns[j])
			mtx[i][j] = d
			mtx[j][i] = d

			if counter%100 == 0 {
				fmt.Println(fmt.Sprintf("%.2f", float64(counter)/float64(n)*100.0) + "% finished with matrix generation.")

			}
			counter++
		}
	}
	return mtx
}

// EditDistance takes two strings and returns the edit distance between them.
func EditDistance(seq1 string, seq2 string) float64 {
	var seq1Size = len(seq1) + 1
	var seq2Size = len(seq2) + 1
	var matrix = InitializeMatrix(seq1Size, seq2Size)

	for i := 0; i < seq1Size; i++ {
		matrix[i][0] = float64(i)
	}
	for j := 0; j < seq2Size; j++ {
		matrix[0][j] = float64(j)
	}

	for i := 1; i < seq1Size; i++ {
		for j := 1; j < seq2Size; j++ {
			if seq1[i-1] == seq2[j-1] {
				matrix[i][j] =
					MinFloat(
						matrix[i-1][j]+5,
						matrix[i-1][j-1],
						matrix[i][j-1]+5)
			} else {
				matrix[i][j] =
					MinFloat(
						matrix[i-1][j]+5,
						matrix[i-1][j-1]+1,
						matrix[i][j-1]+5)
			}
		}
	}
	return matrix[seq1Size-1][seq2Size-1]

}

// InitializeMatrix takes integers m and n and
// returns a DistanceMatrix of dimensions m x n with default values.
func InitializeMatrix(m int, n int) DistanceMatrix {
	mtx := make([][]float64, m)

	for i := range mtx {
		mtx[i] = make([]float64, n)
	}
	return mtx
}

// CreateFrequencyDNAMap takes a slice of strings. It produces a dictionary where dict[i]
// corresponds to slice[i] in the slice. Essentially provides dummy
// labels for unannotated species.
func CreateFrequencyDNAMap(patterns []string) map[string]string {
	var freqMap = CreateFrequencyMap(patterns)
	var keys, _ = GetKeyValuesI(freqMap)

	var dnaMap = make(map[string]string)
	for i := 0; i < len(keys); i++ {
		dnaMap[strconv.Itoa(i)] = keys[i]
	}
	return dnaMap
}

// CreateFrequencyMap takes a collection of strings and returns
// the frequency table of these strings, mapping a string
// to its number of occurrences.
func CreateFrequencyMap(patterns []string) map[string]int {
	freqMap := make(map[string]int)
	for _, val := range patterns {
		freqMap[val]++
	}
	return freqMap
}

/************************************************
  MISCELLANEOUS
************************************************/

// MinFloat returns the minimum of an arbitrary collection of floats.
func MinFloat(nums ...float64) float64 {
	m := 0.0
	// nums gets converted to an array
	for i, val := range nums {
		if val < m || i == 0 {
			// update m
			m = val
		}
	}
	return m
}

// MaxFloat returns the maximum of an arbitrary collection of floats.
func MaxFloat(nums ...float64) float64 {
	m := 0.0
	// nums gets converted to an array
	for i, val := range nums {
		if val > m || i == 0 {
			// update m
			m = val
		}
	}
	return m
}

/************************************************
 VARIANT FINDING
************************************************/
//Takes in a database mapping strings (dates) to lists of genomes
//Returns a map of string keys identifying each sample to kmer frequency maps (map[string]int)
//Keys have the format date_index(01-100)_variant
func KmerMapsFromGenomeDatabase(database map[string]([]string), k int) map[string](map[string]int) {
	allMaps := make(map[string](map[string]int))

	for date, genomes := range database {
		fmt.Println(date, "with", len(genomes), "genomes")
		for i, str := range genomes {
			if i > 10 {
				allMaps[date+"_"+strconv.Itoa(i)] = MakeKmerTable(str, k)
			} else {
				//Add a 0 to all keys in the first 9 genomes so they appear in the table in the order they are in the database
				allMaps[date+"_0"+strconv.Itoa(i)] = MakeKmerTable(str, k)
			}
		}
	}

	return allMaps
}

// Takes in a database mapping strings (dates) to lists of genomes
// Returns a map of string keys identifying each sample to kmer frequency maps (map[string]int)
// Keys have the format date_index(01-100)_variant (labelled meaning labelled using variant classifier)
func KmerMapsFromGenomeDatabaseLabelled(database map[string]([]string), k int) map[string](map[string]int) {
	allMaps := make(map[string](map[string]int))

	for date, genomes := range database {
		fmt.Println(date, "with", len(genomes), "genomes")
		for i, str := range genomes {
			variant := ClassifyVariant(str)
			if i > 10 {
				allMaps[date+"_"+strconv.Itoa(i)+"_"+variant] = MakeKmerTable(str, k)
			} else {
				//Add a 0 to all keys in the first 9 genomes so they appear in the table in the order they are in the database
				allMaps[date+"_0"+strconv.Itoa(i)+"_"+variant] = MakeKmerTable(str, k)
			}
		}
	}

	return allMaps
}

// Takes a string and makes a map with keys being the kmers of the string,
// and values being the frequency of that kmer in the string.
// Only includes kmers containing only A's C's T's and G's
func MakeKmerTable(str string, k int) map[string](int) {
	kmers := KmerComposition(str, k)
	kmerMap := make(map[string]int)
	for _, kmer := range kmers {
		if ValidDNA(kmer) {
			if _, in := kmerMap[kmer]; !in {
				kmerMap[kmer] = 0
			}
			kmerMap[kmer] += 1
		}
	}
	return kmerMap
}

// ValidDNA Checks if a string is made up of only the characters 'A', 'C', 'T', and 'G'
func ValidDNA(str string) bool {
	for _, ch := range str {
		if !(ch == 'A' || ch == 'C' || ch == 'T' || ch == 'G') {
			return false
		}
	}
	return true
}

// KmerComposition returns the k-mer composition (all k-mer substrings) of a given genome.
func KmerComposition(genome string, k int) []string {
	n := len(genome)
	kmers := make([]string, n-k+1)
	// range through and grab all substrings
	for i := 0; i < n-k+1; i++ {
		kmers[i] = genome[i : i+k]
	}

	return kmers
}

// JaccardDistance takes two frequency maps and returns the Jaccard
// distance between them.
func JaccardDistance(map1 map[string]int, map2 map[string]int) float64 {
	c := SumOfMinima(map1, map2)
	t := SumOfMaxima(map1, map2)
	j := 1 - (float64(c) / float64(t))
	return j
}

// SumOfMaxima takes two frequency maps as input.
// It identifies the keys that are shared by the two frequency maps
// and returns the sum of the corresponding values. (a.k.a. "union")
// SumOfMaxima takes two frequency maps as input.
// It identifies the keys that are shared by the two frequency maps
// and returns the sum of the corresponding values. (a.k.a. "union")
func SumOfMaxima(map1 map[string]int, map2 map[string]int) int {
	c := 0

	for key := range map2 {
		_, exists := map1[key]
		if exists {
			c += Max2(map1[key], map2[key])
		} else {
			c += map2[key]
		}
	}
	for key := range map1 {
		_, exists := map2[key]
		if !exists {
			c += map1[key]
		}
	}

	// panic if c is equal to zero since we don't want 0/0
	if c == 0 {
		panic("Error: no species common to two maps given to SumOfMaxima")
	}

	return c
}

// Max2 takes two integers and returns their maximum.
func Max2(n1, n2 int) int {
	if n1 < n2 {
		return n2
	}
	return n1
}

// BetaDiversityMatrix takes a map of frequency maps along with a distance metric
// ("Bray-Curtis" or "Jaccard") as input.
// It returns a slice of strings corresponding to the sorted names of the keys
// in the map, along with a matrix of distances whose (i,j)-th
// element is the distance between the i-th and j-th samples using the input metric.
// Input: a collection of frequency maps samples and a distance metric
// Output: a list of sample names and a distance matrix where D_i,j is the distance between
// sample_i and sample_j according to the given distance metric
func BetaDiversityMatrix(allMaps map[string](map[string]int), distMetric string) ([]string, [][]float64) {

	// first, grab all sample names and sort them
	numSamples := len(allMaps)
	sampleNames := make([]string, 0)

	for name := range allMaps {
		//set current sample name
		sampleNames = append(sampleNames, name)
	}
	//but our slice isn't sorted, so let's do that

	sort.Strings(sampleNames)

	// now form the distance matrix

	mtx := InitializeSquareMatrix(numSamples)

	//range over all pairs of our samples and fill in all the matrix values
	counter := 0
	for i := 0; i < numSamples; i++ {
		for j := i; j < numSamples; j++ {
			var dist float64
			if counter%1000 == 0 {
				fmt.Println(counter, "Jaccard distances taken for diversity matrix")
			}
			counter++
			map1 := allMaps[sampleNames[i]]
			map2 := allMaps[sampleNames[j]]

			//which distance do you want?
			if distMetric == "Bray-Curtis" {
				dist = BrayCurtisDistance(map1, map2)
			} else if distMetric == "Jaccard" {
				dist = JaccardDistance(map1, map2)
			} else {
				panic("Error: invalid distance function name given to BetaDiversityMatrix.")
			}

			//all that remains is to set the appropriate matrix values ...
			mtx[i][j] = dist
			mtx[j][i] = dist
		}
	}

	return sampleNames, mtx
}

// InitializeSquareMatrix takes an integer n and returns an n x n square matrix
// of 0 default floats.
func InitializeSquareMatrix(n int) [][]float64 {
	mtx := make([][]float64, n)
	//now we need to make each row

	for i := range mtx {
		mtx[i] = make([]float64, n)
	}

	return mtx
}

// BrayCurtisDistance takes two frequency maps and returns the Bray-Curtis
// distance between them.
func BrayCurtisDistance(map1 map[string]int, map2 map[string]int) float64 {
	c := SumOfMinima(map1, map2)
	s1 := SampleTotal(map1)
	s2 := SampleTotal(map2)

	//don't allow the situation in which we have zero richness.
	if s1 == 0 || s2 == 0 {
		panic("Error: sample given to BrayCurtisDistance() has no positive values.")
	}

	av := Average(float64(s1), float64(s2))
	return 1 - (float64(c) / av)
}

// Average takes two floats and returns their average.
func Average(x, y float64) float64 {
	return (x + y) / 2.0
}

// SumOfMinima takes two frequency maps as input.
// It identifies the keys that are shared by the two frequency maps
// and returns the sum of the corresponding values.
func SumOfMinima(map1 map[string]int, map2 map[string]int) int {
	c := 0

	for key := range map1 {
		_, exists := map2[key]
		if exists {
			c += Min2(map1[key], map2[key])
		}
	}

	return c
}

// Min2 takes two integers and returns their minimum.
func Min2(n1, n2 int) int {
	if n1 < n2 {
		return n1
	}
	return n2
}

func SampleTotal(freqMap map[string]int) int {
	total := 0
	for _, val := range freqMap {
		total += val
	}
	return total
}

func getCats(labels []string) []string {
	cats := make([]string, 0)
	for _, label := range labels {
		cats = append(cats, strings.Split(label, "|")[0])
	}
	return unique(cats)
}

func unique(strSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range strSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func rearrangeStrings(newLabels []string, oldLabels []string, patternsOld []string) []string {
	patternsNew := make([]string, 0)
	for _, newLabel := range newLabels {
		j := getIndex(oldLabels, newLabel)
		patternsNew = append(patternsNew, patternsOld[j])
	}
	return patternsNew
}

func getIndex(arr []string, target string) int {
	for i, str := range arr {
		if str == target {
			return i
		}
	}
	return 0
}
