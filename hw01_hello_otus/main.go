package main

import "fmt"

func main() {
	//reversWord := stringutil.Reverse("Hello, OTUS!")
	//fmt.Println(reversWord)
	wordsMap := make(map[string]int)

	for i := 0; i < 10; i++ {
		wordsMap[string(rune(i))]++
	}

	for i := range wordsMap {
		fmt.Printf("Key = %s, value = %d\n", i, wordsMap[i])
	}
}
