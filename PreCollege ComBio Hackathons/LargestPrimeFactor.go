func LargestPrimeFactor(n int) int {
    primes := ListPrimes(n) // primes is a slice of ints with the index of the slice corresponding to int vals
    for i := len(primes) - 1; i >= 0; i-- { // iterating from the back of listprimes
		if n % primes[i] == 0 {
			return primes[i]
		}
	}
	return -1 
}

//ListPrimes takes as input integer n
//returns slice containing all prime numbers up to and including n.
func ListPrimes(n int) []int {
	primeList:= make([]int, 0)

	primeBooleans := SieveOfEratosthenes(n)
	for p := range primeBooleans {
		if primeBooleans[p] {
			primeList = append(primeList, p)
		}
	}

	return primeList
}

//SieveOfEratosthenes takes as input integer n
//returns true if prime, false if otherwise
func SieveOfEratosthenes(n int) []bool {
	primeBooleans := make ([]bool, n + 1)
	//everyone is false by default
	//set 2 onward to true
	for p := 2; p <= n; p++ {
		primeBooleans[p] = true
	}

	for p := 2; float64(p) <= math.Sqrt(float64(n)); p++ {
		//if p is prime, cross off its multiples
		primeBooleans = CrossOffMultiples(primeBooleans, p)
	}

	return primeBooleans
}

//CrossOffMultiples takes as input boolean slice primeBooleans + integer p
//returns updated slice in which all variables in array whose indices are multiples of p have been set to false.
func CrossOffMultiples(primeBooleans []bool, p int) []bool {
	n := len(primeBooleans) - 1

	//consider every multiple of p, starting with 2p, and "cross it off" by setting its corresponding entry of slice to false.
	for k:= 2 * p; k <= n; k += p {
		primeBooleans[k] = false
	}

	return primeBooleans
}