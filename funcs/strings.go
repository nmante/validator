package funcs

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/nmante/validator/compare"
	"github.com/nmante/validator/transform"
	"github.com/nmante/validator/types"
)

var (
	String = _string{}
)

type _string struct{}

// IsInRangeInts checks if a string value is between integers
func (s _string) IsInRangeInts(lower int, upper int) Func {
	return IsBetween(transform.StringToInt, compare.Int, lower, upper)
}

type IsStringBetweenInts struct {
	value string
	upper int
	lower int
}

func (i IsStringBetweenInts) Validate() (Response, error) {
	parsed, err := strconv.ParseInt(i.value, 10, 0)
	if err != nil {
		return Response{}, err
	}

	if i.lower <= int(parsed) && int(parsed) <= i.upper {
		return Response{IsValid: true, Error: ""}, nil
	}

	return Response{IsValid: false, Error: fmt.Sprintf("must be between %d and %d", i.lower, i.upper)}, nil
}

// IsEmail checks if a string is an email address via regex
func (s _string) IsEmail(v interface{}) (Response, error) {
	email, ok := v.(string)
	if !ok {
		return Response{}, errors.New("must be a string")
	}

	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	isEmail := re.MatchString(email)

	if isEmail {
		return Response{IsValid: true}, nil
	}

	return Response{IsValid: false, Error: "Must be an email address"}, nil
}

// IsEqualToInt checks if the value within a string is equal to an integer
func (s _string) IsEqualToInt(right int) Func {
	return IsEqual(transform.StringToInt, compare.Int, right)
}

// IsStringInt checks if the value within a string is an integer
func (s _string) IsInt(v interface{}) (Response, error) {
	return IsTransformableTo(transform.StringToInt, types.Int)(v)
}

// IsStringFloat32 checks if the value within a string is a 64 bit float
func (s _string) IsFloat32(v interface{}) (Response, error) {
	return IsTransformableTo(transform.StringToFloat32, types.Float32)(v)
}

// IsStringFloat64 checks if the value within a string is a 64 bit float
func (s _string) IsFloat64(v interface{}) (Response, error) {
	return IsTransformableTo(transform.StringToFloat64, types.Float64)(v)
}

// IsBool checks if the value within a string is a boolean
func (s _string) IsBool(v interface{}) (Response, error) {
	return IsTransformableTo(transform.StringToBool, types.Bool)(v)
}

// IsUint checks if the value within a string is an unsigned integer
func (s _string) IsUint(v interface{}) (Response, error) {
	return IsTransformableTo(transform.StringToUint, types.Uint)(v)
}
