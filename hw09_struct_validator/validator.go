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
	ErrSetValidation    = errors.New("incorrect validation conditions")
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
		fieldType1 := field
		fmt.Printf("Этот тип %v\n", fieldType1)
		if reflectValue.Field(i).Kind() != reflect.Slice {
			if err := validateValue(tag, varValue); err != nil {
				if checkError(err) {
					vErrors = append(vErrors, ValidationError{Field: name, Err: err})
					continue
				}
				return err
			}
			continue
		}
		for _, sliceVal := range varValue.([]string) {
			if err := validateValue(tag, sliceVal); err != nil {
				if checkError(err) {
					vErrors = append(vErrors, ValidationError{Field: name, Err: err})
					continue
				}
				return err
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
			checkStatus, err := contains(inStr, value)
			if !checkStatus && err == nil {
				return ErrIn
			}
			if err != nil {
				return err
			}
		case "regexp":
			matched, err := regexp.MatchString(val, fmt.Sprintf("%v", value))
			if err != nil {
				return err
			}
			if !matched {
				return ErrRegexp
			}
		}
	}
	return nil
}

func contains(s []string, e interface{}) (bool, error) {
	for _, a := range s {
		switch reflect.ValueOf(e).Kind() {
		case reflect.String:
			if a == fmt.Sprint(e) {
				return true, nil
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
			v, err := strconv.Atoi(a)
			if err != nil {
				return false, ErrSetValidation
			}
			if v == e.(int) {
				return true, nil
			}
		default:
		}
	}
	return false, nil
}

func checkError(err error) bool {
	if errors.Is(err, ErrLen) ||
		errors.Is(err, ErrMin) ||
		errors.Is(err, ErrMax) ||
		errors.Is(err, ErrIn) ||
		errors.Is(err, ErrValueIsNotStruct) ||
		errors.Is(err, ErrWrongRule) ||
		errors.Is(err, ErrRegexp) ||
		errors.Is(err, ErrInternal) {
		return true
	}
	return false
}
