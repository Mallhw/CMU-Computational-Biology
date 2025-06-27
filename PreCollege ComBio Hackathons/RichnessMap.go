func RichnessMap(allMaps map[string](map[string]int)) map[string]int {
	richnessMap := make(map[string]int)
	
	//go through maps in allMaps 
	for i, val := range allMaps {
		richnessMap[i] = Richness(val)
	}

	return richnessMap
}
