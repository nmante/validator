package transform

import (
	"errors"
	"strconv"
)

var (
	ErrNotString = errors.New("Value is not a string")
)

var (
	StringToInt     = stringToInt{}
	StringToFloat32 = stringToFloat32{}
	StringToFloat64 = stringToFloat64{}
	StringToUint    = stringToUint{}
	StringToBool    = stringToBool{}
)

type stringToInt struct{}

// Transform converts a string to an integer
func (s stringToInt) Transform(v interface{}) (interface{}, error) {
	val, ok := v.(string)
	if !ok {
		return 0, ErrNotString
	}

	num, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}

	return num, nil
}

type stringToFloat32 struct{}

// Transform converts a string to a float32
func (s stringToFloat32) Transform(v interface{}) (interface{}, error) {
	val, ok := v.(string)
	if !ok {
		return 0, ErrNotString
	}

	num, err := strconv.ParseFloat(val, 32)
	if err != nil {
		return 0, err
	}

	return num, nil
}

type stringToFloat64 struct{}

// Transform converts a string to a float64
func (s stringToFloat64) Transform(v interface{}) (interface{}, error) {
	val, ok := v.(string)
	if !ok {
		return 0, ErrNotString
	}

	num, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0, err
	}

	return num, nil
}

type stringToUint struct{}

// Transform converts a string to a uint
func (s stringToUint) Transform(v interface{}) (interface{}, error) {
	val, ok := v.(string)
	if !ok {
		return 0, ErrNotString
	}

	num, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0, err
	}

	return num, nil
}

type stringToBool struct{}

// Transform converts a string to a bool
func (s stringToBool) Transform(v interface{}) (interface{}, error) {
	val, ok := v.(string)
	if !ok {
		return 0, ErrNotString
	}

	boolean, err := strconv.ParseBool(val)
	if err != nil {
		return 0, err
	}

	return boolean, nil
}
