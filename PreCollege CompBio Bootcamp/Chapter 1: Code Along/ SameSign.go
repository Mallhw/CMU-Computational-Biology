func SameSign(x, y int) bool {
    if x == 0 || y == 0 {
        return true
    }
    return (x > 0) == (y > 0)
}