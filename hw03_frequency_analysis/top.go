package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(inputText string) []string {
	if inputText == "" {
		return []string{}
	}

	allWords := strings.Fields(inputText)

	wordToCountMap := make(map[string]int)
	for _, word := range allWords {
		wordToCountMap[word]++
	}

	resultSlice := getSliceFromMap(wordToCountMap)
	sort.Slice(resultSlice, func(i, j int) bool {
		if wordToCountMap[resultSlice[i]] == wordToCountMap[resultSlice[j]] {
			return resultSlice[i] < resultSlice[j]
		}
		return wordToCountMap[resultSlice[i]] > wordToCountMap[resultSlice[j]]
	})

	max := 10
	if l := len(resultSlice); l < max {
		max = l
	}
	return resultSlice[:max]
}

func getSliceFromMap(inputMap map[string]int) []string {
	result := make([]string, 0, len(inputMap))
	for key := range inputMap {
		result = append(result, key)
	}
	sort.SliceStable(result, func(i, j int) bool {
		return inputMap[result[i]] > inputMap[result[j]]
	})
	return result
}
