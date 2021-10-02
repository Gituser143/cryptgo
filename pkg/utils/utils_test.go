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
