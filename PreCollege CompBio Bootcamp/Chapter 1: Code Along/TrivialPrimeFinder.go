func TrivialPrimeFinder(n int) []bool {
    ans := make([]bool, n+1)
    for x := 1; x < n+1; x++ {
        ans[x] = IsPrime(x)
    }
    return ans
}

//Hint: place your IsPrime() subroutine here
func IsPrime(p int) bool {
    if p == 1 {return false}
    for x := 2; x <= int(math.Sqrt(float64(p))); x++ {
        if(p % x == 0) {return false}
    }
    return true
}