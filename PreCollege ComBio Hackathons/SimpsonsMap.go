func SimpsonsMap(allMaps map[string](map[string]int)) map[string]float64 {
	simpsonsMap := make(map[string]float64)

	for i, val := range allMaps {
		simpsonsMap[i] = SimpsonsIndex(val)
	}

	return simpsonsMap
}
