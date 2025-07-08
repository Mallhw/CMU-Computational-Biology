package main

// RichnessMap takes a map of frequency maps as input.  It returns a map
// whose values are the richness of each sample.
func RichnessMap(allMaps map[string](map[string]int)) map[string]int {
	r := make(map[string]int)

	//range over all of the samples and set their richness values
	for sampleName, currentMap := range allMaps {
		r[sampleName] = Richness(currentMap)
	}

	return r
}
