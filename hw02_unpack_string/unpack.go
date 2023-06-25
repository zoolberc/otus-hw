package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(inputWord string) (string, error) {
	var result strings.Builder
	inputWordRune := []rune(inputWord)
	if len(inputWordRune) == 0 {
		return "", nil
	}
	preValue := ""
	_, err := strconv.Atoi(string(inputWordRune[0]))
	if err == nil {
		return "", ErrInvalidString
	}
	for i, s := range inputWordRune {
		curr, err := strconv.Atoi(string(s))
		if err != nil {
			if i == len(inputWordRune)-1 {
				_, err := strconv.Atoi(preValue)
				if err != nil {
					result.WriteString(preValue)
				}
				result.WriteString(string(s))
			} else {
				_, err := strconv.Atoi(preValue)
				if err == nil {
					preValue = string(s)
				} else {
					result.WriteString(preValue)
					preValue = string(s)
				}
			}
		} else {
			_, err := strconv.Atoi(preValue)
			if err != nil {
				result.WriteString(strings.Repeat(preValue, curr))
				preValue = string(s)
			} else {
				return "", ErrInvalidString
			}
		}
	}
	return result.String(), nil
}
