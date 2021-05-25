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

import (
	"sort"
	"strconv"
	"strings"
)

const (
	UP_ARROW   = "▲"
	DOWN_ARROW = "▼"
)

// SortData helps sort table rows. It sorts the table based on values given
// in the sortIdx column and sorts ascending if sortAsc is true.
// sortCase is set to identify the set of 'less' functions to use to
// sort the selected column by.
func SortData(data [][]string, sortIdx int, sortAsc bool, sortCase string) {

	if sortIdx < 0 {
		return
	}

	// Define less functions
	intSort := func(i, j int) bool {
		x, _ := strconv.Atoi(data[i][sortIdx])
		y, _ := strconv.Atoi(data[j][sortIdx])
		if sortAsc {
			return x < y
		}
		return x > y
	}

	strSort := func(i, j int) bool {
		if sortAsc {
			return data[i][sortIdx] < data[j][sortIdx]
		}
		return data[i][sortIdx] > data[j][sortIdx]
	}

	floatSort := func(i, j int) bool {
		x1 := data[i][sortIdx]
		y1 := data[j][sortIdx]
		x, _ := strconv.ParseFloat(x1, 32)
		y, _ := strconv.ParseFloat(y1, 32)
		if sortAsc {
			return x < y
		}
		return x > y
	}

	changeSort := func(i, j int) bool {
		x1 := strings.Split(data[i][sortIdx], " ")
		y1 := strings.Split(data[j][sortIdx], " ")
		x, _ := strconv.ParseFloat(x1[1], 64)
		if string(x1[0]) == DOWN_ARROW {
			x = -x
		}

		y, _ := strconv.ParseFloat(y1[1], 64)
		if string(y1[0]) == DOWN_ARROW {
			y = -y
		}

		if sortAsc {
			return x < y
		}
		return x > y
	}

	// Set function map
	sortFuncs := make(map[int]func(i, j int) bool)
	switch sortCase {
	case "COINS":
		sortFuncs = map[int]func(i, j int) bool{
			0: intSort,    // Rank
			1: strSort,    // Symbol
			2: floatSort,  // Price
			3: changeSort, // Change %
		}

	case "FAVOURITES":
		sortFuncs = map[int]func(i, j int) bool{
			0: strSort,   // Symbol
			1: floatSort, // Price
		}

	case "PORTFOLIO":
		sortFuncs = map[int]func(i, j int) bool{
			0: strSort,   // Coin
			1: strSort,   // Symbol
			2: floatSort, // Price
			3: floatSort, // Holding
			4: floatSort, // Balance
		}
	default:
		sortFuncs[sortIdx] = strSort
	}

	if _, ok := sortFuncs[sortIdx]; !ok {
		sortIdx = 0
	}

	// Sort data
	sort.Slice(data, sortFuncs[sortIdx])
}
