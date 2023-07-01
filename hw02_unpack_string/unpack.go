package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(inputWord string) (string, error) {
	var result strings.Builder
	if inputWord == "" {
		return "", nil
	}
	_, err := strconv.Atoi(string(inputWord[0]))
	if err == nil {
		return "", ErrInvalidString
	}
	preValue := ""
	lenInputWord := len([]rune(inputWord))
	for i, s := range inputWord {
		curr, err := strconv.Atoi(string(s))
		_, errConvPreValue := strconv.Atoi(preValue)
		if err != nil {
			if errConvPreValue != nil {
				result.WriteString(preValue)
			}
			if i == lenInputWord-1 {
				result.WriteString(string(s))
			}
			preValue = string(s)
		} else {
			if errConvPreValue == nil {
				return "", ErrInvalidString
			}
			result.WriteString(strings.Repeat(preValue, curr))
			preValue = string(s)
		}
	}
	return result.String(), nil
}
