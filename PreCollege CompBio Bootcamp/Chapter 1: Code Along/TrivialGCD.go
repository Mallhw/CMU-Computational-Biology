func TrivialGCD(a, b int) int {
    d := 1
    num := min(a, b)
    for x := 1; x <= num; x++ {
        if(a % x == 0 && b % x == 0) {
            d = x
        }
    }
    return d
}