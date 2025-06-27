func MinIntegerArray(list []int) int {
    ans := list[0]
    for x := 0; x < len(list); x++ {
        if(list[x] < ans) {ans = list[x]}
    }
    return ans
}