package validator

import (
	"github.com/nmante/validator/compare"
	"github.com/nmante/validator/transform"
	"github.com/nmante/validator/types"
)

// IsStringInRangeInts checks if a string value is between integers
func IsStringInRangeInts(lower int, upper int) Func {
	return IsBetween(transform.StringToInt, compare.Int, lower, upper)
}

// IsStringEqualToInt checks if the value within a string is equal to an integer
func IsStringEqualToInt(right int) Func {
	return IsEqual(transform.StringToInt, compare.Int, right)
}

// IsStringInt checks if the value within a string is an integer
func IsStringInt(v interface{}) (FuncResponse, error) {
	return IsTransformableTo(transform.StringToInt, types.Int)(v)
}

// IsStringFloat32 checks if the value within a string is a float
func IsStringFloat32(v interface{}) (FuncResponse, error) {
	return IsTransformableTo(transform.StringToFloat32, types.Float32)(v)
}

// IsStringFloat64 checks if the value within a string is a float
func IsStringFloat64(v interface{}) (FuncResponse, error) {
	return IsTransformableTo(transform.StringToFloat64, types.Float64)(v)
}

// IsStringBool checks if the value within a string is a float
func IsStringBool(v interface{}) (FuncResponse, error) {
	return IsTransformableTo(transform.StringToBool, types.Bool)(v)
}

// IsStringUint checks if the value within a string is a float
func IsStringUint(v interface{}) (FuncResponse, error) {
	return IsTransformableTo(transform.StringToUint, types.Uint)(v)
}
