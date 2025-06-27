func BetaDiversityMatrix(allMaps map[string](map[string]int), distanceMetric string) ([]string, [][]float64) {
	sampleNames := make([]string, 0)
	for names, _ := range allMaps {
		sampleNames = append(sampleNames, names)
	}
	sort.Strings(sampleNames)
	n := len(sampleNames)
	distance := make([][]float64, n)
	for i := range distance {
		distance[i] = make([]float64, n)
	}
	
	if distanceMetric == "Bray-Curtis" {
		for i:= 0; i < n; i++ {
			for j := 0; j < n; j++ {
				distance[i][j] = BrayCurtisDistance(allMaps[sampleNames[i]],allMaps[sampleNames[j]])
			}
		}
	} else if distanceMetric == "Jaccard" {
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				distance[i][j] = JaccardDistance(allMaps[sampleNames[i]], allMaps[sampleNames[j]])
			}
		}
	}

	return sampleNames, distance
}