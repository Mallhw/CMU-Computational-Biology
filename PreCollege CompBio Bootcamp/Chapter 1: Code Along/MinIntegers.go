func MinIntegers(numbers ...int) int {
    if(len(numbers) == 0) {
        return -1
    }
    var ans int
    for i, val := range(numbers) {
        if(i == 0 || val < ans) {
            ans = val
        }
    }
    return ans
}