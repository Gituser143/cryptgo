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

package allcoin

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterRows(t *testing.T) {
	tests := []struct {
		allRows      [][]string
		filter       string
		filteredRows [][]string
	}{
		{
			allRows: [][]string{
				{"", "BTC", "", "", "", "bitcoin"},
				{"", "ETH", "", "", "", "ethereum"},
				{"", "LTC", "", "", "", "litecoin"},
				{"", "BTC", "", "", "", "bitcoin"},
				{"", "BTC", "", "", "", "bitcoin"},
			},
			filter: "bitcoin",
			filteredRows: [][]string{
				{"", "BTC", "", "", "", "bitcoin"},
				{"", "BTC", "", "", "", "bitcoin"},
				{"", "BTC", "", "", "", "bitcoin"},
			},
		},
		{
			allRows: [][]string{
				{"", "BTC", "", "", "", "bitcoin"},
				{"", "ETH", "", "", "", "ethereum"},
				{"", "LTC", "", "", "", "litecoin"},
				{"", "BTC", "", "", "", "bitcoin"},
				{"", "BTC", "", "", "", "bitcoin"},
			},
			filter: "LTC",
			filteredRows: [][]string{
				{"", "LTC", "", "", "", "litecoin"},
			},
		},
		{
			allRows: [][]string{
				{"", "XRP", "", "", "", "ripple"},
				{"", "ETH", "", "", "", "ethereum"},
				{"", "LTC", "", "", "", "litecoin"},
				{"", "BTC", "", "", "", "bitcoin"},
			},
			filter:       "ABA",
			filteredRows: [][]string{},
		},
	}

	for _, test := range tests {
		var mutex = &sync.Mutex{}
		filteredRows := filterRows(test.allRows, test.filter, mutex)
		assert.Equal(t, test.filteredRows, filteredRows)
	}
}
