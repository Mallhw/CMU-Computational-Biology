func SumEven(k int) int {
    ans := 0
    for x := 0; x <= k; x += 2 {
        ans += x
    }
    return ans
}