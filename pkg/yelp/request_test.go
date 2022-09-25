package yelp

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseRequest(t *testing.T) {
	testCases := []struct {
		name            string
		message         string
		expectedRequest SearchRequest
		expectedError   error
	}{
		{
			name:    "search term only",
			message: "dinner",
			expectedRequest: SearchRequest{
				Term:  "dinner",
				Price: []string{},
			},
		},
		{
			name:    "single value location only",
			message: "in NYC",
			expectedRequest: SearchRequest{
				Location: "NYC",
				Price:    []string{},
			},
		},
		{
			name:    "multi value location only",
			message: "in Santa Barbara, CA",
			expectedRequest: SearchRequest{
				Location: "Santa Barbara, CA",
				Price:    []string{},
			},
		},
		{
			name:    "$ only",
			message: "$",
			expectedRequest: SearchRequest{
				Price: []string{"1"},
			},
		},
		{
			name:    "$$ only",
			message: "$$",
			expectedRequest: SearchRequest{
				Price: []string{"2"},
			},
		},
		{
			name:    "$$$ only",
			message: "$$$",
			expectedRequest: SearchRequest{
				Price: []string{"3"},
			},
		},
		{
			name:    "$$$$ only",
			message: "$$$$",
			expectedRequest: SearchRequest{
				Price: []string{"4"},
			},
		},
		{
			name:    "with open only",
			message: "open",
			expectedRequest: SearchRequest{
				OpenNow: true,
				Price:   []string{},
			},
		},
		{
			name:    "with within only with space between length and unit",
			message: "within 10 mi",
			expectedRequest: SearchRequest{
				Within: 16093,
				Price:  []string{},
			},
		},
		{
			name:    "with within that has capital units",
			message: "within 10 MI",
			expectedRequest: SearchRequest{
				Within: 16093,
				Price:  []string{},
			},
		},
		{
			name:    "with within that has capital units no space",
			message: "within 10MI",
			expectedRequest: SearchRequest{
				Within: 16093,
				Price:  []string{},
			},
		},
		{
			name:    "with within only with no space between length and unit",
			message: "within 1000m",
			expectedRequest: SearchRequest{
				Within: 1000,
				Price:  []string{},
			},
		},
		{
			name:    "with within only with value",
			message: "within 1500",
			expectedRequest: SearchRequest{
				Within: 1500,
				Price:  []string{},
			},
		},
		{
			name:          "with invalid within value",
			message:       "within distance",
			expectedError: errors.New("invalid value provided for within distance"),
		},
		{
			name:    "with term and location",
			message: "dinner in NYC",
			expectedRequest: SearchRequest{
				Term:     "dinner",
				Location: "NYC",
				Price:    []string{},
			},
		},
		{
			name:    "with term and long location",
			message: "dinner in Santa Barbara, CA",
			expectedRequest: SearchRequest{
				Term:     "dinner",
				Location: "Santa Barbara, CA",
				Price:    []string{},
			},
		},
		{
			name:    "with term, long location, and price options",
			message: "dinner in Santa Barbara, CA $ $$$",
			expectedRequest: SearchRequest{
				Term:     "dinner",
				Location: "Santa Barbara, CA",
				Price:    []string{"1", "3"},
			},
		},
		{
			name:    "with term, long location, price options, and open",
			message: "dinner in Santa Barbara, CA $ open $$$ $$",
			expectedRequest: SearchRequest{
				Term:     "dinner",
				Location: "Santa Barbara, CA",
				OpenNow:  true,
				Price:    []string{"1", "3", "2"},
			},
		},
		{
			name:    "with term, long location, price options, and open",
			message: "dinner in Santa Barbara, CA $ open $$$ $$",
			expectedRequest: SearchRequest{
				Term:     "dinner",
				Location: "Santa Barbara, CA",
				OpenNow:  true,
				Price:    []string{"1", "3", "2"},
			},
		},
		{
			name:    "with term, long location, price options, open, and within",
			message: "dinner in Santa Barbara, CA $ open $$$ $$ within 10mi",
			expectedRequest: SearchRequest{
				Term:     "dinner",
				Location: "Santa Barbara, CA",
				OpenNow:  true,
				Price:    []string{"1", "3", "2"},
				Within:   16093,
			},
		},
		{
			name:    "with within before in",
			message: "dinner within 1000m in Santa Barbara, CA $ open $$$ $$",
			expectedRequest: SearchRequest{
				Term:     "dinner",
				Location: "Santa Barbara, CA",
				OpenNow:  true,
				Price:    []string{"1", "3", "2"},
				Within:   1000,
			},
		},
		{
			name:          "with in provided without location at end of query",
			message:       "dinner in",
			expectedError: errors.New("no location provided"),
		},
		{
			name:          "with in provided without location",
			message:       "dinner in $ $$ open",
			expectedError: errors.New("no location provided"),
		},
		{
			name:          "with within provided without value at end of query",
			message:       "dinner within",
			expectedError: errors.New("invalid within options provided"),
		},
		{
			name:          "with within provided without value",
			message:       "dinner within $ $$ open",
			expectedError: errors.New("invalid within options provided"),
		},
		{
			name:          "with within provided without value",
			message:       "dinner within $ $$ open",
			expectedError: errors.New("invalid within options provided"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			request, err := ParseRequest(testCase.message)

			if testCase.expectedError != nil && assert.Error(t, err) {
				assert.Equal(t, testCase.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectedRequest, request)
			}
		})
	}
}
