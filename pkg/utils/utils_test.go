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

package utils_test

import (
	"testing"

	"github.com/Gituser143/cryptgo/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestMinFloat64(t *testing.T) {
	tests := []struct {
		inputVal    []float64
		expectedVal float64
	}{
		{
			inputVal:    []float64{10, 20, 30, 40},
			expectedVal: 10,
		},
		{
			inputVal:    []float64{},
			expectedVal: 0,
		},
		{
			inputVal:    []float64{-10, 20, -30, 40},
			expectedVal: -30,
		},
	}

	for _, test := range tests {
		val := utils.MinFloat64(test.inputVal...)
		assert.Equal(t, test.expectedVal, val)
	}
}

func TestMaxFloat64(t *testing.T) {
	tests := []struct {
		inputVal    []float64
		expectedVal float64
	}{
		{
			inputVal:    []float64{10, 20, 30, 40},
			expectedVal: 40,
		},
		{
			inputVal:    []float64{},
			expectedVal: 0,
		},
		{
			inputVal:    []float64{-10, -20, -30, -40},
			expectedVal: -10,
		},
	}

	for _, test := range tests {
		val := utils.MaxFloat64(test.inputVal...)
		assert.Equal(t, test.expectedVal, val)
	}
}

func TestRoundValues(t *testing.T) {
	tests := []struct {
		name           string
		inputNum1      float64
		inputNum2      float64
		expRoundedVals []float64
		expUnit        string
	}{
		{
			name:           "both values smaller than kilo",
			inputNum1:      250,
			inputNum2:      500,
			expRoundedVals: []float64{250, 500},
			expUnit:        "",
		},
		{
			name:           "both values kilo",
			inputNum1:      80000,
			inputNum2:      9000,
			expRoundedVals: []float64{80, 9},
			expUnit:        "K",
		},
		{
			name:           "mega and less than mega",
			inputNum1:      400000,
			inputNum2:      7000000,
			expRoundedVals: []float64{0.4, 7},
			expUnit:        "M",
		},
		{
			name:           "giga and less than giga",
			inputNum1:      100000000,
			inputNum2:      9000000000,
			expRoundedVals: []float64{0.1, 9},
			expUnit:        "B",
		},
		{
			name:           "both values tera",
			inputNum1:      5000000000000000,
			inputNum2:      100000000000000,
			expRoundedVals: []float64{5000, 100},
			expUnit:        "T",
		},
	}

	for _, test := range tests {
		roundedVals, unit := utils.RoundValues(test.inputNum1, test.inputNum2)
		assert.Equal(t, test.expRoundedVals, roundedVals)
		assert.Equal(t, test.expUnit, unit)
	}
}
