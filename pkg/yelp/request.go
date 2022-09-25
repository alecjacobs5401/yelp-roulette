package yelp

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/alecjacobs5401/yelp-roulette/pkg/convert"
)

var withinRegex = regexp.MustCompile(`(\d+)(\w*)`)

// TODO: this should be its own package with better type handling for parsing query strings
// ParseRequest takes a message and parses out the text into a SearchRequest
func ParseRequest(body string) (SearchRequest, error) {
	tokens := strings.Split(body, " ")

	inIndex := -1
	openIndex := -1
	withinIndex := -1
	priceIndices := []int{}
	priceLevels := []string{}
	locationTokens := []string{}
	termTokens := []string{}
	within := 0

	for index, token := range tokens {
		switch strings.ToLower(token) {
		case "in":
			inIndex = index
		case "open":
			openIndex = index
		case "within":
			withinIndex = index
		case "$", "$$", "$$$", "$$$$":
			priceIndices = append(priceIndices, index)
			priceLevels = append(priceLevels, fmt.Sprint(strings.Count(token, "$")))
		default:
			if inIndex >= 0 || withinIndex >= 0 {
				continue
			} else {
				termTokens = append(termTokens, token)
			}
		}
	}

	tokensLength := len(tokens)
	lastTokenIndex := tokensLength - 1
	indices := append(priceIndices, inIndex, withinIndex, openIndex, tokensLength)
	sort.Ints(indices)

	if inIndex >= 0 {
		index := indexOf(indices, inIndex)
		if indices[index] >= lastTokenIndex {
			return SearchRequest{}, errors.New("no location provided")
		}
		nextIndex := indices[index+1]
		locationTokens = tokens[inIndex+1 : nextIndex]
		if len(locationTokens) == 0 {
			return SearchRequest{}, errors.New("no location provided")
		}
	}

	if withinIndex >= 0 {
		index := indexOf(indices, withinIndex)
		if indices[index] >= lastTokenIndex {
			return SearchRequest{}, errors.New("invalid within options provided")
		}
		nextIndex := indices[index+1]
		withinTokens := tokens[withinIndex+1 : nextIndex]

		switch len(withinTokens) {
		case 1:
			matches := withinRegex.FindStringSubmatch(withinTokens[0])
			if matches == nil {
				return SearchRequest{}, errors.New("invalid value provided for within distance")
			}
			value, _ := strconv.ParseFloat(matches[1], 64)
			unit := matches[2]
			if unit == "" {
				within = int(value)
			} else {
				converted, err := convert.ToMeters(value, strings.ToLower(matches[2]))
				if err != nil {
					return SearchRequest{}, err
				}
				within = int(converted)
			}
		case 2:
			value, err := strconv.ParseFloat(withinTokens[0], 64)
			if err != nil {
				return SearchRequest{}, errors.New("invalid value provided for within distance")
			}
			converted, err := convert.ToMeters(value, strings.ToLower(withinTokens[1]))
			if err != nil {
				return SearchRequest{}, err
			}
			within = int(converted)
		default:
			return SearchRequest{}, errors.New("invalid within options provided")
		}
	}

	return SearchRequest{
		Term:     strings.Join(termTokens, " "),
		Location: strings.Join(locationTokens, " "),
		OpenNow:  openIndex >= 0,
		Price:    priceLevels,
		Within:   within,
	}, nil
}

// currently assumes sorted array
func indexOf(items []int, value int) int {
	for index, item := range items {
		if item == value {
			return index
		}

		if item > value {
			return -1
		}
	}
	return -1
}
