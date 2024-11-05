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

func (v ValidationErrors) Error() string {
	errs := make([]string, 0, len(v))
	for _, err := range v {
		errs = append(errs, fmt.Sprintf("%s: %s\n", err.Field, err.Err))
	}
	return strings.Join(errs, "")
}

func Validate(v interface{}) error {
	var validationErrors ValidationErrors

	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		return errors.New("validation input must be a struct type")
	}

	// iterating each field of the struct
	for i := 0; i < val.NumField(); i++ {
		valTypeField := val.Type().Field(i)
		valField := val.Field(i)

		tag := valTypeField.Tag.Get("validate")
		if tag == "" {
			continue
		}

		if err := validateRules(strings.Split(tag, "|"), valField, validationErrors, valTypeField.Name); err != nil {
			return err
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

func validateRules(rules []string, valField reflect.Value, valErrors ValidationErrors, varTypeFieldName string) error {
	// iterate over each rule
	for _, rule := range rules {
		ruleKeyVal := strings.Split(rule, ":") // ex: ["min"]["18"]
		if len(ruleKeyVal) != 2 {
			return fmt.Errorf("invalid rule format in: %s", rule)
		}

		//nolint:exhaustive
		switch valField.Kind() {
		case reflect.String:
			if err := validateString(valField.String(), ruleKeyVal[0], ruleKeyVal[1]); err != nil {
				valErrors = append(valErrors, ValidationError{Field: varTypeFieldName, Err: err})
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if err := validateInteger(int(valField.Int()), ruleKeyVal[0], ruleKeyVal[1]); err != nil {
				valErrors = append(valErrors, ValidationError{Field: varTypeFieldName, Err: err})
			}
		case reflect.Slice:
			for j := 0; j < valField.Len(); j++ {
				elem := valField.Index(j)
				if elem.Kind() == reflect.String {
					if err := validateString(elem.String(), ruleKeyVal[0], ruleKeyVal[1]); err != nil {
						valErrors = append(valErrors, ValidationError{Field: varTypeFieldName, Err: err})
					}
				} else if elem.Kind() == reflect.Int {
					if err := validateInteger(int(elem.Int()), ruleKeyVal[0], ruleKeyVal[1]); err != nil {
						valErrors = append(valErrors, ValidationError{Field: varTypeFieldName, Err: err})
					}
				}
			}
		default:
			return fmt.Errorf("unsupported field type: %s", valField.Kind())
		}
	}
	return nil
}

func validateString(s, key, val string) error {
	switch key {
	case "len":
		length, err := strconv.Atoi(val)
		if err != nil {
			return fmt.Errorf("invalid len value: %s", val)
		}
		if len(s) != length {
			return fmt.Errorf("length must be %d", length)
		}
	case "regexp":
		re, err := regexp.Compile(val)
		if err != nil {
			return fmt.Errorf("invalid regexp: %s", val)
		}
		if !re.MatchString(s) {
			return fmt.Errorf("must match regexp %s", val)
		}
	case "in":
		options := strings.Split(val, ",")
		for _, option := range options {
			if s == option {
				return nil
			}
		}
		return fmt.Errorf("must be one of [%s]", strings.Join(options, ", "))
	}
	return nil
}

func validateInteger(n int, key, val string) error {
	switch key {
	case "min":
		minVal, err := strconv.Atoi(val)
		if err != nil {
			return fmt.Errorf("invalid min value: %s", val)
		}
		if n < minVal {
			return fmt.Errorf("must be >= %d", minVal)
		}
	case "max":
		maxVal, err := strconv.Atoi(val)
		if err != nil {
			return fmt.Errorf("invalid max value: %s", val)
		}
		if n > maxVal {
			return fmt.Errorf("must be <= %d", maxVal)
		}
	case "in":
		options := strings.Split(val, ",")
		for _, option := range options {
			optVal, err := strconv.Atoi(option)
			if err != nil {
				return fmt.Errorf("invalid in value: %s", option)
			}
			if n == optVal {
				return nil
			}
		}
		return fmt.Errorf("must be one of [%s]", strings.Join(options, ", "))
	}
	return nil
}
