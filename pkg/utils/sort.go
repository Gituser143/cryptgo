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

	default:
		sortFuncs[sortIdx] = strSort
	}

	// Sort data
	sort.Slice(data, sortFuncs[sortIdx])
}
