func IsPrime(p int) bool {
    if p == 1 {return false}
    for x := 2; x <= int(math.Sqrt(float64(p))); x++ {
        if(p % x == 0) {return false}
    }
    return true
}