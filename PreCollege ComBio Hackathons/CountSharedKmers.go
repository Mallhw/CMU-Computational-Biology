func CountSharedKmers(str1, str2 string, k int) int {
	count := 0
	map1 := make(map[string]int)
	map2 := make(map[string]int)
	for i := 0; i < len(str1)-k+1; i++ {
		map1[str1[i:i+k]]++
	}
	for i := 0; i < len(str2)-k+1; i++ {
		map2[str2[i:i+k]]++
	}
	for key, value := range map1 {
		count += min(value, map2[key])
	}
	return count
}