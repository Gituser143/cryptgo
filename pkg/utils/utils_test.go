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
