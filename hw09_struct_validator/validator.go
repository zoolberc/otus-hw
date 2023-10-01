package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

var (
	ErrLen              = errors.New("invalid value length")
	ErrMin              = errors.New("value is less than the minimum allowable")
	ErrMax              = errors.New("value is greater than the maximum allowable")
	ErrIn               = errors.New("invalid value passed")
	ErrValueIsNotStruct = errors.New("value is not a structure")
	ErrWrongRule        = errors.New("wrong rule")
	ErrRegexp           = errors.New("value is invalid according to the regular sequence")
	ErrInternal         = errors.New("internal error")
)

func (v ValidationErrors) Error() string {
	buff := strings.Builder{}
	for _, valErr := range v {
		buff.WriteString(fmt.Sprintf("field: %s, error: %v\n", valErr.Field, valErr.Err))
	}
	return buff.String()
}

func Validate(v interface{}) error {
	reflectValue := reflect.ValueOf(v)

	if reflectValue.Kind() != reflect.Struct {
		return ErrValueIsNotStruct
	}

	vErrors := ValidationErrors{}

	for i := 0; i < reflectValue.NumField(); i++ {
		field := reflectValue.Type().Field(i)
		name := field.Name
		tag := field.Tag
		if len(tag) == 0 {
			continue
		}

		varValue := reflectValue.Field(i).Interface()

		if reflectValue.Field(i).Kind() != reflect.Slice {
			err := validateValue(tag, varValue)
			if err != nil {
				vErrors = append(vErrors, ValidationError{Field: name, Err: err})
			}
		} else {
			for _, sliceVal := range varValue.([]string) {
				err := validateValue(tag, sliceVal)
				if err != nil {
					vErrors = append(vErrors, ValidationError{Field: name, Err: err})
				}
			}
		}
	}

	if len(vErrors) > 0 {
		return vErrors
	}

	return nil
}

func validateValue(tag reflect.StructTag, value interface{}) error {
	tagVal := tag.Get("validate")
	if len(tagVal) == 0 {
		return nil
	}

	valRules := strings.Split(tagVal, "|")

	for _, rawRule := range valRules {
		valRule := strings.Split(rawRule, ":")

		if len(valRule) != 2 {
			return ErrWrongRule
		}

		rule := valRule[0]
		val := valRule[1]
		intVal, _ := strconv.Atoi(val)
		switch rule {
		case "len":
			if intVal != len(value.(string)) {
				return ErrLen
			}
		case "min":
			if intVal > value.(int) {
				return ErrMin
			}
		case "max":
			if intVal < value.(int) {
				return ErrMax
			}
		case "in":
			inStr := strings.Split(val, ",")
			if !contains(inStr, fmt.Sprintf("%v", value)) {
				return ErrIn
			}
		case "regexp":
			matched, _ := regexp.MatchString(val, fmt.Sprintf("%v", value))
			if !matched {
				return ErrRegexp
			}
		}
	}
	return nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
