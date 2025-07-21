package main

//If it is a start codon, then we transcribe until one of two things happens.
//We either hit the end of the string, in which case we return "",
//or we hit a stop codon, in which we case we return the translated amino acid string.
func Translate(rnaStrand string, startIndex int) string {
	//codons_table is a map of strings to strings that maps
	//RNA codons to single amino acids according to the genetic code.
	var codons_table = map[string]string{
		"UUU": "F", "UUC": "F", "UUA": "L", "UUG": "L",
		"CUU": "L", "CUC": "L", "CUA": "L", "CUG": "L",
		"AUU": "I", "AUC": "I", "AUA": "I", "AUG": "M",
		"GUU": "V", "GUC": "V", "GUA": "V", "GUG": "V",
		"UCU": "S", "UCC": "S", "UCA": "S", "UCG": "S",
		"CCU": "P", "CCC": "P", "CCA": "P", "CCG": "P",
		"ACU": "T", "ACC": "T", "ACA": "T", "ACG": "T",
		"GCU": "A", "GCC": "A", "GCA": "A", "GCG": "A",
		"UAU": "Y", "UAC": "Y", "UAA": "x", "UAG": "x",
		"CAU": "H", "CAC": "H", "CAA": "Q", "CAG": "Q",
		"AAU": "N", "AAC": "N", "AAA": "K", "AAG": "K",
		"GAU": "D", "GAC": "D", "GAA": "E", "GAG": "E",
		"UGU": "C", "UGC": "C", "UGA": "x", "UGG": "W",
		"CGU": "R", "CGC": "R", "CGA": "R", "CGG": "R",
		"AGU": "S", "AGC": "S", "AGA": "R", "AGG": "R",
		"GGU": "G", "GGC": "G", "GGA": "G", "GGG": "G",
	}

	if len(rnaStrand) < 3 {
		panic("Error: length of genome given to Translate is shorter than 3.")
	}
	if startIndex > len(rnaStrand)-3 {
		panic("Error: startIndex given to Translate is too close to the end of the genome.")
	}

	//declaring our translated string
	var translatedPeptide string

	// range through the strand and tack on amino acids until we
	// hit the end of the strand or a stop codon.
	for i := startIndex; i <= len(rnaStrand)-3; i += 3 {
		_, ok := codons_table[rnaStrand[i:i+3]]
		if !ok {
			translatedPeptide += "-"
		} else {
			if codons_table[rnaStrand[i:i+3]] == "x" { // stop codon
				return translatedPeptide
			} else {
				translatedPeptide += codons_table[rnaStrand[i:i+3]]
			}
		}

	}
	// if we reach here, then we hit the end of the string without encountering a stop codon.
	// Return what we currently have.
	return translatedPeptide
}
