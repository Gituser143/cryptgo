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

// MinFloat64 returns minimum float from a given number of floats
func MinFloat64(a ...float64) float64 {
	var min float64
	if len(a) > 0 {
		min = a[0]
	} else {
		return 0
	}

	for _, val := range a {
		if val < min {
			min = val
		}
	}
	return min
}

// MaxFloat64 returns maximum float from a given number of floats
func MaxFloat64(a ...float64) float64 {
	var max float64
	if len(a) > 0 {
		max = a[0]
	} else {
		return 0
	}
	for _, val := range a {
		if val > max {
			max = val
		}
	}
	return max
}
