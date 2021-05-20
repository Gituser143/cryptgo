package utils

func MinFloat64(a ...float64) float64 {
	min := a[0]
	for _, val := range a {
		if val < min {
			min = val
		}
	}
	return min
}

func MaxFloat64(a ...float64) float64 {
	max := a[0]
	for _, val := range a {
		if val > max {
			max = val
		}
	}
	return max
}
