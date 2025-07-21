package main

//ReverseComplement takes a DNA string as input and returns its reverse complement.
func ReverseComplement(text string) string {
	return Reverse(Complement(text))
}

//Reverse takes a string and returns the reversed string.
func Reverse(text string) string {
	n := len(text)
	symbols := make([]byte, n)
	for i := range text {
		symbols[i] = text[n-i-1]
	}
	return string(symbols)
}

//Complement takes a DNA string and returns a string in which each base has
//been replaced by its complementary base.
func Complement(text string) string {
	// as with arrays, we can use "range"

	n := len(text)
	symbols := make([]byte, n)

	for i := range text {
		if text[i] == 'A' {
			symbols[i] = 'T'
		} else if text[i] == 'T' {
			symbols[i] = 'A'
		} else if text[i] == 'C' {
			symbols[i] = 'G'
		} else if text[i] == 'G' {
			symbols[i] = 'C'
		}
	}

	return string(symbols)
}

//DNAToRNA takes a DNA string as input and returns the RNA string transcribed from
//this DNA string (i.e., 'T' is replaced with 'U').
func DNAToRNA(genome string) string {
	rna := make([]rune, len(genome))
	for i, symbol := range genome {
		if symbol == 'T' {
			rna[i] = 'U'
		} else {
			rna[i] = symbol
		}
	}
	return string(rna)
}
