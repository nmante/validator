package validator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"

	"github.com/nmante/validator/types"
)

var (
	ErrInvalidKind = errors.New("Can't call 'len' on this value")
	validKinds     = map[reflect.Kind]int{
		reflect.Array:  1,
		reflect.Chan:   1,
		reflect.Map:    1,
		reflect.Slice:  1,
		reflect.String: 1,
	}
)

// Func is function type that all validator functions must follow
type Func func(interface{}) (FuncResponse, error)

// FuncResponse contains info around if a validator function was valid. If it isn't valid, an
// error message is also returned
type FuncResponse struct {
	IsValid bool
	Error   string
}

// IsTranformableToInt checks if a value is transformable to an int
func IsTransformableToInt(transformer Transformer) Func {
	return IsTransformableTo(transformer, types.Int)
}

// IsTransformableTo checks if a value of type 'A' is transformable to type 'B'
func IsTransformableTo(transformer Transformer, _type reflect.Type) Func {
	return func(v interface{}) (FuncResponse, error) {
		t, err := transformer.Transform(v)
		if err != nil {
			return FuncResponse{}, err
		}

		if reflect.TypeOf(t) != _type {
			return FuncResponse{
				IsValid: false,
				Error:   fmt.Sprintf("%v not transformable to %s", v, _type.Name()),
			}, nil
		}

		return FuncResponse{IsValid: true, Error: ""}, nil
	}
}

// IsEqual transforms a 'v' to a type, and checks if it's equal to 'right'
func IsEqual(transformer Transformer, comparer Comparer, right interface{}) Func {
	return func(v interface{}) (FuncResponse, error) {
		value, err := transformer.Transform(v)
		if err != nil {
			return FuncResponse{}, err
		}

		if ok, err := isTypesEqual(value, right); !ok {
			return FuncResponse{}, err
		}

		if comparer.Compare(value, right) == 0 {
			return FuncResponse{IsValid: true, Error: ""}, nil
		}

		return FuncResponse{IsValid: false, Error: fmt.Sprintf("must be equal to %d", right)}, nil
	}
}

// IsBetween checks if a value is between a lower and an upper value
func IsBetween(transformer Transformer, comparer Comparer, lower interface{}, upper interface{}) Func {
	return func(v interface{}) (FuncResponse, error) {
		value, err := transformer.Transform(v)
		if err != nil {
			return FuncResponse{}, err
		}

		if ok, err := isTypesEqual(value, lower, upper); !ok {
			return FuncResponse{}, err
		}

		if comparer.Compare(lower, value) < 1 && comparer.Compare(value, upper) > -1 {
			return FuncResponse{IsValid: true, Error: ""}, nil
		}

		return FuncResponse{IsValid: false, Error: fmt.Sprintf("must be between %d and %d", lower, upper)}, nil
	}
}

// IsLength checks if the length of an item equals a value
func IsLength(length int) Func {
	return func(v interface{}) (FuncResponse, error) {
		value := reflect.ValueOf(v)
		if _, ok := validKinds[value.Kind()]; !ok {
			return FuncResponse{}, ErrInvalidKind
		}

		if length == value.Len() {
			return FuncResponse{IsValid: true}, nil
		}

		return FuncResponse{
			IsValid: false,
			Error:   fmt.Sprintf("Must have length %d", length),
		}, nil
	}
}

// IsLengthBetween checks if the length of an item is within a range
func IsLengthBetween(lower int, upper int) Func {
	return func(v interface{}) (FuncResponse, error) {
		value := reflect.ValueOf(v)
		if _, ok := validKinds[value.Kind()]; !ok {
			return FuncResponse{}, ErrInvalidKind
		}

		if lower <= value.Len() && value.Len() <= upper {
			return FuncResponse{IsValid: true}, nil
		}

		return FuncResponse{
			IsValid: false,
			Error:   fmt.Sprintf("Must be between length %d and %d", lower, upper),
		}, nil
	}
}

func IsEmail(v interface{}) (FuncResponse, error) {
	email, ok := v.(string)
	if !ok {
		return FuncResponse{}, errors.New("must be a string")
	}

	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	isEmail := re.MatchString(email)

	if isEmail {
		return FuncResponse{IsValid: true}, nil
	}

	return FuncResponse{IsValid: false, Error: "Must be an email"}, nil
}

// isTypesEqual checks if a list of values all have the same type
func isTypesEqual(values ...interface{}) (bool, error) {
	for i := range values {
		if i == len(values)-1 {
			break
		}

		if reflect.TypeOf(values[i]) != reflect.TypeOf(values[i+1]) {
			return false, errors.New("Values must all have the same type")
		}
	}

	return true, nil
}

// IsType checks if a value is of a certain type
func IsType(_type reflect.Type) Func {
	return func(v interface{}) (FuncResponse, error) {
		if reflect.TypeOf(v) != _type {
			return FuncResponse{IsValid: false, Error: fmt.Sprintf("must be a %s", _type.Name())}, nil
		}

		return FuncResponse{IsValid: true, Error: ""}, nil
	}
}

func IsInt(v interface{}) (FuncResponse, error) {
	return IsType(types.Int)(v)
}

func IsInt8(v interface{}) (FuncResponse, error) {
	return IsType(types.Int8)(v)
}
func IsInt16(v interface{}) (FuncResponse, error) {
	return IsType(types.Int16)(v)
}

func IsInt32(v interface{}) (FuncResponse, error) {
	return IsType(types.Int32)(v)
}

func IsInt64(v interface{}) (FuncResponse, error) {
	return IsType(types.Int64)(v)
}

func IsUint(v interface{}) (FuncResponse, error) {
	return IsType(types.Uint)(v)
}

func IsUint8(v interface{}) (FuncResponse, error) {
	return IsType(types.Uint8)(v)
}
func IsUint16(v interface{}) (FuncResponse, error) {
	return IsType(types.Uint16)(v)
}

func IsUint32(v interface{}) (FuncResponse, error) {
	return IsType(types.Uint32)(v)
}

func IsUint64(v interface{}) (FuncResponse, error) {
	return IsType(types.Uint64)(v)
}

func IsUintptr(v interface{}) (FuncResponse, error) {
	return IsType(types.Uintptr)(v)
}

func IsByte(v interface{}) (FuncResponse, error) {
	return IsType(types.Byte)(v)
}

func IsRune(v interface{}) (FuncResponse, error) {
	return IsType(types.Rune)(v)
}

func IsFloat32(v interface{}) (FuncResponse, error) {
	return IsType(types.Float32)(v)
}

func IsFloat64(v interface{}) (FuncResponse, error) {
	return IsType(types.Float64)(v)
}

func IsComplex64(v interface{}) (FuncResponse, error) {
	return IsType(types.Complex64)(v)
}

func IsComplex128(v interface{}) (FuncResponse, error) {
	return IsType(types.Complex128)(v)
}
