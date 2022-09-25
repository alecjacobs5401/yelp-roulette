package convert

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToMeters(t *testing.T) {
	testCases := []struct {
		name          string
		value         float64
		originalUnit  string
		expectedValue float64
		expectedError error
	}{
		{
			name:          "from miles",
			value:         25,
			originalUnit:  "miles",
			expectedValue: 40233.60,
		},
		{
			name:          "from mile",
			value:         10,
			originalUnit:  "miles",
			expectedValue: 16093.44,
		},
		{
			name:          "from mi",
			value:         15,
			originalUnit:  "mi",
			expectedValue: 24140.16,
		},
		{
			name:          "from km",
			value:         25,
			originalUnit:  "km",
			expectedValue: 25000,
		},
		{
			name:          "from kilometer",
			value:         25,
			originalUnit:  "kilometer",
			expectedValue: 25000,
		},
		{
			name:          "from kilometre",
			value:         25,
			originalUnit:  "kilometre",
			expectedValue: 25000,
		},
		{
			name:          "from kilometers",
			value:         10,
			originalUnit:  "kilometers",
			expectedValue: 10000,
		},
		{
			name:          "from kilometres",
			value:         15,
			originalUnit:  "kilometres",
			expectedValue: 15000,
		},
		{
			name:          "from m",
			value:         25,
			originalUnit:  "m",
			expectedValue: 25,
		},
		{
			name:          "from meter",
			value:         25,
			originalUnit:  "meter",
			expectedValue: 25,
		},
		{
			name:          "from metre",
			value:         25,
			originalUnit:  "metre",
			expectedValue: 25,
		},
		{
			name:          "from meters",
			value:         10,
			originalUnit:  "meters",
			expectedValue: 10,
		},
		{
			name:          "from metres",
			value:         15,
			originalUnit:  "metres",
			expectedValue: 15,
		},
		{
			name:          "from unknown unit",
			value:         15,
			originalUnit:  "random-unit",
			expectedError: errors.New("unknown unit provided: \"random-unit\""),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			value, err := ToMeters(testCase.value, testCase.originalUnit)
			if testCase.expectedError != nil && assert.Error(t, err) {
				assert.Equal(t, err, testCase.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectedValue, value)
			}
		})
	}
}
