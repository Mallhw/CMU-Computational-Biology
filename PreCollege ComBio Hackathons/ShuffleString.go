func ShuffleStrings(patterns []string) []string {
	ans := make([]string, 0) 
	perm := rand.Perm(len(patterns))
	for i := 0; i < len(patterns); i++ {
        ans = append(ans, patterns[perm[i]])
	}
	return ans
}