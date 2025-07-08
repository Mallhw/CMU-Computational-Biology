package main

func InitializeMatrix(numRows int, numCols int) [][]int {
	mtx := make([][]int, numRows)

	for i := 0; i < numRows; i++ {
		mtx[i] = make([]int, numCols)
	}

	return mtx
}

// Min takes an arbitrary number of integers and returns their minimum.
func Min(nums ...int) int {
	if len(nums) == 0 {
		panic("Error: no inputs given to Max.")
	}
	m := nums[0]
	// nums gets converted to an array
	for _, val := range nums {
		if val < m {
			// update m
			m = val
		}
	}
	return m
}

// MaxInts
// Input: an arbitrary number of integers
// Output: their maximum
func MaxInts(nums ...int) int {
	//nums is a slice
	if len(nums) == 0 {
		panic("No integers given.")
	}
	m := 0

	for i, val := range nums {
		if i == 0 || val > m {
			// update m if we are at the first integer
			// OR if current element is bigger than known max
			m = val
		}
	}

	return m
}

// InitializeFloatMatrix
// Input: integers numRows and numCols
// Output: a two-dimensional slice of floats, initialized to zero.
func InitializeFloatMatrix(numRows, numCols int) [][]float64 {
	mtx := make([][]float64, numRows)
	for i := range mtx {
		mtx[i] = make([]float64, numCols)
	}

	return mtx
}

// MaxFloat takes an arbitrary number of floats and returns their maximum.
func MaxFloat(nums ...float64) float64 {
	//nums now is an array of integers
	if len(nums) == 0 {
		panic("Error: need to give at least one value as input.")
	}
	m := 0.0

	// tricky issue: make sure that this works if every input integer is < 0
	//if nums is an array, let's range over it
	for i, val := range nums {
		if i == 0 || val > m {
			//update m if we are at the first parameter
			//or if the current value beats the largest known integer
			m = val
		}
	}

	return m
}
