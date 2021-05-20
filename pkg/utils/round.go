package utils

import "math"

var (
	K = math.Pow(10, 3)
	M = math.Pow(10, 6)
	G = math.Pow(10, 9)
	T = math.Pow(10, 12)
)

func roundOffNearestTen(num float64, divisor float64) float64 {
	x := num / divisor
	return math.Round(x*10) / 10
}

func RoundValues(num1, num2 float64) ([]float64, string) {
	nums := []float64{}
	var units string
	var n float64
	if num1 > num2 {
		n = num1
	} else {
		n = num2
	}

	switch {
	case n < K:
		nums = append(nums, num1)
		nums = append(nums, num2)
		units = ""

	case n < M:
		nums = append(nums, roundOffNearestTen(num1, K))
		nums = append(nums, roundOffNearestTen(num2, K))
		units = "K"

	case n < G:
		nums = append(nums, roundOffNearestTen(num1, M))
		nums = append(nums, roundOffNearestTen(num2, M))
		units = "M"

	case n < T:
		nums = append(nums, roundOffNearestTen(num1, G))
		nums = append(nums, roundOffNearestTen(num2, G))
		units = "B"
	}

	return nums, units
}
