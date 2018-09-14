package main

import (
	"fmt"
	"strconv"
)

// IsStringBetweenInts checks if a string value is between integers
func IsStringBetweenInts(lower int, upper int) Func {
	return func(v interface{}) (FuncResponse, error) {
		r, err := IsStringInt(v)
		if err != nil {
			return FuncResponse{}, err
		} else if !r.IsValid {
			return r, nil
		}

		num, _ := strconv.Atoi(v.(string))
		if lower < num && num < upper {
			return FuncResponse{true, ""}, nil
		}

		return FuncResponse{false, fmt.Sprintf("must be between %d and %d", lower, upper)}, nil
	}
}

// IsStringEqual checks if two values are equal
func IsStringEqualToInt(b int) Func {
	return func(v interface{}) (FuncResponse, error) {
		r, err := IsStringInt(v)
		if err != nil {
			return FuncResponse{}, err
		} else if !r.IsValid {
			return r, nil
		}

		if num, _ := strconv.Atoi(v.(string)); num == b {
			return FuncResponse{true, ""}, nil
		}

		return FuncResponse{false, fmt.Sprintf("must be equal to %d", b)}, nil
	}
}

// IsStringInt checks if a string is an integer
func IsStringInt(v interface{}) (FuncResponse, error) {
	val, ok := v.(string)
	if !ok {
		return FuncResponse{}, ErrTypeMismatch{v, "string"}
	}

	if _, err := strconv.Atoi(val); err != nil {
		return FuncResponse{false, "must be an integer"}, nil
	}

	return FuncResponse{true, ""}, nil
}
