func CrossOffMultiples(primeBooleans []bool, p int) []bool {
    for k := 2 * p; k < len(primeBooleans); k += p {
        primeBooleans[k] = false
    }
    return primeBooleans
}