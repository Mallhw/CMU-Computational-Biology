func AnotherFactorial(n int) int {
    ans := 1
    for x := 1; x <= n; x++ {
        ans *= x
    }
    return ans
}