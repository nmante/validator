package validator

import (
	"strconv"
)

type StringToInt struct{}

// Transform converts a string to an integer
func (s StringToInt) Transform(v interface{}) (interface{}, error) {
	val, ok := v.(string)
	if !ok {
		return 0, ErrTypeMismatch{v, "string"}
	}

	num, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}

	return num, nil
}

// IsStringInRangeInts checks if a string value is between integers
func IsStringInRangeInts(lower int, upper int) Func {
	return IsBetween(StringToInt{}, IntComparer{}, lower, upper)
}

// IsStringEqualToInt checks if the value within a string is equal to an integer
func IsStringEqualToInt(right int) Func {
	return IsEqual(StringToInt{}, IntComparer{}, right)
}

// IsStringInt checks if the value within a string is an integer
func IsStringInt(v interface{}) (FuncResponse, error) {
	return IsTransformableToInt(StringToInt{})(v)
}
