package main

func IdentifyVariantPercentages(genomeDatabase map[string]([]string)) map[string](map[string]float64) {
	countsByDate := make(map[string](map[string]int))
	for date := range genomeDatabase {
		//make inner map
		countsByDate[date] = make(map[string]int)
		for i := range genomeDatabase[date] {
			variant := ClassifyVariant(genomeDatabase[date][i])
			countsByDate[date][variant]++
		}
	}

	return NormalizeVariantPercentages(countsByDate)
}

func NormalizeVariantPercentages(countsByDate map[string](map[string]int)) map[string](map[string]float64) {
	frequenciesByDate := make(map[string](map[string]float64))

	for date := range countsByDate {
		//make inner map
		frequenciesByDate[date] = make(map[string]float64)

		//find total number of elements at this date
		total := 0
		for variant := range countsByDate[date] {
			total += countsByDate[date][variant]
		}

		//normalize by the total
		for variant := range countsByDate[date] {
			frequenciesByDate[date][variant] = float64(countsByDate[date][variant]) / float64(total)
		}
	}

	return frequenciesByDate
}
