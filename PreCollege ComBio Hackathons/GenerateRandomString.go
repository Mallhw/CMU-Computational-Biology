func GenerateRandomGenome(length int) string {
	ans := ""
	for i := 0; i < length; i++ {
		rand := rand.Intn(4) + 1
		if rand == 1 {
			ans += "A"
		} else if rand == 2 {
			ans += "T"
		} else if rand == 3 {
			ans += "G"
		} else if rand == 4 {
			ans += "C"
		}
	}
	return ans
}