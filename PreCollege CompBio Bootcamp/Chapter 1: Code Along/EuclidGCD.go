func EuclidGCD(a, b int) int {
    if b == 0 {
        return a
    }
    return EuclidGCD(b, a%b)
}