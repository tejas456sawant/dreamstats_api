package utils

func ArraySum(arr []int) int {
	sum := 0
	for _, v := range arr {
		sum += v
	}
	return sum
}

func ArraySumFloat(arr []float64) float64 {
	sum := 0.0
	for _, v := range arr {
		sum += v
	}
	return sum
}

func ArrayAvgFloat(arr []float64) float64 {
	sum := 0.0
	for _, v := range arr {
		sum += v
	}
	return sum / float64(len(arr))
}
