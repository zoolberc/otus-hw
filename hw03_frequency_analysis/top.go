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

	sort.Slice(allWords, func(i, j int) bool {
		if wordToCountMap[allWords[i]] == wordToCountMap[allWords[j]] {
			return allWords[i] < allWords[j]
		}
		return wordToCountMap[allWords[i]] > wordToCountMap[allWords[j]]
	})

	max := 10
	if l := len(allWords); l < max {
		max = l
	}
	return allWords[:max]
	//countToWordsMap := make(map[int][]string)
	//for word, count := range wordToCountMap {
	//	words := countToWordsMap[count]
	//	words = append(words, word)
	//	countToWordsMap[count] = words
	//}
	//
	//counts := make([]int, 0, len(countToWordsMap))
	//for count := range countToWordsMap {
	//	counts = append(counts, count)
	//}
	//sort.Ints(counts)
	//
	//result := make([]string, 0, 10)
	//for i := len(counts) - 1; i >= 0 && len(result) < 10; i-- {
	//	words := countToWordsMap[counts[i]]
	//	sort.Strings(words)
	//	for _, word := range words {
	//		if len(result) < 10 {
	//			result = append(result, word)
	//		} else {
	//			break
	//		}
	//	}
	//}

	//return result
}
