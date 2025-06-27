func SieveOfEratosthenes(n int) []bool {
    ans := make([]bool, n+1)
    for x := 2; x < n+1; x++ {
        ans[x] = true
    }
    for i := 2; i < n+1; i++ {
        if ans[i] == true {
            ans = CrossOffMultiples(ans, i)
        }
    }
    return ans
}

// Hint: insert your CrossOffMultiples function here.
func CrossOffMultiples(primeBooleans []bool, p int) []bool {
    for k := 2 * p; k < len(primeBooleans); k += p {
        primeBooleans[k] = false
    }
    return primeBooleans
}