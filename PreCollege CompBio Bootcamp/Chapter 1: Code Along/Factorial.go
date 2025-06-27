func Factorial(n int) int {
    x := 1
    ans := 1
    for (x <= n) { 
        ans *= x
        x += 1
    }
    return ans
}