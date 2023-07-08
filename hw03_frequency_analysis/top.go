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
	wordsMap := make(map[string]int)
	for _, word := range allWords {
		wordsMap[word]++
	}
	tempSliceOfNotSortedWords := getSliceFromMap(wordsMap)[:10]
	tempSlice := make([]string, 0)
	result := make([]string, 0)
	for i, k := range tempSliceOfNotSortedWords {
		if len(tempSlice) == 0 && i == len(tempSliceOfNotSortedWords)-1 {
			result = append(result, k)
		}
		if len(tempSlice) == 0 {
			tempSlice = append(tempSlice, k)
			continue
		}
		if i == len(tempSliceOfNotSortedWords)-1 && len(tempSlice) > 0 {
			tempSlice = append(tempSlice, k)
			sort.Strings(tempSlice)
			result = append(result, tempSlice...)
			tempSlice = nil
		}
		if len(tempSlice) > 0 && wordsMap[tempSlice[len(tempSlice)-1]] == wordsMap[k] {
			tempSlice = append(tempSlice, k)
			continue
		}
		if len(tempSlice) > 0 && wordsMap[tempSlice[len(tempSlice)-1]] != wordsMap[k] {
			sort.Strings(tempSlice)
			result = append(result, tempSlice...)
			tempSlice = nil
			tempSlice = append(tempSlice, k)
		}
	}
	return result
}

func getSliceFromMap(inputMap map[string]int) []string {
	results := make([]string, 0, len(inputMap))
	for key := range inputMap {
		results = append(results, key)
	}
	sort.SliceStable(results, func(i, j int) bool {
		return inputMap[results[i]] > inputMap[results[j]]
	})
	return results
}
