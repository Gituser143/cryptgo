/*
Copyright © 2021 Bhargav SNV bhargavsnv100@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

import "math"

var (
	K = math.Pow(10, 3)
	M = math.Pow(10, 6)
	G = math.Pow(10, 9)
	T = math.Pow(10, 12)
	Q = math.Pow(10, 15)
)

func roundOffNearestTen(num float64, divisor float64) float64 {
	x := num / divisor
	return math.Round(x*10) / 10
}

// Round values rounds off a pair of given floats to Thousands (K),
// Millions (M) or Billions (B).
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

	case n >= Q:
		nums = append(nums, roundOffNearestTen(num1, T))
		nums = append(nums, roundOffNearestTen(num2, T))
		units = "T"

	default:
		return []float64{num1, num2}, ""
	}

	return nums, units
}
