func FactorialArray(n int) []int {
    a := make([]int, n+1)
    a[0] = 1
    for x := 1; x < n+1; x++ {
        a[x] = x * a[x-1]
    }
    return a
}